package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/simulation"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type simulationScalarFlags struct {
	File                 string `flag:"file" desc:"JSON array file containing one simulation character"`
	Name                 string `flag:"name" desc:"Use a real character as the base"`
	Iterations           int    `flag:"iterations" default:"10" desc:"Number of simulations to run"`
	Level                int    `flag:"level" default:"1" desc:"Fake character level (1-50)"`
	Logs                 bool   `flag:"logs" desc:"Include combat logs in the output"`
	PlayerCritical       int    `flag:"player-critical" default:"-1" desc:"Override player critical chance (0-100)"`
	MonsterCritical      int    `flag:"monster-critical" default:"-1" desc:"Override monster critical chance (0-100)"`
	Utility1SlotQuantity int    `flag:"utility1_slot_quantity" default:"1" desc:"Utility 1 quantity"`
	Utility2SlotQuantity int    `flag:"utility2_slot_quantity" default:"1" desc:"Utility 2 quantity"`
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

var simulationCmd = &cobra.Command{
	Use:   "simulation",
	Short: "Run combat simulations",
}

type simulationCommandOutput struct {
	Results     []schemas.CombatResultSchema `json:"results"`
	Wins        int                          `json:"wins"`
	Losses      int                          `json:"losses"`
	Winrate     float32                      `json:"winrate"`
	Diagnostics simulationDiagnostics        `json:"diagnostics"`
}

func init() {
	rootCmd.AddCommand(simulationCmd)
	err := registerSimulationScalarFlags(simulationCmd)
	if err != nil {
		panic(err)
	}
	registerSimulationSlotFlags(simulationCmd, make(map[string]*string), true)
}

func resolveSimulationCharacter(cmd *cobra.Command, flags simulationScalarFlags) (schemas.FakeCharacterSchema, error) {
	explicitSlots := map[string]string{}
	for slot, value := range readSimulationSlotFlags(cmd) {
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
	return simulation.ResolveCharacter(simulation.CharacterOptions{
		File:              flags.File,
		Name:              flags.Name,
		Level:             flags.Level,
		ExplicitSlots:     explicitSlots,
		UtilityQuantities: explicitQuantities,
	})
}

func registerSimulationScalarFlags(cmd *cobra.Command) error {
	flags := cmd.PersistentFlags()
	flags.String("file", "", "JSON array file containing one simulation character")
	flags.String("name", "", "Use a real character as the base")
	flags.Int("iterations", 10, "Number of simulations to run")
	flags.Int("level", 1, "Fake character level (1-50)")
	flags.Bool("logs", false, "Include combat logs in the output")
	flags.Int("player-critical", -1, "Override player critical chance (0-100)")
	flags.Int("monster-critical", -1, "Override monster critical chance (0-100)")
	flags.Int("utility1_slot_quantity", 1, "Utility 1 quantity")
	flags.Int("utility2_slot_quantity", 1, "Utility 2 quantity")
	return nil
}

func readSimulationScalarFlags(cmd *cobra.Command) (simulationScalarFlags, error) {
	return utils.ReadFlags[simulationScalarFlags](cmd)
}

func validateSimulationScalarFlags(flags simulationScalarFlags) error {
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

func simulationCriticalOptions(flags simulationScalarFlags) fight.SimulationOptions {
	options := fight.SimulationOptions{}
	if flags.PlayerCritical >= 0 {
		options.Critical.PlayerChance = &flags.PlayerCritical
	}
	if flags.MonsterCritical >= 0 {
		options.Critical.MonsterChance = &flags.MonsterCritical
	}
	return options
}

func registerSimulationSlotFlags(cmd *cobra.Command, values map[string]*string, defaults bool) {
	for _, slot := range simulationSlotNames {
		value := ""
		if defaults && slot == "weapon" {
			value = "wooden_stick"
		}
		values[slot] = &value
		flagName := slot + "_slot"
		cmd.PersistentFlags().StringVar(&value, flagName, value, "Item code for the "+slot+" slot")
		_ = cmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return database.Items.Keys(), cobra.ShellCompDirectiveNoFileComp
		})
	}
}

func readSimulationSlotFlags(cmd *cobra.Command) map[string]string {
	result := map[string]string{}
	for _, slot := range simulationSlotNames {
		value, err := cmd.Flags().GetString(slot + "_slot")
		if err == nil && value != "" {
			result[slot] = value
		}
	}
	return result
}
