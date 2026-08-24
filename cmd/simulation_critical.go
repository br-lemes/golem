package cmd

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/simulation"
	"github.com/spf13/cobra"
)

var attackLogRE = regexp.MustCompile(`^Turn ([0-9]+): (.+?) used .* attack`)

func criticalSequenceFromLogs(logs []string) []bool {
	sequence := make([]bool, 0)
	lastKey := ""
	for _, log := range logs {
		match := attackLogRE.FindStringSubmatch(log)
		if match == nil {
			continue
		}
		key := match[1] + "\x00" + match[2]
		if key == lastKey {
			continue
		}
		lastKey = key
		sequence = append(sequence, strings.Contains(log, "Critical strike"))
	}
	return sequence
}

var simulationCriticalCmd = &cobra.Command{
	Use:   "critical <monster>",
	Short: "Compare local fights using the API critical sequence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := readSimulationScalarFlags(cmd)
		if err != nil {
			return err
		}
		err = validateSimulationScalarFlags(flags)
		if err != nil {
			return err
		}
		request, err := simulationRequest(cmd, args[0], flags)
		if err != nil {
			return err
		}
		request.Iterations = 1
		remote, err := api.SimulationFight(request)
		if err != nil {
			return err
		}
		if len(remote.Results) == 0 {
			return fmt.Errorf("API returned no combat result")
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
		sequence := criticalSequenceFromLogs(remote.Results[0].Logs)
		iterations := flags.Iterations
		if iterations < 1 {
			iterations = 1
		}
		localOptions := simulationCriticalOptions(flags)
		localOptions.Iterations = iterations
		localOptions.CriticalSequence = sequence
		localOptions.RNG = rand.New(rand.NewSource(time.Now().UnixNano())).Float64
		local := fight.SimulateMany(fighter, *monster, localOptions)
		result := simulationComparison{
			TotalIterations: iterations,
			API: simulationStats{
				Wins:   countSimulationWins(remote.Results),
				Losses: len(remote.Results) - countSimulationWins(remote.Results),
			},
			Local: simulationStats{
				Wins:    local.Wins,
				Losses:  local.Losses,
				Winrate: local.Winrate,
			},
		}
		result.API.Winrate = float32(result.API.Wins) * 100 / float32(len(remote.Results))
		apiAverageTurns, apiAverageFinalHP := simulationAverages(remote.Results)
		localAverageTurns, localAverageFinalHP := simulationAverages(local.Results)
		apiMetrics := fight.Metrics(fighter, character.Level, *monster, remote.Results)
		localMetrics := fight.Metrics(fighter, character.Level, *monster, local.Results)
		result.API.Diagnostics = simulationDiagnosticsFor(remote.Results)
		result.API.Diagnostics.AverageTurns, result.API.Diagnostics.AverageFinalHP = apiAverageTurns, apiAverageFinalHP
		result.API.Diagnostics.AverageFightCooldown, result.API.Diagnostics.XP, result.API.Diagnostics.XPPerCycle = apiMetrics.AverageFightCooldown, apiMetrics.XP, apiMetrics.XPPerCycle
		result.Local.Diagnostics = simulationDiagnosticsFor(local.Results)
		result.Local.Diagnostics.AverageTurns, result.Local.Diagnostics.AverageFinalHP = localAverageTurns, localAverageFinalHP
		result.Local.Diagnostics.AverageFightCooldown, result.Local.Diagnostics.XP, result.Local.Diagnostics.XPPerCycle = localMetrics.AverageFightCooldown, localMetrics.XP, localMetrics.XPPerCycle
		if flags.Logs {
			result.APILogs = remote.Results[0].Logs
			if len(local.Results) > 0 {
				result.LocalLogs = local.Results[0].Logs
			}
		}
		return console.Auto(result)
	},
}

func init() { simulationCmd.AddCommand(simulationCriticalCmd) }
