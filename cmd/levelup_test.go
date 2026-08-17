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

	if result["low"]["combat"] != 11 {
		t.Fatalf("low combat level = %d, want 11", result["low"]["combat"])
	}
	if result["low"]["alchemy"] != 1 {
		t.Fatalf("low alchemy level = %d, want 1", result["low"]["alchemy"])
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
