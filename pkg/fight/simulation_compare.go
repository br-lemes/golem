package fight

import (
	"math/rand"
	"time"

	"github.com/br-lemes/golem/pkg/schemas"
)

type SimulationComparison struct {
	API         SimulationStats        `json:"api"`
	Local       SimulationStats        `json:"local"`
	Differences []SimulationDifference `json:"relevant_differences,omitempty"`
}

type SimulationStats struct {
	Wins        int                   `json:"wins"`
	Losses      int                   `json:"losses"`
	Winrate     float32               `json:"winrate"`
	Diagnostics SimulationDiagnostics `json:"diagnostics"`
	Logs        []string              `json:"logs,omitempty"`
}

type SimulationDifference struct {
	Metric string      `json:"metric"`
	API    interface{} `json:"api"`
	Local  interface{} `json:"local"`
}

func CompareSimulations(request schemas.CombatSimulationRequestSchema, monster schemas.MonsterSchema, options SimulationOptions, includeLogs bool) (SimulationComparison, error) {
	apiResults, err := SimulateAPI(request)
	if err != nil {
		return SimulationComparison{}, err
	}
	character := request.Characters[0]
	fighter := FromLoadout(character.Level, CharacterSlots(character), CharacterUtilities(character))
	options.Iterations = request.Iterations
	if options.RNG == nil {
		options.RNG = rand.New(rand.NewSource(time.Now().UnixNano())).Float64
	}
	local := SimulateMany(fighter, monster, options)
	comparison := SimulationComparison{
		API: SimulationStats{
			Wins:   CountSimulationWins(apiResults),
			Losses: len(apiResults) - CountSimulationWins(apiResults),
		},
		Local: SimulationStats{
			Wins:    local.Wins,
			Losses:  local.Losses,
			Winrate: local.Winrate,
		},
	}
	if len(apiResults) > 0 {
		comparison.API.Winrate = float32(comparison.API.Wins) * 100 / float32(len(apiResults))
	}
	apiTurns, apiHP := SimulationAverages(apiResults)
	localTurns, localHP := SimulationAverages(local.Results)
	apiMetrics := Metrics(fighter, character.Level, monster, apiResults)
	localMetrics := Metrics(fighter, character.Level, monster, local.Results)
	comparison.API.Diagnostics = SimulationDiagnosticsFor(apiResults)
	comparison.API.Diagnostics.AverageTurns, comparison.API.Diagnostics.AverageFinalHP = apiTurns, apiHP
	comparison.API.Diagnostics.AverageFightCooldown, comparison.API.Diagnostics.XP, comparison.API.Diagnostics.XPPerCycle = apiMetrics.AverageFightCooldown, apiMetrics.XP, apiMetrics.XPPerCycle
	comparison.Local.Diagnostics = SimulationDiagnosticsFor(local.Results)
	comparison.Local.Diagnostics.AverageTurns, comparison.Local.Diagnostics.AverageFinalHP = localTurns, localHP
	comparison.Local.Diagnostics.AverageFightCooldown, comparison.Local.Diagnostics.XP, comparison.Local.Diagnostics.XPPerCycle = localMetrics.AverageFightCooldown, localMetrics.XP, localMetrics.XPPerCycle
	if includeLogs && len(apiResults) > 0 && len(local.Results) > 0 {
		comparison.API.Logs = apiResults[0].Logs
		comparison.Local.Logs = local.Results[0].Logs
	}
	if comparison.API.Wins != comparison.Local.Wins {
		comparison.Differences = append(comparison.Differences, SimulationDifference{
			Metric: "wins",
			API:    comparison.API.Wins,
			Local:  comparison.Local.Wins,
		})
	}
	if comparison.API.Losses != comparison.Local.Losses {
		comparison.Differences = append(comparison.Differences, SimulationDifference{
			Metric: "losses",
			API:    comparison.API.Losses,
			Local:  comparison.Local.Losses,
		})
	}
	return comparison, nil
}
