package cmd

import (
	"testing"

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
