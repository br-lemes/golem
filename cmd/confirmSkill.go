package cmd

import "fmt"

var characterMap = map[string][]string{
	"br_lemes": {"fighting"},
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
	fmt.Fprintf(writer, "Skill %s is not configured for %s\n", skill, character)
	fmt.Fprintf(writer, "Available skills: %v\n", skills)
	return confirm("Do you want to continue?")
}
