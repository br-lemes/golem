package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

func TestFightItemFlags(t *testing.T) {
	command := &cobra.Command{}
	err := utils.RegisterFlags[fightItemFlags](command)
	if err != nil {
		t.Fatalf("register fight item flags: %v", err)
	}
	err = command.Flags().Set("without", "iron_sword")
	if err != nil {
		t.Fatalf("set without flag: %v", err)
	}
	err = command.Flags().Set("use", "copper_sword")
	if err != nil {
		t.Fatalf("set use flag: %v", err)
	}
	flags, err := utils.ReadFlags[fightItemFlags](command)
	if err != nil {
		t.Fatalf("read fight item flags: %v", err)
	}
	if len(flags.Without) != 1 || flags.Without[0] != "iron_sword" {
		t.Fatalf("without = %v", flags.Without)
	}
	if len(flags.Use) != 1 || flags.Use[0] != "copper_sword" {
		t.Fatalf("use = %v", flags.Use)
	}
}

func TestValidateFightItems(t *testing.T) {
	bank := map[string]int{"copper_dagger": 1}
	valid := fightItemFlags{
		Without: []string{"iron_sword"},
		Use:     []string{"copper_dagger"},
	}
	err := validateFightItems(valid, bank)
	if err != nil {
		t.Fatalf("valid items rejected: %v", err)
	}
	missing := fightItemFlags{Without: []string{"missing_item"}}
	err = validateFightItems(missing, bank)
	if err == nil {
		t.Fatal("missing item accepted")
	}
	conflict := fightItemFlags{
		Without: []string{"copper_dagger"},
		Use:     []string{"copper_dagger"},
	}
	err = validateFightItems(conflict, bank)
	if err == nil {
		t.Fatal("conflicting item accepted")
	}
	unavailable := fightItemFlags{Use: []string{"iron_sword"}}
	err = validateFightItems(unavailable, bank)
	if err == nil {
		t.Fatal("item missing from bank accepted")
	}
	insufficient := fightItemFlags{
		Use: []string{"copper_dagger", "copper_dagger"},
	}
	err = validateFightItems(insufficient, bank)
	if err == nil {
		t.Fatal("insufficient bank quantity accepted")
	}
}
