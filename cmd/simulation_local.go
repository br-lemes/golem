package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/simulation"
	"github.com/spf13/cobra"
)

var simulationLocalCmd = &cobra.Command{
	Args:              cobra.ExactArgs(1),
	Use:               "local <monster>",
	Short:             "Simulate a fight using the local simulator",
	ValidArgsFunction: completion.Monster(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := readSimulationScalarFlags(cmd)
		if err != nil {
			return err
		}
		err = validateSimulationScalarFlags(flags)
		if err != nil {
			return err
		}
		monster, ok := database.Monsters.Get(args[0])
		if !ok {
			return fmt.Errorf("invalid monster: %s", args[0])
		}
		character, err := resolveSimulationCharacter(cmd, flags)
		if err != nil {
			return err
		}
		fighter := fight.FromLoadout(character.Level, simulation.CharacterSlots(character), simulation.CharacterUtilities(character))
		options := simulationCriticalOptions(flags)
		options.Iterations = flags.Iterations
		options.RNG = rand.New(rand.NewSource(time.Now().UnixNano())).Float64
		summary := fight.SimulateMany(fighter, *monster, options)
		if !flags.Logs {
			for i := range summary.Results {
				summary.Results[i].Logs = []string{}
			}
		}
		metrics := fight.Metrics(fighter, character.Level, *monster, summary.Results)
		averageTurns, averageFinalHP := simulationAverages(summary.Results)
		output := simulationCommandOutput{
			Results: summary.Results,
			Wins:    summary.Wins,
			Losses:  summary.Losses,
			Winrate: summary.Winrate,
		}
		output.Diagnostics = simulationDiagnosticsFor(summary.Results)
		output.Diagnostics.AverageTurns = averageTurns
		output.Diagnostics.AverageFinalHP = averageFinalHP
		output.Diagnostics.AverageFightCooldown = metrics.AverageFightCooldown
		output.Diagnostics.XP = metrics.XP
		output.Diagnostics.XPPerCycle = metrics.XPPerCycle
		return console.Auto(output)
	},
}

func init() {
	simulationCmd.AddCommand(simulationLocalCmd)
}
