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

var groupSimulationCompareCmd = &cobra.Command{
	Args:              cobra.ExactArgs(1),
	Use:               "compare <monster>",
	Short:             "Compare the group simulators",
	ValidArgsFunction: completion.Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := readGroupSimulationFlags(cmd)
		if err != nil {
			return err
		}
		monster, err := groupSimulationMonster(args)
		if err != nil {
			return err
		}
		fighters, err := group.FightersFromFile(flags.File)
		if err != nil {
			return err
		}
		summary := group.SimulateMany(fighters, group.MonsterFromSchema(*monster), group.SimulationOptions{
			Iterations: flags.Iterations,
		})
		localResults := groupSimulationResults(summary.Results, true)
		characters, err := group.CharactersFromFile(flags.File)
		if err != nil {
			return err
		}
		request := schemas.CombatSimulationRequestSchema{
			Characters: characters,
			Monster:    args[0],
			Iterations: flags.Iterations,
		}
		remote, err := api.SimulationFight(request)
		if err != nil {
			return err
		}
		diagnostics := groupSimulationDiagnostics(remote.Results)
		if !flags.Logs {
			for i := range remote.Results {
				remote.Results[i].Logs = []string{}
			}
		}
		apiOutput := groupSimulationOutput{
			Losses:      remote.Losses,
			Results:     remote.Results,
			Winrate:     remote.Winrate,
			Wins:        remote.Wins,
			Diagnostics: diagnostics,
		}
		local := groupSimulationReport(localResults, flags.Logs)
		comparison := groupSimulationComparison{API: apiOutput, Local: local}
		if apiOutput.Wins != local.Wins {
			comparison.Differences = append(comparison.Differences, "wins")
		}
		if apiOutput.Losses != local.Losses {
			comparison.Differences = append(comparison.Differences, "losses")
		}
		return console.Auto(comparison)
	},
}

func init() {
	groupSimulationCmd.AddCommand(groupSimulationCompareCmd)
	err := utils.RegisterFlags[groupSimulationFlags](groupSimulationCompareCmd)
	if err != nil {
		panic(err)
	}
}
