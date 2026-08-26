package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/spf13/cobra"
)

var simulationAPICmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "api <monster>",
	Short: "Simulate a fight using the official API",
	Long: `Simulate a fight using the official API

Arguments:
  monster   The code of the monster.`,
	ValidArgsFunction: completion.Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		monster := args[0]
		input, err := readSimulationInput(cmd, monster, false)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return simulationAPIRun(input)
	},
}

func simulationAPIRun(input simulationInput) error {
	request, err := simulationRequest(input.Monster.Code, input.Character, input.Flags)
	if err != nil {
		return err
	}
	results, err := fight.SimulateAPI(request)
	if err != nil {
		return err
	}
	fighter := fight.FromLoadout(input.Character.Level, fight.CharacterSlots(input.Character), fight.CharacterUtilities(input.Character))
	result := fight.Report(fighter, input.Character.Level, input.Monster, results)
	if !input.Flags.Logs {
		for i := range result.Results {
			result.Results[i].Logs = []string{}
		}
	}
	return console.Auto(result)
}

func init() {
	simulationCmd.AddCommand(simulationAPICmd)
	err := registerSimulationFlags(simulationAPICmd)
	if err != nil {
		panic(err)
	}
}
