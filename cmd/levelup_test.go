package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSkillLevelupRequirements(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{
			Name:              "high",
			Level:             27,
			AlchemyLevel:      20,
			MiningLevel:       9,
			GearcraftingLevel: 50,
		},
		{
			Name:              "low",
			Level:             11,
			AlchemyLevel:      1,
			MiningLevel:       1,
			GearcraftingLevel: 1,
		},
	}

	result := skillLevelupRequirements(characters)

	if result["low"]["combat"] != "11 -> 20" {
		t.Fatalf("low combat requirement = %q, want %q", result["low"]["combat"], "11 -> 20")
	}
	if result["low"]["alchemy"] != "1 -> 20" {
		t.Fatalf("low alchemy requirement = %q, want %q", result["low"]["alchemy"], "1 -> 20")
	}
	_, exists := result["low"]["mining"]
	if exists {
		t.Fatal("low mining should be absent")
	}
	_, exists = result["low"]["gearcrafting"]
	if exists {
		t.Fatal("gearcrafting should be excluded")
	}
	_, exists = result["high"]
	if exists {
		t.Fatal("high character should have no requirements")
	}
}
