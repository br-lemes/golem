package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

type MoveOptions struct {
	AllowGold  bool
	NoTeleport bool
}

func Move(character schemas.CharacterSchema, code string, options MoveOptions) (schemas.CharacterSchema, error) {
	//+gocover:ignore:block production wrapper over tested implementation
	return move(defaultDeps, character, code, options)
}

func move(d deps, character schemas.CharacterSchema, code string, options MoveOptions) (schemas.CharacterSchema, error) {
	events, alreadyAtTarget := currentTarget(d, character, code)
	if alreadyAtTarget {
		return restockTravelPotion(d, character, options)
	}
	potions, bankPotions, atBank, err := availableTeleportPotions(d, character, options)
	if err != nil {
		return character, err
	}
	results := find(d, character, code, potions, events)
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
		character, err = move(d, character, "bank", options)
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
		character, err = move(d, character, "bank", options)
		if err != nil {
			return character, err
		}
		withdraw, err := d.myActionBankWithdrawGold(character.Name, gold-character.Gold)
		if err != nil {
			return character, err
		}
		return move(d, withdraw.Character, code, options)
	}
	if result.Potion != nil {
		character, err = prepareTeleportPotion(d, character, result.Potion.Code, bankPotions, atBank)
		if err != nil {
			return character, err
		}
		if d.myActionUse == nil {
			return character, fmt.Errorf("cannot use teleport potion %s", result.Potion.Code)
		}
		useData, err := d.myActionUse(character.Name, schemas.SimpleItemSchema{
			Code:     result.Potion.Code,
			Quantity: 1,
		})
		if err != nil {
			return character, err
		}
		character = useData.Character
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
			return restockTravelPotion(d, character, options)
		}
	}
	character, err = makeMove(d, character, result.Target)
	if err != nil {
		return character, err
	}
	return restockTravelPotion(d, character, options)
}

func availableTeleportPotions(d deps, character schemas.CharacterSchema, options MoveOptions) ([]schemas.SimpleItemSchema, map[string]int, bool, error) {
	bankPotions := make(map[string]int)
	if options.NoTeleport {
		return nil, bankPotions, false, nil
	}
	quantities := make(map[string]int)
	var order []string
	addPotion := func(code string, quantity int) {
		if quantity <= 0 || !isTeleportPotion(code) {
			return
		}
		_, exists := quantities[code]
		if !exists {
			order = append(order, code)
		}
		quantities[code] += quantity
	}
	if character.Inventory != nil {
		for _, item := range *character.Inventory {
			addPotion(item.Code, item.Quantity)
		}
	}
	atBank := characterAtBank(character)
	if atBank && d.myBankItems != nil {
		items, err := d.myBankItems()
		if err != nil {
			return nil, nil, false, err
		}
		for _, item := range items {
			if isTeleportPotion(item.Code) && item.Quantity > 0 {
				bankPotions[item.Code] += item.Quantity
				if totalItems(character) < character.InventoryMaxItems {
					addPotion(item.Code, item.Quantity)
				}
			}
		}
	}
	potions := make([]schemas.SimpleItemSchema, 0, len(order))
	for _, code := range order {
		potions = append(potions, schemas.SimpleItemSchema{
			Code:     code,
			Quantity: quantities[code],
		})
	}
	return potions, bankPotions, atBank, nil
}

func isTeleportPotion(code string) bool {
	item, exists := catalog.Items().Get(code)
	if !exists || item.Effects == nil {
		return false
	}
	for _, effect := range *item.Effects {
		if effect.Code == "teleport" {
			return true
		}
	}
	return false
}

func characterAtBank(character schemas.CharacterSchema) bool {
	tile, exists := catalog.Maps.Get(catalog.Point{
		X:     character.X,
		Y:     character.Y,
		Layer: character.Layer,
	})
	return exists && tile.Interactions.Content != nil && tile.Interactions.Content.Type == "bank"
}

func prepareTeleportPotion(d deps, character schemas.CharacterSchema, code string, bankPotions map[string]int, atBank bool) (schemas.CharacterSchema, error) {
	inventoryQuantity := inventoryItemQuantity(character, code)
	required := max(0, 1-inventoryQuantity)
	needed := required
	if atBank && isReservedTravelPotion(code) {
		needed = max(needed, 2-inventoryQuantity)
	}
	if needed <= 0 {
		return character, nil
	}
	freeSpace := max(0, character.InventoryMaxItems-totalItems(character))
	withdrawQuantity := min(needed, bankPotions[code], freeSpace)
	if withdrawQuantity < required {
		return character, fmt.Errorf("teleport potion is not available in inventory or bank: %s", code)
	}
	if withdrawQuantity == 0 {
		return character, nil
	}
	if d.myActionBankWithdrawItem == nil {
		return character, fmt.Errorf("cannot withdraw teleport potion %s", code)
	}
	withdraw, err := d.myActionBankWithdrawItem(character.Name, []schemas.SimpleItemSchema{
		{Code: code, Quantity: withdrawQuantity},
	})
	if err != nil {
		return character, err
	}
	character = withdraw.Character
	if inventoryItemQuantity(character, code) < 1 {
		return character, fmt.Errorf("teleport potion is not available in inventory: %s", code)
	}
	return character, nil
}

func isReservedTravelPotion(code string) bool {
	for _, reserved := range reservedTravelPotionCodes {
		if code == reserved {
			return true
		}
	}
	return false
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

func equippedItemQuantity(character schemas.CharacterSchema, code string) int {
	return itemQuantity(character, code) - inventoryItemQuantity(character, code)
}

func makeMove(d deps, character schemas.CharacterSchema, target schemas.MapSchema) (schemas.CharacterSchema, error) {
	if target.X == character.X && target.Y == character.Y && target.Layer == character.Layer {
		return character, nil
	}
	moveData, err := d.myActionMove(character.Name, target.X, target.Y)
	if err != nil {
		return character, err
	}
	return moveData.Character, nil
}
