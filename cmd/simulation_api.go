package cmd

import (
	"fmt"
	"time"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/simulation"
	"github.com/spf13/cobra"
)

var simulationAPICmd = &cobra.Command{
	Use:   "api <monster>",
	Short: "Simulate a fight using the official API",
	Long: `Simulate a fight using the official member-only simulator.

Use --file to submit a JSON array containing one character. The positional
monster and --iterations flag are supplied by the command. Without --file,
the flags describe one fake character and are useful for quick experiments.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := readSimulationScalarFlags(cmd)
		if err != nil {
			return err
		}
		err = validateSimulationScalarFlags(flags)
		if err != nil {
			return err
		}
		if cmd.Flags().Changed("player-critical") || cmd.Flags().Changed("monster-critical") {
			return fmt.Errorf("critical overrides are supported by simulation local and compare, but not by the API")
		}
		request, err := simulationRequest(cmd, args[0], flags)
		if err != nil {
			return err
		}
		var result simulationCommandOutput
		var lastRequest time.Time
		for _, iterations := range simulation.IterationChunks(flags.Iterations) {
			request.Iterations = iterations
			if !lastRequest.IsZero() && time.Since(lastRequest) < time.Second {
				time.Sleep(time.Second - time.Since(lastRequest))
			}
			lastRequest = time.Now()
			batchResult, requestErr := api.SimulationFight(request)
			if requestErr != nil {
				return requestErr
			}
			if result.Results == nil {
				result.Results = batchResult.Results
			} else {
				result.Results = append(result.Results, batchResult.Results...)
			}
		}
		result.Wins = countSimulationWins(result.Results)
		result.Losses = len(result.Results) - result.Wins
		if len(result.Results) > 0 {
			result.Winrate = float32(result.Wins) * 100 / float32(len(result.Results))
		}
		if !flags.Logs {
			for i := range result.Results {
				result.Results[i].Logs = nil
			}
		}
		monsterData, ok := database.Monsters.Get(args[0])
		if !ok {
			return fmt.Errorf("invalid monster: %s", args[0])
		}
		c := request.Characters[0]
		fighter := fight.FromLoadout(c.Level, simulation.CharacterSlots(c), simulation.CharacterUtilities(c))
		metrics := fight.Metrics(fighter, c.Level, *monsterData, result.Results)
		averageTurns, averageFinalHP := simulationAverages(result.Results)
		result.Diagnostics = simulationDiagnosticsFor(result.Results)
		result.Diagnostics.AverageTurns = averageTurns
		result.Diagnostics.AverageFinalHP = averageFinalHP
		result.Diagnostics.AverageFightCooldown = metrics.AverageFightCooldown
		result.Diagnostics.XP = metrics.XP
		result.Diagnostics.XPPerCycle = metrics.XPPerCycle
		return console.Auto(result)
	},
}

func simulationRequest(cmd *cobra.Command, monster string, flags simulationScalarFlags) (schemas.CombatSimulationRequestSchema, error) {
	character, err := resolveSimulationCharacter(cmd, flags)
	if err != nil {
		return schemas.CombatSimulationRequestSchema{}, err
	}
	request := schemas.CombatSimulationRequestSchema{
		Characters: []schemas.FakeCharacterSchema{character},
		Monster:    monster,
		Iterations: flags.Iterations,
	}
	return request, simulation.ValidateRequest(request)
}

func init() {
	simulationCmd.AddCommand(simulationAPICmd)
}
