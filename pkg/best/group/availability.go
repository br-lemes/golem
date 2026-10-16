// Package group finds equipment for combat parties.
package group

import "github.com/br-lemes/golem/pkg/schemas"

func AvailableItems(bank []schemas.SimpleItemSchema, party []schemas.CharacterSchema) map[string]int {
	available := make(map[string]int)
	for _, item := range bank {
		addItem(available, item.Code, item.Quantity)
	}
	for _, character := range party {
		for _, item := range inventoryItems(character) {
			addItem(available, item.Code, item.Quantity)
		}
		for _, code := range equippedItems(character) {
			addItem(available, code, 1)
		}
		addItem(available, character.Utility1Slot, character.Utility1SlotQuantity)
		addItem(available, character.Utility2Slot, character.Utility2SlotQuantity)
	}
	return available
}

func inventoryItems(character schemas.CharacterSchema) []schemas.InventorySlotSchema {
	if character.Inventory == nil {
		return nil
	}
	return *character.Inventory
}

func equippedItems(character schemas.CharacterSchema) []string {
	return []string{
		character.AmuletSlot,
		character.Artifact1Slot,
		character.Artifact2Slot,
		character.Artifact3Slot,
		character.BagSlot,
		character.BodyArmorSlot,
		character.BootsSlot,
		character.HelmetSlot,
		character.LegArmorSlot,
		character.Ring1Slot,
		character.Ring2Slot,
		character.RuneSlot,
		character.ShieldSlot,
		character.WeaponSlot,
	}
}

func addItem(items map[string]int, code string, quantity int) {
	if code == "" || quantity < 1 {
		return
	}
	items[code] += quantity
}
