package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

type MoveOptions struct {
	AllowGold bool
}

func Move(character schemas.CharacterSchema, code string, options MoveOptions) (schemas.CharacterSchema, error) {
	//+gocover:ignore:block production wrapper over tested implementation
	return move(defaultDeps, character, code, options)
}

func move(d deps, character schemas.CharacterSchema, code string, options MoveOptions) (schemas.CharacterSchema, error) {
	results := find(d, character, code, nil)
	var result *Result
	for index := range results {
		if canUseRequirements(results[index].Requirements, options.AllowGold) {
			result = &results[index]
			break
		}
	}
	if result == nil {
		return character, fmt.Errorf("no coordinates found for code %s", code)
	}
	items := missingItems(character, result.Requirements)
	if len(items) > 0 {
		if code == "bank" {
			return character, fmt.Errorf("required item is not available: %s", items[0].Code)
		}
		var err error
		character, err = move(d, character, "bank", MoveOptions{})
		if err != nil {
			return character, err
		}
		transaction, err := d.myActionBankWithdrawItem(character.Name, items)
		if err != nil {
			return character, err
		}
		return move(d, transaction.Character, code, options)
	}
	gold := requiredGold(result.Requirements)
	if gold > character.Gold {
		if !options.AllowGold {
			return character, fmt.Errorf("route requires %d gold; enable gold usage to continue", gold)
		}
		if code == "bank" {
			return character, fmt.Errorf("required gold is not available before reaching the bank")
		}
		var err error
		character, err = move(d, character, "bank", MoveOptions{})
		if err != nil {
			return character, err
		}
		withdraw, err := d.myActionBankWithdrawGold(character.Name, gold-character.Gold)
		if err != nil {
			return character, err
		}
		return move(d, withdraw.Character, code, options)
	}
	for _, transition := range result.Transitions {
		character, err := makeMove(d, character, transition)
		if err != nil {
			return character, err
		}
		transitionData, err := d.myActionTransition(character.Name)
		if err != nil {
			return character, err
		}
		character = transitionData.Character
	}
	if len(result.Transitions) > 0 {
		lastTransition := result.Transitions[len(result.Transitions)-1]
		if lastTransition.X == result.Target.X && lastTransition.Y == result.Target.Y {
			return character, nil
		}
	}
	return makeMove(d, character, result.Target)
}

func canUseRequirements(conditions []schemas.ConditionSchema, allowGold bool) bool {
	for _, condition := range conditions {
		if condition.Operator == "has_item" || condition.Operator == "achievement_unlocked" {
			continue
		}
		if condition.Operator == "cost" && condition.Code != "gold" {
			continue
		}
		if condition.Code == "gold" && condition.Operator == "cost" && allowGold {
			continue
		}
		return false
	}
	return true
}

func requiredGold(conditions []schemas.ConditionSchema) int {
	total := 0
	for _, condition := range conditions {
		if condition.Code == "gold" && condition.Operator == "cost" {
			total += condition.Value
		}
	}
	return total
}

func missingItems(character schemas.CharacterSchema, conditions []schemas.ConditionSchema) []schemas.SimpleItemSchema {
	hasItemQuantities := make(map[string]int)
	costQuantities := make(map[string]int)
	for _, condition := range conditions {
		switch {
		case condition.Operator == "has_item":
			if condition.Value > hasItemQuantities[condition.Code] {
				hasItemQuantities[condition.Code] = condition.Value
			}
		case condition.Operator == "cost" && condition.Code != "gold":
			costQuantities[condition.Code] += condition.Value
		}
	}
	items := make([]schemas.SimpleItemSchema, 0, len(costQuantities)+len(hasItemQuantities))
	for code, hasItemQuantity := range hasItemQuantities {
		hasItemQuantity -= equippedItemQuantity(character, code)
		if hasItemQuantity < 0 {
			hasItemQuantity = 0
		}
		if costQuantities[code] > hasItemQuantity {
			hasItemQuantity = costQuantities[code]
		}
		quantity := hasItemQuantity - inventoryItemQuantity(character, code)
		if quantity > 0 {
			items = append(items, schemas.SimpleItemSchema{
				Code:     code,
				Quantity: quantity,
			})
		}
	}
	for code, quantity := range costQuantities {
		_, handled := hasItemQuantities[code]
		if handled {
			continue
		}
		quantity -= inventoryItemQuantity(character, code)
		if quantity > 0 {
			items = append(items, schemas.SimpleItemSchema{
				Code:     code,
				Quantity: quantity,
			})
		}
	}
	return items
}

func inventoryItemQuantity(character schemas.CharacterSchema, code string) int {
	quantity := 0
	if character.Inventory == nil {
		return quantity
	}
	for _, item := range *character.Inventory {
		if item.Code == code {
			quantity += item.Quantity
		}
	}
	return quantity
}

func equippedItemQuantity(character schemas.CharacterSchema, code string) int {
	return itemQuantity(character, code) - inventoryItemQuantity(character, code)
}

func makeMove(d deps, character schemas.CharacterSchema, target schemas.MapSchema) (schemas.CharacterSchema, error) {
	if target.X == character.X && target.Y == character.Y {
		return character, nil
	}
	moveData, err := d.myActionMove(character.Name, target.X, target.Y)
	if err != nil {
		return character, err
	}
	return moveData.Character, nil
}
