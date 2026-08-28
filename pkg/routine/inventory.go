package routine

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Inventory(character schemas.CharacterSchema, keepTypes []string) (schemas.CharacterSchema, error) {
	totalItems := totalItems(character)
	if totalItems+5 < character.InventoryMaxItems {
		return character, nil
	}
	character, err := Move(character, "bank")
	if err != nil {
		return character, err
	}
	items := GetInventoryItems(character, keepTypes)
	if len(items) > 0 {
		transaction, err := api.MyActionBankDepositItem(character.Name, items)
		if err != nil {
			return character, err
		}
		character = transaction.Character
	}
	return character, nil
}

func totalItems(character schemas.CharacterSchema) int {
	if character.Inventory == nil {
		return 0
	}
	total := 0
	for _, item := range *character.Inventory {
		total += item.Quantity
	}
	return total
}

func GetInventoryItems(character schemas.CharacterSchema, keepTypes []string) []schemas.SimpleItemSchema {
	items := []schemas.SimpleItemSchema{}
	if character.Inventory == nil {
		return items
	}
	for _, item := range *character.Inventory {
		if item.Code == "" || item.Quantity == 0 {
			continue
		}
		if shouldKeepItem(item.Code, keepTypes) {
			continue
		}
		items = append(items, schemas.SimpleItemSchema{
			Code:     item.Code,
			Quantity: item.Quantity,
		})
	}
	return items
}

func shouldKeepItem(code string, keepTypes []string) bool {
	if len(keepTypes) == 0 {
		return false
	}
	item, found := database.Items().Get(code)
	if !found {
		return false
	}
	for _, keepType := range keepTypes {
		if item.Type == keepType || item.Subtype == keepType {
			return true
		}
	}
	return false
}
