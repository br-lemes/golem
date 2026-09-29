package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCostCraftingCharacter(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{Name: "higher", GearcraftingLevel: 30, GearcraftingXp: 10},
		{Name: "later", GearcraftingLevel: 20, GearcraftingXp: 50},
		{Name: "first", GearcraftingLevel: 20, GearcraftingXp: 10},
	}

	character, compatible := costCraftingCharacter(characters, "gearcrafting", 20)
	if !compatible || character.Name != "first" {
		t.Fatalf("costCraftingCharacter() = (%#v, %t), want first compatible character", character, compatible)
	}

	character, compatible = costCraftingCharacter(characters, "gearcrafting", 40)
	if compatible || character.Name != "higher" {
		t.Fatalf("costCraftingCharacter() = (%#v, %t), want highest incompatible character", character, compatible)
	}
}
