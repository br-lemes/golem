package routine

import (
	"slices"

	"github.com/br-lemes/golem/pkg/schemas"
)

type DepositOptions struct {
	Movement  MoveOptions
	KeepTypes []string
}

// Deposit moves a character to the bank and stores its unkept items and gold.
func Deposit(character schemas.CharacterSchema, options DepositOptions) (schemas.CharacterSchema, error) {
	// +gocover:ignore:block production wrapper over tested implementation
	return deposit(defaultDeps, character, options)
}

func deposit(d deps, character schemas.CharacterSchema, options DepositOptions) (schemas.CharacterSchema, error) {
	var err error
	character, err = move(d, character, "bank", options.Movement)
	if err != nil {
		return character, err
	}
	items := GetInventoryItems(character, InventoryItemsOptions{
		KeepTypes:         options.KeepTypes,
		KeepTravelPotions: !options.Movement.NoTeleport,
	})
	if len(items) > 0 {
		result, err := d.myActionBankDepositItem(character.Name, items)
		if err != nil {
			return character, err
		}
		character = result.Character
	}
	character, err = restockTravelPotion(d, character, options.Movement)
	if err != nil {
		return character, err
	}
	if character.Gold > 0 && !slices.Contains(options.KeepTypes, "gold") {
		result, err := d.myActionBankDepositGold(character.Name, character.Gold)
		if err != nil {
			return character, err
		}
		character = result.Character
	}
	return character, nil
}
