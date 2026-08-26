package fight

import "github.com/br-lemes/golem/pkg/schemas"

type SimulationReport struct {
	Results     []schemas.CombatResultSchema `json:"results"`
	Wins        int                          `json:"wins"`
	Losses      int                          `json:"losses"`
	Winrate     float32                      `json:"winrate"`
	Diagnostics SimulationDiagnostics        `json:"diagnostics"`
}

func Report(player Fighter, level int, monster schemas.MonsterSchema, results []schemas.CombatResultSchema) SimulationReport {
	metrics := Metrics(player, level, monster, results)
	averageTurns, averageFinalHP := SimulationAverages(results)
	diagnostics := SimulationDiagnosticsFor(results)
	wins := CountSimulationWins(results)
	diagnostics.AverageTurns = averageTurns
	diagnostics.AverageFinalHP = averageFinalHP
	diagnostics.AverageFightCooldown = metrics.AverageFightCooldown
	diagnostics.XP = metrics.XP
	diagnostics.XPPerCycle = metrics.XPPerCycle
	winrate := float32(0)
	if len(results) > 0 {
		winrate = float32(wins) * 100 / float32(len(results))
	}
	return SimulationReport{
		Results:     results,
		Wins:        wins,
		Losses:      len(results) - wins,
		Winrate:     winrate,
		Diagnostics: diagnostics,
	}
}
