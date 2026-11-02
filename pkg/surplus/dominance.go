package surplus

import (
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

var lowerIsBetter = map[string]bool{
	"alchemy_crafting":      true,
	"alchemy_gathering":     true,
	"cooking_crafting":      true,
	"fishing_gathering":     true,
	"mining_crafting":       true,
	"mining_gathering":      true,
	"woodcutting_crafting":  true,
	"woodcutting_gathering": true,
}

var elementalDamageEffects = []string{
	"dmg_fire",
	"dmg_water",
	"dmg_earth",
	"dmg_air",
}

func Dominates(superior, inferior schemas.ItemSchema) bool {
	return dominatesEffects(effectValues(superior), effectValues(inferior))
}

func dominatesEffects(superiorEffects, inferiorEffects map[string]int) bool {
	strictlyBetter := false
	for code, inferiorValue := range inferiorEffects {
		superiorValue := superiorEffects[code]
		if !atLeastAsGood(code, superiorValue, inferiorValue) {
			return false
		}
		if superiorValue != inferiorValue {
			strictlyBetter = true
		}
	}
	for code, superiorValue := range superiorEffects {
		if inferiorEffects[code] == superiorValue {
			continue
		}
		if !atLeastAsGood(code, superiorValue, inferiorEffects[code]) {
			return false
		}
		strictlyBetter = true
	}
	return strictlyBetter
}

func CanReplace(superior, inferior schemas.ItemSchema, character schemas.CharacterSchema) bool {
	if !compatibleItems(superior, inferior) {
		return false
	}
	if !CanUse(superior, character) || !CanUse(inferior, character) {
		return false
	}
	if superior.Subtype == "tool" {
		return dominatesEffects(toolEffectValues(superior), toolEffectValues(inferior))
	}
	return Dominates(superior, inferior)
}

func DominatesEquipment(superior, inferior schemas.ItemSchema) bool {
	if !compatibleItems(superior, inferior) {
		return false
	}
	if superior.Subtype == "tool" {
		return dominatesEffects(toolEffectValues(superior), toolEffectValues(inferior))
	}
	return Dominates(superior, inferior)
}

func NonDominated(items []schemas.ItemSchema, character schemas.CharacterSchema) []schemas.ItemSchema {
	result := make([]schemas.ItemSchema, 0, len(items))
	for index, item := range items {
		dominated := false
		for otherIndex, other := range items {
			if index == otherIndex {
				continue
			}
			if CanReplace(other, item, character) {
				dominated = true
				break
			}
		}
		if !dominated {
			result = append(result, item)
		}
	}
	return result
}

func compatibleItems(first, second schemas.ItemSchema) bool {
	if first.Subtype == "tool" || second.Subtype == "tool" {
		return first.Subtype == "tool" && second.Subtype == "tool" && toolSkill(first) != "" && toolSkill(first) == toolSkill(second)
	}
	return first.Type == second.Type && first.Subtype == second.Subtype
}

func toolEffectValues(item schemas.ItemSchema) map[string]int {
	values := map[string]int{}
	if item.Effects == nil {
		//+gocover:ignore:block cannot occur in real catalog
		return values
	}
	for _, effect := range *item.Effects {
		if lowerIsBetter[effect.Code] {
			values[effect.Code] += effect.Value
		}
	}
	return values
}

func CanUse(item schemas.ItemSchema, character schemas.CharacterSchema) bool {
	if item.Subtype == "tool" {
		skill := toolSkill(item)
		level, ok := utils.GetCharacterSkillLevel(character, skill)
		return ok && level >= item.Level
	}
	return character.Level >= item.Level
}

func effectValues(item schemas.ItemSchema) map[string]int {
	values := map[string]int{}
	if item.Effects == nil {
		//+gocover:ignore:block cannot occur in real catalog
		return values
	}
	for _, effect := range *item.Effects {
		if effect.Code == "dmg" {
			for _, elementalEffect := range elementalDamageEffects {
				values[elementalEffect] += effect.Value
			}
			continue
		}
		values[effect.Code] += effect.Value
	}
	return values
}

func atLeastAsGood(code string, superior, inferior int) bool {
	if lowerIsBetter[code] {
		return superior <= inferior
	}
	return superior >= inferior
}
