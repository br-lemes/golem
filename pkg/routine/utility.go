package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

const utilityMaxStack = 100
const utilityMinThreshold = 10

func Utility(character schemas.CharacterSchema, utility1, utility2 string) (schemas.CharacterSchema, error) {
	return utility(defaultDeps, character, utility1, utility2)
}

type utilitySlot struct {
	slot    schemas.ItemSlot
	code    string
	curCode string
	curQty  int
}

func utility(d deps, character schemas.CharacterSchema, utility1, utility2 string) (schemas.CharacterSchema, error) {
	slots := []*utilitySlot{
		{
			"utility1",
			utility1,
			character.Utility1Slot,
			character.Utility1SlotQuantity,
		},
		{
			"utility2",
			utility2,
			character.Utility2Slot,
			character.Utility2SlotQuantity,
		},
	}

	var needsAttention []*utilitySlot
	for _, s := range slots {
		if s.code == "" {
			continue
		}
		if s.curCode != s.code || s.curQty < utilityMinThreshold {
			needsAttention = append(needsAttention, s)
		}
	}
	if len(needsAttention) == 0 {
		return character, nil
	}

	bankItems, err := d.myBankItems()
	if err != nil {
		return character, err
	}
	bankQty := map[string]int{}
	for _, item := range bankItems {
		bankQty[item.Code] += item.Quantity
	}

	predicted := func(s *utilitySlot) int {
		current := 0
		if s.curCode == s.code {
			current = s.curQty
		}
		return current + bankQty[s.code]
	}
	err = checkDepleted(utility1, utility2, slots, predicted)
	if err != nil {
		return character, err
	}

	var toFill []*utilitySlot
	for _, s := range needsAttention {
		if bankQty[s.code] <= 0 {
			console.Printf("  No %s available in bank to fill %s\n", s.code, s.slot)
			continue
		}
		toFill = append(toFill, s)
	}
	if len(toFill) == 0 {
		return character, nil
	}

	character, err = Move(character, "bank")
	if err != nil {
		return character, err
	}

	for _, s := range toFill {
		if s.curCode != "" && s.curCode != s.code {
			qty := s.curQty
			unequipData, err := d.myActionUnequip(character.Name,
				[]schemas.UnequipSchema{{Slot: s.slot, Quantity: &qty}})
			if err != nil {
				return character, err
			}
			character = unequipData.Character
			s.curCode = ""
			s.curQty = 0
		}

		remaining := utilityMaxStack - s.curQty
		filled := 0
		for remaining > 0 && bankQty[s.code] > 0 {
			freeSpace := character.InventoryMaxItems - totalItems(character)
			if freeSpace <= 0 {
				console.Printf("  No inventory space to withdraw more %s for %s\n",
					s.code, s.slot)
				break
			}
			chunk := remaining
			if bankQty[s.code] < chunk {
				chunk = bankQty[s.code]
			}
			if freeSpace < chunk {
				chunk = freeSpace
			}

			withdrawData, err := d.myActionBankWithdrawItem(character.Name,
				[]schemas.SimpleItemSchema{{Code: s.code, Quantity: chunk}})
			if err != nil {
				return character, err
			}
			character = withdrawData.Character
			bankQty[s.code] -= chunk

			equipData, err := d.myActionEquip(character.Name, []schemas.EquipSchema{
				{Code: s.code, Quantity: &chunk, Slot: s.slot}})
			if err != nil {
				return character, err
			}
			character = equipData.Character

			remaining -= chunk
			filled += chunk
		}
		if remaining > 0 {
			console.Printf("  Only %d/%d %s filled for %s\n",
				filled, utilityMaxStack-s.curQty, s.code, s.slot)
		}
	}

	return character, nil
}

func checkDepleted(utility1, utility2 string, slots []*utilitySlot, qty func(*utilitySlot) int) error {
	if utility1 != "" && utility2 != "" && utility1 == utility2 {
		if qty(slots[0])+qty(slots[1]) == 0 {
			return fmt.Errorf("utility %s depleted", utility1)
		}
		return nil
	}
	for _, s := range slots {
		if s.code == "" {
			continue
		}
		if qty(s) == 0 {
			return fmt.Errorf("utility %s depleted in %s", s.code, s.slot)
		}
	}
	return nil
}
