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
	if flags.Loadout {
		t.Fatal("loadout flag should default to false")
	}
	if flags.AllowUnsafe {
		t.Fatal("allow-unsafe flag should default to false")
	}
	err = command.Flags().Set("full", "true")
	if err != nil {
		t.Fatalf("set full flag: %v", err)
	}
	err = command.Flags().Set("allow-unsafe", "true")
	if err != nil {
		t.Fatalf("set allow-unsafe flag: %v", err)
	}
	flags, err = utils.ReadFlags[bestFightFlags](command)
	if err != nil {
		t.Fatalf("read best fight flags: %v", err)
	}
	if !flags.Full {
		t.Fatal("full flag was not read")
	}
	if !flags.AllowUnsafe {
		t.Fatal("allow-unsafe flag was not read")
	}
	err = command.Flags().Set("without", "iron_sword")
	if err != nil {
		t.Fatalf("set without flag: %v", err)
	}
	err = command.Flags().Set("with", "wooden_stick")
	if err != nil {
		t.Fatalf("set with flag: %v", err)
	}
	err = command.Flags().Set("use", "copper_sword")
	if err != nil {
		t.Fatalf("set use flag: %v", err)
	}
	flags, err = utils.ReadFlags[bestFightFlags](command)
	if err != nil {
		t.Fatalf("read item flags: %v", err)
	}
	if len(flags.Without) != 1 || flags.Without[0] != "iron_sword" {
		t.Fatalf("without = %v", flags.Without)
	}
	if len(flags.With) != 1 || flags.With[0] != "wooden_stick" {
		t.Fatalf("with = %v", flags.With)
	}
	if len(flags.Use) != 1 || flags.Use[0] != "copper_sword" {
		t.Fatalf("use = %v", flags.Use)
	}
}

func TestBestFightValidateItems(t *testing.T) {
	valid := bestFightFlags{Without: []string{"iron_sword"}}
	err := bestFightValidateItems(valid)
	if err != nil {
		t.Fatalf("valid without item rejected: %v", err)
	}
	missing := bestFightFlags{Without: []string{"missing_item"}}
	err = bestFightValidateItems(missing)
	if err == nil {
		t.Fatal("missing item accepted")
	}
	conflict := bestFightFlags{
		Without: []string{"iron_sword"},
		Use:     []string{"iron_sword"},
	}
	err = bestFightValidateItems(conflict)
	if err == nil {
		t.Fatal("conflicting item accepted")
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

func TestBestFightLoadout(t *testing.T) {
	loadout := bestFightLoadout(best.Result{
		Level: 12,
		FinalEquipment: map[string]string{
			"weapon":   "iron_sword",
			"helmet":   "iron_helmet",
			"utility1": "",
			"utility2": "health_potion",
		},
	})
	if len(loadout) != 1 {
		t.Fatalf("loadout length = %d", len(loadout))
	}
	character := loadout[0]
	if character.Level != 12 {
		t.Fatalf("level = %d", character.Level)
	}
	if character.WeaponSlot == nil || *character.WeaponSlot != "iron_sword" {
		t.Fatalf("weapon = %v", character.WeaponSlot)
	}
	if character.HelmetSlot == nil || *character.HelmetSlot != "iron_helmet" {
		t.Fatalf("helmet = %v", character.HelmetSlot)
	}
	if character.Utility1Slot != nil {
		t.Fatalf("empty utility was included: %v", character.Utility1Slot)
	}
	if character.Utility2Slot == nil || *character.Utility2Slot != "health_potion" {
		t.Fatalf("utility2 = %v", character.Utility2Slot)
	}
	if character.Utility2SlotQuantity == nil || *character.Utility2SlotQuantity != bestFightUtilityQuantity {
		t.Fatalf("utility2 quantity = %v", character.Utility2SlotQuantity)
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
