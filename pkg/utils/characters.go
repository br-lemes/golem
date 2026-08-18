package utils

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

var characters []string

func GetCharacters() []string {
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

func init() {
	characters = cache.GetCharacters()
}
