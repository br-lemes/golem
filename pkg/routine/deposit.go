package routine

import (
	"slices"

	"github.com/br-lemes/golem/pkg/schemas"
)

// keepTypes preserves matching item types, subtypes, and "gold".
func Deposit(character schemas.CharacterSchema, keepTypes []string) (schemas.CharacterSchema, error) {
	// +gocover:ignore:block production wrapper over tested implementation
	return deposit(defaultDeps, character, keepTypes)
}

func deposit(d deps, character schemas.CharacterSchema, keepTypes []string) (schemas.CharacterSchema, error) {
	var err error
	character, err = move(d, character, "bank")
	if err != nil {
		return character, err
	}
	items := GetInventoryItems(character, keepTypes)
	if len(items) > 0 {
		result, err := d.myActionBankDepositItem(character.Name, items)
		if err != nil {
			return character, err
		}
		character = result.Character
	}
	if character.Gold > 0 && !slices.Contains(keepTypes, "gold") {
		result, err := d.myActionBankDepositGold(character.Name, character.Gold)
		if err != nil {
			return character, err
		}
		character = result.Character
	}
	return character, nil
}
