package cmd

import (
	"github.com/br-lemes/golem/pkg/database"
	. "github.com/br-lemes/golem/pkg/schemas"
)

func handleInventory(character CharacterSchema, keepFood bool) (CharacterSchema, error) {
	totalItems := _totalItems(character)
	if totalItems+5 < character.InventoryMaxItems {
		return character, nil
	}
	bank := database.FindClosest(character, "bank")
	moveData, err := apiActionMove(character.Name, bank.X, bank.Y)
	if err != nil {
		return character, err
	}
	character = moveData.Character
	transactionData, err := apiActionBankDepositItem(character.Name, _getItems(character, keepFood))
	if err != nil {
		return character, err
	}
	character = transactionData.Character

	return character, nil
}

func _totalItems(character CharacterSchema) int {
	if character.Inventory == nil {
		return 0
	}
	total := 0
	for _, item := range *character.Inventory {
		total += item.Quantity
	}
	return total
}

func _getItems(character CharacterSchema, keepFood bool) []SimpleItemSchema {
	items := []SimpleItemSchema{}
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
		simpleItem := SimpleItemSchema{
			Code:     item.Code,
			Quantity: item.Quantity,
		}
		items = append(items, simpleItem)
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
