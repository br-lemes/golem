package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func handleInventory(character schemas.CharacterSchema, keepFood bool) (schemas.CharacterSchema, error) {
	totalItems := _totalItems(character)
	if totalItems+5 < character.InventoryMaxItems {
		return character, nil
	}
	bank := database.FindClosest(character, "bank")
	moveData, err := api.MyActionMove(character.Name, bank.X, bank.Y)
	if err != nil {
		return character, err
	}
	character = moveData.Character
	transactionData, err := api.MyActionBankDepositItem(character.Name, _getItems(character, keepFood))
	if err != nil {
		return character, err
	}
	character = transactionData.Character

	return character, nil
}

func _totalItems(character schemas.CharacterSchema) int {
	if character.Inventory == nil {
		return 0
	}
	total := 0
	for _, item := range *character.Inventory {
		total += item.Quantity
	}
	return total
}

func _getItems(character schemas.CharacterSchema, keepFood bool) []schemas.SimpleItemSchema {
	items := []schemas.SimpleItemSchema{}
	if character.Inventory == nil {
		return items
	}
	for _, item := range *character.Inventory {
		if item.Code == "" || item.Quantity == 0 {
			continue
		}
		if keepFood && _isFood(item.Code) {
			continue
		}
		items = append(items, schemas.SimpleItemSchema{
			Code:     item.Code,
			Quantity: item.Quantity,
		})
	}
	return items
}

func _isFood(code string) bool {
	item, found := database.GetItem(code)
	if !found {
		return false
	}
	return item.Type == "consumable" && item.Subtype == "food"
}
