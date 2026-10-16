package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var groupSimulationLocalCmd = &cobra.Command{
	Args:              cobra.ExactArgs(1),
	Use:               "local <monster>",
	Short:             "Run the local group simulator",
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
		results := groupSimulationResults(summary.Results, true)
		return console.Auto(groupSimulationReport(results, flags.Logs))
	},
}

func init() {
	groupSimulationCmd.AddCommand(groupSimulationLocalCmd)
	err := utils.RegisterFlags[groupSimulationFlags](groupSimulationLocalCmd)
	if err != nil {
		panic(err)
	}
}
