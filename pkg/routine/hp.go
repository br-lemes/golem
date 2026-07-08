package routine

import (
	"sort"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

type FoodItem struct {
	Code     string
	Heal     int
	Quantity int
}

func Hp(character schemas.CharacterSchema, minHp int) (schemas.CharacterSchema, error) {
	if character.Hp > minHp {
		return character, nil
	}
	if character.Inventory == nil {
		return handleRest(character)
	}
	foods := []FoodItem{}
	for _, slot := range *character.Inventory {
		if slot.Quantity <= 0 {
			continue
		}

		item, found := database.GetItem(slot.Code)
		if !found || item.Type != "consumable" || item.Subtype != "food" ||
			item.Level > character.Level || item.Effects == nil {
			continue
		}

		healValue := 0
		for _, effect := range *item.Effects {
			if effect.Code == "heal" {
				healValue += effect.Value
				break
			}
		}

		if healValue <= 0 {
			continue
		}
		foods = append(foods, FoodItem{
			Code:     slot.Code,
			Heal:     healValue,
			Quantity: slot.Quantity,
		})
	}
	sort.Slice(foods, func(i, j int) bool {
		return foods[i].Heal > foods[j].Heal
	})

	for i, food := range foods {
		neededHp := character.MaxHp - character.Hp
		if neededHp <= 0 {
			break
		}

		qtyToUse := neededHp / food.Heal
		if qtyToUse > food.Quantity {
			qtyToUse = food.Quantity
		}

		if qtyToUse <= 0 {
			continue
		}
		useData, err := api.MyActionUse(character.Name,
			schemas.SimpleItemSchema{Code: food.Code, Quantity: qtyToUse})
		if err != nil {
			return schemas.CharacterSchema{}, err
		}
		character = useData.Character
		foods[i].Quantity = foods[i].Quantity - qtyToUse
	}

	if character.Hp > minHp {
		return character, nil
	}

	for i := len(foods) - 1; i >= 0; i-- {
		food := foods[i]
		if food.Quantity <= 0 {
			continue
		}

		useData, err := api.MyActionUse(character.Name,
			schemas.SimpleItemSchema{Code: food.Code, Quantity: 1})
		if err != nil {
			return schemas.CharacterSchema{}, err
		}
		character = useData.Character
		break
	}

	if character.Hp > minHp {
		return character, nil
	}

	return handleRest(character)
}

func handleRest(character schemas.CharacterSchema) (schemas.CharacterSchema, error) {
	data, err := api.MyActionRest(character.Name)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	return data.Character, nil
}
