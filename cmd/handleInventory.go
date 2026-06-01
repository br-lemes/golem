package cmd

import (
	"fmt"
	"time"
)

const bankX = 4
const bankY = 1

func handleInventory(character CharacterSchema) (CharacterSchema, error) {
	totalItems := _totalItems(character)
	if totalItems+5 < character.InventoryMaxItems {
		return character, nil
	}
	initialX := character.X
	initialY := character.Y
	fmt.Fprintf(writer, "[%s] Inventory full (%d/%d)...\n",
		time.Now().Format("15:04:05"), totalItems, character.InventoryMaxItems)

	fmt.Fprintf(writer, "[%s] Moving to bank (%d, %d)...\n",
		time.Now().Format("15:04:05"), bankX, bankY)
	moveData, err := apiActionMove(character.Name, bankX, bankY)
	if err != nil {
		return CharacterSchema{}, err
	}
	character = moveData.Character

	fmt.Fprintf(writer, "[%s] Depositing items...\n",
		time.Now().Format("15:04:05"))
	transactionData, err := apiActionBankDepositItem(character.Name, _getItems(character))
	if err != nil {
		return CharacterSchema{}, err
	}
	character = transactionData.Character

	fmt.Fprintf(writer, "[%s] Moving back to initial position (%d, %d)...\n",
		time.Now().Format("15:04:05"), initialX, initialY)
	moveData, err = apiActionMove(character.Name, initialX, initialY)
	if err != nil {
		return CharacterSchema{}, err
	}
	character = moveData.Character

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

func _getItems(character CharacterSchema) []SimpleItemSchema {
	items := []SimpleItemSchema{}
	if character.Inventory == nil {
		return items
	}
	for _, item := range *character.Inventory {
		if item.Code == "" || item.Quantity == 0 {
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
