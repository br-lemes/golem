package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var groupSimulationAPICmd = &cobra.Command{
	Args:              cobra.ExactArgs(1),
	Use:               "api <monster>",
	Short:             "Run the official group simulator API",
	ValidArgsFunction: completion.Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := readGroupSimulationFlags(cmd)
		if err != nil {
			return err
		}
		_, err = groupSimulationMonster(args)
		if err != nil {
			return err
		}
		characters, err := group.CharactersFromFile(flags.File)
		if err != nil {
			return err
		}
		request := schemas.CombatSimulationRequestSchema{
			Characters: characters,
			Monster:    args[0],
			Iterations: flags.Iterations,
		}
		result, err := api.SimulationFight(request)
		if err != nil {
			return err
		}
		diagnostics := groupSimulationDiagnostics(result.Results)
		if !flags.Logs {
			for i := range result.Results {
				result.Results[i].Logs = []string{}
			}
		}
		output := groupSimulationOutput{
			Losses:      result.Losses,
			Results:     result.Results,
			Winrate:     result.Winrate,
			Wins:        result.Wins,
			Diagnostics: diagnostics,
		}
		return console.Auto(output)
	},
}

func init() {
	groupSimulationCmd.AddCommand(groupSimulationAPICmd)
	err := utils.RegisterFlags[groupSimulationFlags](groupSimulationAPICmd)
	if err != nil {
		panic(err)
	}
}
