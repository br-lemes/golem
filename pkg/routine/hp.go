package routine

import (
	"math"
	"sort"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

type FoodItem struct {
	Code     string
	Heal     int
	Quantity int
}

type HpOptions struct {
	MinHP    int
	UseFood  bool
	FoodOnly string
}

func Hp(character schemas.CharacterSchema, options HpOptions) (schemas.CharacterSchema, error) {
	return hp(defaultDeps, character, options)
}

func hp(d deps, character schemas.CharacterSchema, hpOptions HpOptions) (schemas.CharacterSchema, error) {
	minHp := hpOptions.MinHP
	if character.Hp >= minHp {
		return character, nil
	}
	neededHp := character.MaxHp - character.Hp
	missingHpPercent := (float64(neededHp) / float64(character.MaxHp)) * 100
	estimatedRestCooldown := math.Ceil(missingHpPercent)
	if estimatedRestCooldown < 3 {
		estimatedRestCooldown = 3
	}
	if !hpOptions.UseFood || character.Inventory == nil || estimatedRestCooldown <= 3 {
		return rest(d, character)
	}
	foods := []FoodItem{}
	for _, slot := range *character.Inventory {
		if slot.Quantity <= 0 { //+gocover:ignore:block should not happen
			continue
		}
		item, found := database.Items().Foods().Get(slot.Code)
		if hpOptions.FoodOnly != "" && slot.Code != hpOptions.FoodOnly {
			continue
		}
		if !found || !utils.MeetsItemConditions(character, *item) || item.Effects == nil {
			//+gocover:ignore:block should not happen
			continue
		}
		healValue := 0
		for _, effect := range *item.Effects {
			if effect.Code == "heal" {
				healValue += effect.Value
				break
			}
		}
		if healValue <= 0 { //+gocover:ignore:block should not happen
			continue
		}
		foods = append(foods, FoodItem{
			Code:     slot.Code,
			Heal:     healValue,
			Quantity: slot.Quantity,
		})
	}
	if len(foods) == 0 {
		return rest(d, character)
	}
	sort.Slice(foods, func(i, j int) bool {
		return foods[i].Heal > foods[j].Heal
	})
	for _, food := range foods {
		neededHp = character.MaxHp - character.Hp
		if neededHp <= 0 {
			break
		}
		qtyToUse := (neededHp + food.Heal - 1) / food.Heal
		if qtyToUse > food.Quantity {
			qtyToUse = food.Quantity
		}
		if qtyToUse <= 0 { //+gocover:ignore:block should not happen
			continue
		}
		useData, err := d.myActionUse(character.Name, schemas.SimpleItemSchema{
			Code:     food.Code,
			Quantity: qtyToUse,
		})
		if err != nil {
			return schemas.CharacterSchema{}, err
		}
		character = useData.Character
	}
	if character.Hp < character.MaxHp {
		return rest(d, character)
	}
	return character, nil
}

func rest(d deps, character schemas.CharacterSchema) (schemas.CharacterSchema, error) {
	data, err := d.myActionRest(character.Name)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	return data.Character, nil
}
