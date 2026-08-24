package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

func TestBestFlags(t *testing.T) {
	command := &cobra.Command{}
	err := utils.RegisterFlags[bestFlags](command)
	if err != nil {
		t.Fatalf("register best flags: %v", err)
	}
	flags, err := utils.ReadFlags[bestFlags](command)
	if err != nil {
		t.Fatalf("read default best flags: %v", err)
	}
	if flags.AllowDuplicateAdeptRing {
		t.Fatal("allow duplicate adept ring should default to false")
	}
	err = command.Flags().Set("allow-duplicate-adept-ring", "true")
	if err != nil {
		t.Fatalf("set best flag: %v", err)
	}
	flags, err = utils.ReadFlags[bestFlags](command)
	if err != nil {
		t.Fatalf("read best flags: %v", err)
	}
	if !flags.AllowDuplicateAdeptRing {
		t.Fatal("allow duplicate adept ring flag was not read")
	}
}

func TestBestCraftingValidate(t *testing.T) {
	err := bestCraftingValidate("copper_ring")
	if err != nil {
		t.Fatalf("valid craftable item rejected: %v", err)
	}
	err = bestCraftingValidate("copper_ore")
	if err == nil {
		t.Fatal("non-craftable item accepted")
	}
	err = bestCraftingValidate("missing_item")
	if err == nil {
		t.Fatal("missing item accepted")
	}
}

func TestBestGatheringValidate(t *testing.T) {
	err := bestGatheringValidate("ash_tree")
	if err != nil {
		t.Fatalf("valid resource rejected: %v", err)
	}
	err = bestGatheringValidate("missing_resource")
	if err == nil {
		t.Fatal("missing resource accepted")
	}
}

func TestCraftSkillLevel(t *testing.T) {
	character := schemas.CharacterSchema{
		AlchemyLevel:         1,
		CookingLevel:         2,
		GearcraftingLevel:    3,
		JewelrycraftingLevel: 4,
		MiningLevel:          5,
		WeaponcraftingLevel:  6,
		WoodcuttingLevel:     7,
	}
	levels := map[string]int{
		"alchemy":         1,
		"cooking":         2,
		"gearcrafting":    3,
		"jewelrycrafting": 4,
		"mining":          5,
		"weaponcrafting":  6,
		"woodcutting":     7,
		"unknown":         0,
	}
	for skill, want := range levels {
		got := craftSkillLevel(character, skill)
		if got != want {
			t.Errorf("craftSkillLevel(%q) = %d, want %d", skill, got, want)
		}
	}
}

func TestGatheringSkillLevel(t *testing.T) {
	character := schemas.CharacterSchema{
		AlchemyLevel:     1,
		FishingLevel:     2,
		MiningLevel:      3,
		WoodcuttingLevel: 4,
	}
	levels := map[string]int{
		"alchemy":     1,
		"fishing":     2,
		"mining":      3,
		"woodcutting": 4,
		"unknown":     0,
	}
	for skill, want := range levels {
		got := gatheringSkillLevel(character, skill)
		if got != want {
			t.Errorf("gatheringSkillLevel(%q) = %d, want %d", skill, got, want)
		}
	}
}
