package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/spf13/cobra"
)

var simulationCompareCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "compare <monster>",
	Short: "Compare the official and local simulators",
	Long: `Compare the official and local simulators

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
		return simulationCompareRun(input)
	},
}

func simulationCompareRun(input simulationInput) error {
	request, err := simulationRequest(input.Monster.Code, input.Character, input.Flags)
	if err != nil {
		return err
	}
	comparison, err := fight.CompareSimulations(request, input.Monster, simulationCriticalOptions(input.Flags), input.Flags.Logs)
	if err != nil {
		return err
	}
	return console.Auto(comparison)
}

func init() {
	simulationCmd.AddCommand(simulationCompareCmd)
	err := registerSimulationFlags(simulationCompareCmd)
	if err != nil {
		panic(err)
	}
}
