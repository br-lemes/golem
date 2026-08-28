package utils

import (
	"slices"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

var characters []string

func GetCharacters() []string {
	// +gocover:ignore:block cache accessor
	return characters
}

func GetCharacterSkillLevel(character schemas.CharacterSchema, skill string) (int, bool) {
	switch skill {
	case "alchemy":
		return character.AlchemyLevel, true
	case "combat":
		return character.Level, true
	case "cooking":
		return character.CookingLevel, true
	case "fishing":
		return character.FishingLevel, true
	case "gearcrafting":
		return character.GearcraftingLevel, true
	case "jewelrycrafting":
		return character.JewelrycraftingLevel, true
	case "mining":
		return character.MiningLevel, true
	case "weaponcrafting":
		return character.WeaponcraftingLevel, true
	case "woodcutting":
		return character.WoodcuttingLevel, true
	default:
		return 0, false
	}
}

func GetCharacterGatheringSkillLevel(character schemas.CharacterSchema, skill string) (int, bool) {
	if !slices.Contains(database.Enum("GatheringSkill"), skill) {
		return 0, false
	}
	return GetCharacterSkillLevel(character, skill)
}

func GetCharacterCraftingSkillLevel(character schemas.CharacterSchema, skill string) (int, bool) {
	if !slices.Contains(database.Enum("CraftSkill"), skill) {
		return 0, false
	}
	return GetCharacterSkillLevel(character, skill)
}

func GetCharacterConditionLevel(character schemas.CharacterSchema, code string) (int, bool) {
	switch code {
	case "level":
		return character.Level, true
	case "alchemy_level":
		return character.AlchemyLevel, true
	case "cooking_level":
		return character.CookingLevel, true
	case "fishing_level":
		return character.FishingLevel, true
	case "gearcrafting_level":
		return character.GearcraftingLevel, true
	case "jewelrycrafting_level":
		return character.JewelrycraftingLevel, true
	case "mining_level":
		return character.MiningLevel, true
	case "weaponcrafting_level":
		return character.WeaponcraftingLevel, true
	case "woodcutting_level":
		return character.WoodcuttingLevel, true
	default:
		return 0, false
	}
}

func MeetsItemConditions(character schemas.CharacterSchema, item schemas.ItemSchema) bool {
	if item.Conditions == nil {
		return true
	}
	for _, condition := range *item.Conditions {
		currentLevel, ok := GetCharacterConditionLevel(character, condition.Code)
		if !ok {
			//+gocover:ignore:block TODO: support achievement_unlocked
			continue
		}
		satisfied := false
		switch condition.Operator {
		case schemas.Gt:
			satisfied = currentLevel > condition.Value
		case schemas.Eq:
			satisfied = currentLevel == condition.Value
		case schemas.Lt:
			satisfied = currentLevel < condition.Value
		case schemas.Ne:
			satisfied = currentLevel != condition.Value
		}
		if !satisfied {
			return false
		}
	}
	return true
}

func init() {
	characters = cache.GetCharacters()
}
