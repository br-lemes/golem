package cmd

import "github.com/br-lemes/golem/pkg/console"

var characterMap = map[string][]string{
	"br_lemes": {"fighting", "woodcutting"},
	"fb_lemes": {"weaponcrafting", "woodcutting"},
	"bf_lemes": {"gearcrafting", "mining"},
	"mr_lemes": {"cooking", "fishing"},
	"kr_lemes": {"alchemy", "jewelrycrafting"},
}

func getCharacters() []string {
	var characters []string
	for character := range characterMap {
		characters = append(characters, character)
	}
	return characters
}

func getSkills(character string) []string {
	return characterMap[character]
}

func confirmSkill(character string, skill string) bool {
	skills := getSkills(character)
	for _, s := range skills {
		if s == skill {
			return true
		}
	}
	console.Printf("Skill %s is not configured for %s\n", skill, character)
	console.Printf("Available skills: %v\n", skills)
	return console.Confirm("Do you want to continue?")
}
