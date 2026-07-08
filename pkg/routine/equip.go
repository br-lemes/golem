package routine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Equip(name string, equipments []string) (schemas.CharacterSchema, error) {
	slotsReserved := make(map[string]string)
	for _, req := range equipments {
		if strings.Contains(req, "@0") {
			return schemas.CharacterSchema{},
				fmt.Errorf("invalid slot number 0")
		}
		itemCode, slotNum, _ := parseEquipmentArg(req)
		item, exists := database.GetItem(itemCode)
		if !exists {
			return schemas.CharacterSchema{},
				fmt.Errorf("item not found in database: %s", itemCode)
		}
		slots, hasSlot := database.EquipmentTypeToSlots[item.Type]
		if !hasSlot {
			return schemas.CharacterSchema{},
				fmt.Errorf("item type cannot be equipped: %s", item.Type)
		}
		if len(slots) == 1 && slotNum != 0 {
			return schemas.CharacterSchema{},
				fmt.Errorf("cannot specify slot number for item with single"+
					" slot: %s", itemCode)
		}
		if len(slots) > 1 && slotNum == 0 {
			return schemas.CharacterSchema{},
				fmt.Errorf("must specify slot (e.g. @1 or @2) for item with"+
					" multiple slots: %s", itemCode)
		}
		if len(slots) > 1 && (slotNum < 1 || slotNum > len(slots)) {
			return schemas.CharacterSchema{},
				fmt.Errorf("invalid slot number %d for item %s", slotNum,
					itemCode)
		}
		targetSlotStr := slots[0]
		if slotNum > 0 && slotNum <= len(slots) {
			targetSlotStr = slots[slotNum-1]
		}
		conflictingItem, reserved := slotsReserved[targetSlotStr]
		if reserved {
			return schemas.CharacterSchema{},
				fmt.Errorf("slot conflict: both %s and %s are targeting %s",
					conflictingItem, itemCode, targetSlotStr)
		}
		slotsReserved[targetSlotStr] = itemCode
	}

	character, err := api.Characters(name)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	Cooldown(character)

	for _, req := range equipments {
		itemCode, _, _ := parseEquipmentArg(req)
		item, _ := database.GetItem(itemCode)
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
			return schemas.CharacterSchema{},
				fmt.Errorf("character %s does not meet requirement for %s"+
					" (has %d, requires %d)", character.Name, item.Name,
					currentLevel, item.Level)
		}
	}
	currentSlots := map[string]string{
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
	currentQuantities := map[string]int{
		"utility1": character.Utility1SlotQuantity,
		"utility2": character.Utility2SlotQuantity,
	}
	var neededEquipments []string
	for _, req := range equipments {
		itemCode, slotNum, qty := parseEquipmentArg(req)
		item, _ := database.GetItem(itemCode)
		slots := database.EquipmentTypeToSlots[item.Type]
		targetSlot := slots[0]
		if slotNum > 0 && slotNum <= len(slots) {
			targetSlot = slots[slotNum-1]
		}
		currentItem := currentSlots[targetSlot]
		currentQty := 1
		if q, ok := currentQuantities[targetSlot]; ok {
			currentQty = q
		}
		if currentItem == itemCode && currentQty == qty {
			continue
		}
		neededEquipments = append(neededEquipments, req)
	}

	if len(neededEquipments) == 0 {
		return character, nil
	}

	bankItems, err := api.MyBankItems()
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	for _, req := range neededEquipments {
		itemCode, _, qty := parseEquipmentArg(req)
		totalAvailable := 0
		if character.Inventory != nil {
			for _, invSlot := range *character.Inventory {
				if invSlot.Code == itemCode {
					totalAvailable += invSlot.Quantity
				}
			}
		}
		for _, bankItem := range bankItems {
			if bankItem.Code == itemCode {
				totalAvailable += bankItem.Quantity
			}
		}
		if totalAvailable < qty {
			return schemas.CharacterSchema{},
				fmt.Errorf("insufficient items: %s requires %d, but you only"+
					" have %d (inventory+bank)", itemCode, qty, totalAvailable)
		}
	}
	var withdrawPayload []schemas.SimpleItemSchema
	withdrawMap := make(map[string]int)
	for _, req := range neededEquipments {
		itemCode, _, qty := parseEquipmentArg(req)
		qtyInInventory := 0
		if character.Inventory != nil {
			for _, invSlot := range *character.Inventory {
				if invSlot.Code == itemCode {
					qtyInInventory = invSlot.Quantity
					break
				}
			}
		}
		if qtyInInventory < qty {
			neededQty := qty - qtyInInventory
			withdrawMap[itemCode] += neededQty
		}
	}
	for code, qty := range withdrawMap {
		withdrawPayload = append(withdrawPayload, schemas.SimpleItemSchema{
			Code:     code,
			Quantity: qty,
		})
	}
	if len(withdrawPayload) > 0 {
		character, err = Move(character, "bank")
		if err != nil {
			return character, err
		}
		transaction, err := api.MyActionBankWithdrawItem(name, withdrawPayload)
		if err != nil {
			return character, err
		}
		character = transaction.Character
	}

	var equipPayload []schemas.EquipSchema
	for _, req := range neededEquipments {
		itemCode, slotNum, qty := parseEquipmentArg(req)
		item, _ := database.GetItem(itemCode)
		slots := database.EquipmentTypeToSlots[item.Type]
		targetSlotStr := slots[0]
		if slotNum > 0 && slotNum <= len(slots) {
			targetSlotStr = slots[slotNum-1]
		}
		targetSlot := schemas.ItemSlot(targetSlotStr)
		apiQty := qty
		equipPayload = append(equipPayload, schemas.EquipSchema{
			Code:     itemCode,
			Slot:     targetSlot,
			Quantity: &apiQty,
		})
	}
	if len(equipPayload) > 0 {
		result, err := api.MyActionEquip(name, equipPayload)
		if err != nil {
			return character, err
		}
		character = result.Character
	}

	return character, nil
}

func parseEquipmentArg(arg string) (string, int, int) {
	qty := 1
	slot := 0
	itemPart := arg
	if strings.Contains(arg, "@") {
		parts := strings.Split(arg, "@")
		itemPart = parts[0]
		rightPart := parts[1]
		if strings.Contains(rightPart, "x") {
			subParts := strings.Split(rightPart, "x")
			slot, _ = strconv.Atoi(subParts[0])
			qty, _ = strconv.Atoi(subParts[1])
		} else {
			slot, _ = strconv.Atoi(rightPart)
		}
	}
	return itemPart, slot, qty
}
