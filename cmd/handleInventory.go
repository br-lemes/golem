package cmd

import (
	"fmt"
	"time"

	"github.com/br-lemes/golem/pkg/database"
	. "github.com/br-lemes/golem/pkg/schemas"
)

func handleInventory(character CharacterSchema) (CharacterSchema, error) {
	totalItems := _totalItems(character)
	if totalItems+5 < character.InventoryMaxItems {
		return character, nil
	}
	initialX := character.X
	initialY := character.Y
	bank := database.FindClosest(character, "bank")
	fmt.Fprintf(writer, "[%s] Inventory full (%d/%d)...\n",
		time.Now().Format("15:04:05"), totalItems, character.InventoryMaxItems)

	fmt.Fprintf(writer, "[%s] Moving to bank (%d, %d)...\n",
		time.Now().Format("15:04:05"), bank.X, bank.Y)
	moveData, err := apiActionMove(character.Name, bank.X, bank.Y)
	if err != nil {
		return character, err
	}
	character = moveData.Character

	fmt.Fprintf(writer, "[%s] Depositing items...\n",
		time.Now().Format("15:04:05"))
	transactionData, err := apiActionBankDepositItem(character.Name, _getItems(character))
	if err != nil {
		return character, err
	}
	character = transactionData.Character

	fmt.Fprintf(writer, "[%s] Moving back to initial position (%d, %d)...\n",
		time.Now().Format("15:04:05"), initialX, initialY)
	moveData, err = apiActionMove(character.Name, initialX, initialY)
	if err != nil {
		return character, err
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
