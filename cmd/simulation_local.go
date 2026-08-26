package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/spf13/cobra"
)

var simulationLocalCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "local <monster>",
	Short: "Simulate a fight using the local simulator",
	Long: `Simulate a fight using the local simulator

Arguments:
  monster   The code of the monster.`,
	ValidArgsFunction: completion.Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		monster := args[0]
		input, err := readSimulationInput(cmd, monster, true)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return simulationLocalRun(input)
	},
}

func simulationLocalRun(input simulationInput) error {
	fighter := fight.FromLoadout(input.Character.Level, fight.CharacterSlots(input.Character), fight.CharacterUtilities(input.Character))
	options := simulationCriticalOptions(input.Flags)
	options.Iterations = input.Flags.Iterations
	output := fight.SimulateLocal(fighter, input.Character.Level, input.Monster, options, input.Flags.Logs)
	return console.Auto(output)
}

func init() {
	simulationCmd.AddCommand(simulationLocalCmd)
	err := registerSimulationFlags(simulationLocalCmd)
	if err != nil {
		panic(err)
	}
}
