package utils

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

var characters []string

func GetCharacters() []string {
	return characters
}

func GetCharacterSkillLevel(character schemas.CharacterSchema, skill string) int {
	switch skill {
	case "alchemy":
		return character.AlchemyLevel
	case "combat":
		return character.Level
	case "cooking":
		return character.CookingLevel
	case "fishing":
		return character.FishingLevel
	case "gearcrafting":
		return character.GearcraftingLevel
	case "jewelrycrafting":
		return character.JewelrycraftingLevel
	case "mining":
		return character.MiningLevel
	case "weaponcrafting":
		return character.WeaponcraftingLevel
	case "woodcutting":
		return character.WoodcuttingLevel
	default:
		return character.Level
	}
}

func init() {
	characters = cache.GetCharacters()
}
