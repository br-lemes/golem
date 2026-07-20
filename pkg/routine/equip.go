package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Equip(name string, equipments []schemas.EquipSchema) (schemas.CharacterSchema, error) {
	return equip(defaultDeps, name, equipments)
}

func equip(d deps, name string, equipments []schemas.EquipSchema) (schemas.CharacterSchema, error) {
	err := validateEquipments(equipments)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	character, err := d.characters(name)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	err = checkLevelRequirements(character, equipments)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	needed, err := filterNeededEquipments(character, equipments)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	if len(needed) == 0 {
		return character, nil
	}
	err = validateTotalStock(d, character, needed)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}

	Cooldown(character)

	missing := calculateMissingItems(character, needed)
	if len(missing) > 0 {
		character, err = move(d, character, "bank")
		if err != nil {
			return character, err
		}
		transaction, err := d.myActionBankWithdrawItem(name, missing)
		if err != nil {
			return character, err
		}
		character = transaction.Character
	}
	result, err := d.myActionEquip(name, needed)
	if err != nil {
		return character, err
	}
	character = result.Character
	return character, nil
}

func validateEquipments(equipments []schemas.EquipSchema) error {
	slotsReserved := map[schemas.ItemSlot]string{}
	for _, equipment := range equipments {
		if equipment.Code == "" || equipment.Slot == "" {
			return fmt.Errorf("invalid equipment request: missing code or slot")
		}
		item, exists := database.GetItem(equipment.Code)
		if !exists {
			return fmt.Errorf("item not found in database: %s", equipment.Code)
		}
		slots, hasSlot := database.EquipmentTypeToSlots[item.Type]
		if !hasSlot {
			return fmt.Errorf("item type cannot be equipped: %s", item.Type)
		}
		validSlot := false
		for _, s := range slots {
			if schemas.ItemSlot(s) == equipment.Slot {
				validSlot = true
				break
			}
		}
		if !validSlot {
			return fmt.Errorf("item %s cannot be equipped in slot %s",
				equipment.Code, equipment.Slot)
		}
		conflictingItem, reserved := slotsReserved[equipment.Slot]
		if reserved {
			return fmt.Errorf("slot conflict: both %s and %s are targeting %s",
				conflictingItem, equipment.Code, equipment.Slot)
		}
		slotsReserved[equipment.Slot] = equipment.Code
	}
	return nil
}

func checkLevelRequirements(character schemas.CharacterSchema, equipments []schemas.EquipSchema) error {
	for _, equipment := range equipments {
		item, _ := database.GetItem(equipment.Code)
		currentLevel := character.Level
		if item.Subtype == "tool" && item.Effects != nil {
			for _, effect := range *item.Effects {
				switch effect.Code {
				case "alchemy":
					currentLevel = character.AlchemyLevel
				case "fishing":
					currentLevel = character.FishingLevel
				case "mining":
					currentLevel = character.MiningLevel
				case "woodcutting":
					currentLevel = character.WoodcuttingLevel
				}
			}
		}
		if currentLevel < item.Level {
			return fmt.Errorf(
				"does not meet requirement for %s (has %d, requires %d)",
				item.Name, currentLevel, item.Level)
		}
	}
	return nil
}

func filterNeededEquipments(character schemas.CharacterSchema, equipments []schemas.EquipSchema) ([]schemas.EquipSchema, error) {
	currentSlots := map[schemas.ItemSlot]string{
		"amulet":     character.AmuletSlot,
		"artifact1":  character.Artifact1Slot,
		"artifact2":  character.Artifact2Slot,
		"artifact3":  character.Artifact3Slot,
		"bag":        character.BagSlot,
		"body_armor": character.BodyArmorSlot,
		"boots":      character.BootsSlot,
		"helmet":     character.HelmetSlot,
		"leg_armor":  character.LegArmorSlot,
		"ring1":      character.Ring1Slot,
		"ring2":      character.Ring2Slot,
		"rune":       character.RuneSlot,
		"shield":     character.ShieldSlot,
		"utility1":   character.Utility1Slot,
		"utility2":   character.Utility2Slot,
		"weapon":     character.WeaponSlot,
	}
	currentQuantities := map[schemas.ItemSlot]int{
		"utility1": character.Utility1SlotQuantity,
		"utility2": character.Utility2SlotQuantity,
	}

	var needed []schemas.EquipSchema
	for _, equipment := range equipments {
		currentQty := 1
		currentItem := currentSlots[equipment.Slot]
		q, ok := currentQuantities[equipment.Slot]
		if ok {
			currentQty = q
		}
		targetQty := 1
		if equipment.Quantity != nil {
			targetQty = *equipment.Quantity
		}
		if currentItem == equipment.Code && currentQty == targetQty {
			continue
		}
		needed = append(needed, equipment)
	}
	return needed, nil
}

func validateTotalStock(d deps, character schemas.CharacterSchema, needed []schemas.EquipSchema) error {
	bankItems, err := d.myBankItems()
	if err != nil {
		return err
	}
	totalRequired := map[string]int{}
	for _, equipment := range needed {
		targetQty := 1
		if equipment.Quantity != nil {
			targetQty = *equipment.Quantity
		}
		totalRequired[equipment.Code] += targetQty
	}
	for code, requiredQty := range totalRequired {
		totalAvailable := 0
		if character.Inventory != nil {
			for _, invSlot := range *character.Inventory {
				if invSlot.Code == code {
					totalAvailable += invSlot.Quantity
				}
			}
		}
		for _, bankItem := range bankItems {
			if bankItem.Code == code {
				totalAvailable += bankItem.Quantity
			}
		}
		if totalAvailable < requiredQty {
			return fmt.Errorf(
				"insufficient items: %s requires %d, but you only have %d",
				code, requiredQty, totalAvailable)
		}
	}
	return nil
}

func calculateMissingItems(character schemas.CharacterSchema, needed []schemas.EquipSchema) []schemas.SimpleItemSchema {
	totalRequired := map[string]int{}
	for _, equipment := range needed {
		targetQty := 1
		if equipment.Quantity != nil {
			targetQty = *equipment.Quantity
		}
		totalRequired[equipment.Code] += targetQty
	}
	var missing []schemas.SimpleItemSchema
	for code, requiredQty := range totalRequired {
		qtyInInventory := 0
		if character.Inventory != nil {
			for _, invSlot := range *character.Inventory {
				if invSlot.Code == code {
					qtyInInventory = invSlot.Quantity
					break
				}
			}
		}
		if qtyInInventory < requiredQty {
			missing = append(missing, schemas.SimpleItemSchema{
				Code:     code,
				Quantity: requiredQty - qtyInInventory,
			})
		}
	}
	return missing
}
