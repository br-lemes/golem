package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

type WithdrawOptions struct {
	Movement MoveOptions
}

func Withdraw(character schemas.CharacterSchema, items []schemas.SimpleItemSchema, options WithdrawOptions) (schemas.CharacterSchema, error) {
	//+gocover:ignore:block production wrapper over tested implementation
	return withdraw(defaultDeps, character, items, options)
}

func withdraw(d deps, character schemas.CharacterSchema, items []schemas.SimpleItemSchema, options WithdrawOptions) (schemas.CharacterSchema, error) {
	if len(items) == 0 {
		return character, nil
	}
	err := validateWithdrawItems(items)
	if err != nil {
		return character, err
	}
	err = validateWithdrawStock(d, items)
	if err != nil {
		return character, err
	}
	Cooldown(character)
	character, err = move(d, character, "bank", options.Movement)
	if err != nil {
		return character, err
	}
	err = validateWithdrawSpace(character, items)
	if err != nil {
		return character, err
	}
	result, err := d.myActionBankWithdrawItem(character.Name, items)
	if err != nil {
		return character, err
	}
	return result.Character, nil
}

func validateWithdrawItems(items []schemas.SimpleItemSchema) error {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if seen[item.Code] {
			return fmt.Errorf("duplicate item: %s", item.Code)
		}
		seen[item.Code] = true
	}
	return nil
}

func validateWithdrawStock(d deps, items []schemas.SimpleItemSchema) error {
	bankItems, err := d.myBankItems()
	if err != nil {
		return err
	}
	quantities := make(map[string]int, len(bankItems))
	for _, item := range bankItems {
		quantities[item.Code] += item.Quantity
	}
	for _, item := range items {
		if quantities[item.Code] < item.Quantity {
			return fmt.Errorf("insufficient bank items: %s requires %d, but you only have %d", item.Code, item.Quantity, quantities[item.Code])
		}
	}
	return nil
}

func validateWithdrawSpace(character schemas.CharacterSchema, items []schemas.SimpleItemSchema) error {
	required := 0
	for _, item := range items {
		required += item.Quantity
	}
	available := character.InventoryMaxItems - totalItems(character)
	if required > available {
		return fmt.Errorf("not enough inventory space: required %d, available %d", required, available)
	}
	return nil
}
