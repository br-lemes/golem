package task

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Inventory(character schemas.CharacterSchema, keepFood bool) (schemas.CharacterSchema, error) {
	totalItems := totalItems(character)
	if totalItems+5 < character.InventoryMaxItems {
		return character, nil
	}
	character, err := Move(character, "bank")
	if err != nil {
		return character, err
	}
	transaction, err := api.MyActionBankDepositItem(character.Name,
		GetInventoryItems(character, keepFood))
	if err != nil {
		return character, err
	}
	character = transaction.Character
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

func GetInventoryItems(character schemas.CharacterSchema, keepFood bool) []schemas.SimpleItemSchema {
	items := []schemas.SimpleItemSchema{}
	if character.Inventory == nil {
		return items
	}
	for _, item := range *character.Inventory {
		if item.Code == "" || item.Quantity == 0 {
			continue
		}
		if keepFood && isFood(item.Code) {
			continue
		}
		items = append(items, schemas.SimpleItemSchema{
			Code:     item.Code,
			Quantity: item.Quantity,
		})
	}
	return items
}

func isFood(code string) bool {
	item, found := database.GetItem(code)
	if !found {
		return false
	}
	return item.Type == "consumable" && item.Subtype == "food"
}
