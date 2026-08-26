package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type simulationFlags struct {
	File                 string `flag:"file" desc:"JSON array file containing one simulation character"`
	Name                 string `flag:"name" desc:"Use a real character as the base"`
	Iterations           int    `flag:"iterations" default:"10" desc:"Number of simulations to run"`
	Level                int    `flag:"level" default:"1" desc:"Fake character level (1-50)"`
	Logs                 bool   `flag:"logs" desc:"Include combat logs in the output"`
	PlayerCritical       int    `flag:"player-critical" default:"-1" desc:"Override player critical chance (0-100)"`
	MonsterCritical      int    `flag:"monster-critical" default:"-1" desc:"Override monster critical chance (0-100)"`
	Utility1SlotQuantity int    `flag:"utility1_slot_quantity" default:"1" desc:"Utility 1 quantity"`
	Utility2SlotQuantity int    `flag:"utility2_slot_quantity" default:"1" desc:"Utility 2 quantity"`
	WeaponSlot           string `flag:"weapon_slot" default:"wooden_stick" desc:"Item code for the weapon slot"`
	RuneSlot             string `flag:"rune_slot" desc:"Item code for the rune slot"`
	ShieldSlot           string `flag:"shield_slot" desc:"Item code for the shield slot"`
	HelmetSlot           string `flag:"helmet_slot" desc:"Item code for the helmet slot"`
	BodyArmorSlot        string `flag:"body_armor_slot" desc:"Item code for the body armor slot"`
	LegArmorSlot         string `flag:"leg_armor_slot" desc:"Item code for the leg armor slot"`
	BootsSlot            string `flag:"boots_slot" desc:"Item code for the boots slot"`
	Ring1Slot            string `flag:"ring1_slot" desc:"Item code for the first ring slot"`
	Ring2Slot            string `flag:"ring2_slot" desc:"Item code for the second ring slot"`
	AmuletSlot           string `flag:"amulet_slot" desc:"Item code for the amulet slot"`
	Artifact1Slot        string `flag:"artifact1_slot" desc:"Item code for the first artifact slot"`
	Artifact2Slot        string `flag:"artifact2_slot" desc:"Item code for the second artifact slot"`
	Artifact3Slot        string `flag:"artifact3_slot" desc:"Item code for the third artifact slot"`
	Utility1Slot         string `flag:"utility1_slot" desc:"Item code for the first utility slot"`
	Utility2Slot         string `flag:"utility2_slot" desc:"Item code for the second utility slot"`
}

var simulationCmd = &cobra.Command{
	Use:   "simulation",
	Short: "Run combat simulations",
}

var simulationSlotNames = []string{
	"weapon",
	"rune",
	"shield",
	"helmet",
	"body_armor",
	"leg_armor",
	"boots",
	"ring1",
	"ring2",
	"amulet",
	"artifact1",
	"artifact2",
	"artifact3",
	"utility1",
	"utility2",
}

type simulationInput struct {
	Monster   schemas.MonsterSchema
	Character schemas.FakeCharacterSchema
	Flags     simulationFlags
}

func readSimulationInput(cmd *cobra.Command, monster string, allowCriticalOverrides bool) (simulationInput, error) {
	flags, err := utils.ReadFlags[simulationFlags](cmd)
	if err != nil {
		return simulationInput{}, err
	}
	err = validateSimulationFlags(flags)
	if err != nil {
		return simulationInput{}, err
	}
	if !allowCriticalOverrides && (cmd.Flags().Changed("player-critical") || cmd.Flags().Changed("monster-critical")) {
		return simulationInput{}, fmt.Errorf("critical overrides are supported by simulation local, compare, and critical, but not by the API")
	}
	monsterData, ok := database.Monsters.Get(monster)
	if !ok {
		return simulationInput{}, fmt.Errorf("invalid monster: %s", monster)
	}

	explicitSlots := map[string]string{}
	slotValues := map[string]string{
		"weapon":     flags.WeaponSlot,
		"rune":       flags.RuneSlot,
		"shield":     flags.ShieldSlot,
		"helmet":     flags.HelmetSlot,
		"body_armor": flags.BodyArmorSlot,
		"leg_armor":  flags.LegArmorSlot,
		"boots":      flags.BootsSlot,
		"ring1":      flags.Ring1Slot,
		"ring2":      flags.Ring2Slot,
		"amulet":     flags.AmuletSlot,
		"artifact1":  flags.Artifact1Slot,
		"artifact2":  flags.Artifact2Slot,
		"artifact3":  flags.Artifact3Slot,
		"utility1":   flags.Utility1Slot,
		"utility2":   flags.Utility2Slot,
	}
	for slot, value := range slotValues {
		if cmd.Flags().Changed(slot + "_slot") {
			explicitSlots[slot] = value
		}
	}
	explicitQuantities := map[string]int{}
	if cmd.Flags().Changed("utility1_slot_quantity") {
		explicitQuantities["utility1"] = flags.Utility1SlotQuantity
	}
	if cmd.Flags().Changed("utility2_slot_quantity") {
		explicitQuantities["utility2"] = flags.Utility2SlotQuantity
	}
	character, err := fight.ResolveCharacter(fight.CharacterOptions{
		File:              flags.File,
		Name:              flags.Name,
		Level:             flags.Level,
		ExplicitSlots:     explicitSlots,
		UtilityQuantities: explicitQuantities,
	})
	if err != nil {
		return simulationInput{}, err
	}
	return simulationInput{
		Monster:   *monsterData,
		Character: character,
		Flags:     flags,
	}, nil
}

func validateSimulationFlags(flags simulationFlags) error {
	if flags.Iterations < 1 {
		return fmt.Errorf("iterations must be at least 1")
	}
	if flags.Level < 1 || flags.Level > 50 {
		return fmt.Errorf("level must be between 1 and 50")
	}
	if flags.Utility1SlotQuantity < 1 || flags.Utility1SlotQuantity > 100 || flags.Utility2SlotQuantity < 1 || flags.Utility2SlotQuantity > 100 {
		return fmt.Errorf("utility quantities must be between 1 and 100")
	}
	if flags.PlayerCritical < -1 || flags.PlayerCritical > 100 || flags.MonsterCritical < -1 || flags.MonsterCritical > 100 {
		return fmt.Errorf("critical overrides must be between 0 and 100, or -1 to preserve the original value")
	}
	return nil
}

func simulationRequest(monster string, character schemas.FakeCharacterSchema, flags simulationFlags) (schemas.CombatSimulationRequestSchema, error) {
	request := schemas.CombatSimulationRequestSchema{
		Characters: []schemas.FakeCharacterSchema{character},
		Monster:    monster,
		Iterations: flags.Iterations,
	}
	return request, fight.ValidateRequest(request)
}

func simulationCriticalOptions(flags simulationFlags) fight.SimulationOptions {
	options := fight.SimulationOptions{}
	if flags.PlayerCritical >= 0 {
		options.Critical.PlayerChance = &flags.PlayerCritical
	}
	if flags.MonsterCritical >= 0 {
		options.Critical.MonsterChance = &flags.MonsterCritical
	}
	return options
}

func registerSimulationFlags(cmd *cobra.Command) error {
	err := utils.RegisterFlags[simulationFlags](cmd)
	if err != nil {
		return err
	}
	for _, slot := range simulationSlotNames {
		flagName := slot + "_slot"
		err = cmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return database.Items.Keys(), cobra.ShellCompDirectiveNoFileComp
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(simulationCmd)
}
