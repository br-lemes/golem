package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

const utilityMaxStack = 100
const utilityMinThreshold = 10

type utilitySlot struct {
	slot    schemas.ItemSlot
	code    string
	curCode string
	curQty  int
}

func ClearUtilities(character schemas.CharacterSchema, slots []schemas.ItemSlot) (schemas.CharacterSchema, error) {
	return clearUtilities(defaultDeps, character, slots)
}

func clearUtilities(d deps, character schemas.CharacterSchema, slots []schemas.ItemSlot) (schemas.CharacterSchema, error) {
	unequips := make([]schemas.UnequipSchema, 0, 2)
	for _, slot := range slots {
		code := character.Utility1Slot
		quantity := character.Utility1SlotQuantity
		if slot == schemas.Utility2 {
			code = character.Utility2Slot
			quantity = character.Utility2SlotQuantity
		}
		if code == "" {
			continue
		}
		unequips = append(unequips, schemas.UnequipSchema{
			Slot:     slot,
			Quantity: &quantity,
		})
	}
	if len(unequips) == 0 {
		return character, nil
	}
	result, err := d.myActionUnequip(character.Name, unequips)
	if err != nil {
		return character, err
	}
	return result.Character, nil
}

func utilitySlots(character schemas.CharacterSchema, utility1, utility2 string) []*utilitySlot {
	return []*utilitySlot{
		{
			schemas.Utility1,
			utility1,
			character.Utility1Slot,
			character.Utility1SlotQuantity,
		},
		{
			schemas.Utility2,
			utility2,
			character.Utility2Slot,
			character.Utility2SlotQuantity,
		},
	}
}

func utilityCheck(character schemas.CharacterSchema, utility1, utility2 string, bankQty map[string]int) (bool, error) {
	slots := utilitySlots(character, utility1, utility2)
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
		return false, nil
	}
	predicted := func(s *utilitySlot) int {
		current := 0
		if s.curCode == s.code {
			current = s.curQty
		}
		return current + bankQty[s.code]
	}
	err := checkDepleted(utility1, utility2, slots, predicted)
	if err != nil {
		return false, err
	}
	for _, s := range needsAttention {
		if bankQty[s.code] > 0 {
			return true, nil
		}
		console.Printf("  No %s available in bank to fill %s\n", s.code, s.slot)
	}
	return false, nil
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

func utilityRestock(d deps, character schemas.CharacterSchema, utility1, utility2 string, bankQty map[string]int) (schemas.CharacterSchema, error) {
	slots := utilitySlots(character, utility1, utility2)
	for _, s := range slots {
		if s.code == "" || bankQty[s.code] <= 0 {
			continue
		}
		if s.curCode == s.code && s.curQty >= utilityMaxStack {
			continue
		}
		if s.curCode != "" && s.curCode != s.code {
			qty := s.curQty
			unequipData, err := d.myActionUnequip(character.Name, []schemas.UnequipSchema{
				{Slot: s.slot, Quantity: &qty},
			})
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
				console.Printf("No inventory space to withdraw more %s for %s\n", s.code, s.slot)
				break
			}
			chunk := remaining
			if bankQty[s.code] < chunk {
				chunk = bankQty[s.code]
			}
			if freeSpace < chunk {
				chunk = freeSpace
			}
			withdrawData, err := d.myActionBankWithdrawItem(character.Name, []schemas.SimpleItemSchema{
				{Code: s.code, Quantity: chunk},
			})
			if err != nil {
				return character, err
			}
			character = withdrawData.Character
			bankQty[s.code] -= chunk
			equipData, err := d.myActionEquip(character.Name, []schemas.EquipSchema{
				{Code: s.code, Quantity: &chunk, Slot: s.slot},
			})
			if err != nil {
				return character, err
			}
			character = equipData.Character
			remaining -= chunk
			filled += chunk
		}
		if remaining > 0 {
			console.Printf("  Only %d/%d %s filled for %s\n", filled, utilityMaxStack-s.curQty, s.code, s.slot)
		}
	}
	return character, nil
}
