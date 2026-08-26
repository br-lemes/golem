package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/spf13/cobra"
)

var simulationCriticalCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "critical <monster>",
	Short: "Compare local fights using the API critical sequence",
	Long: `Compare local fights using the API critical sequence

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
		return simulationCriticalRun(input)
	},
}

func simulationCriticalRun(input simulationInput) error {
	request, err := simulationRequest(input.Monster.Code, input.Character, input.Flags)
	if err != nil {
		return err
	}
	result, err := fight.CompareCritical(request, input.Monster, simulationCriticalOptions(input.Flags), input.Flags.Logs)
	if err != nil {
		return err
	}
	return console.Auto(result)
}

func init() {
	simulationCmd.AddCommand(simulationCriticalCmd)
	err := registerSimulationFlags(simulationCriticalCmd)
	if err != nil {
		panic(err)
	}
}
