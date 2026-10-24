package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/best"
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

func TestBestFightFlags(t *testing.T) {
	command := &cobra.Command{}
	err := utils.RegisterFlags[bestFightFlags](command)
	if err != nil {
		t.Fatalf("register best fight flags: %v", err)
	}
	flags, err := utils.ReadFlags[bestFightFlags](command)
	if err != nil {
		t.Fatalf("read default best fight flags: %v", err)
	}
	if flags.Full {
		t.Fatal("full flag should default to false")
	}
	err = command.Flags().Set("full", "true")
	if err != nil {
		t.Fatalf("set full flag: %v", err)
	}
	flags, err = utils.ReadFlags[bestFightFlags](command)
	if err != nil {
		t.Fatalf("read best fight flags: %v", err)
	}
	if !flags.Full {
		t.Fatal("full flag was not read")
	}
}

func TestBestFightFullResult(t *testing.T) {
	result := bestFightFullResult(best.Result{
		Equipment: map[string]string{"weapon": "old_sword"},
		Utilities: map[string]string{"utility1": "old_potion"},
		FinalEquipment: map[string]string{
			"weapon":   "iron_sword",
			"helmet":   "iron_helmet",
			"utility1": "health_potion",
			"utility2": "",
		},
	})
	if result.Equipment["weapon"] != "iron_sword" || result.Equipment["helmet"] != "iron_helmet" {
		t.Fatalf("full equipment = %#v", result.Equipment)
	}
	_, utilityIncluded := result.Equipment["utility1"]
	if utilityIncluded {
		t.Fatalf("utility was included in equipment: %#v", result.Equipment)
	}
	if result.Utilities["utility1"] != "health_potion" || result.Utilities["utility2"] != "" {
		t.Fatalf("full utilities = %#v", result.Utilities)
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
