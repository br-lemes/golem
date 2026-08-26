package fight

import "github.com/br-lemes/golem/pkg/schemas"

func CountSimulationWins(results []schemas.CombatResultSchema) int {
	wins := 0
	for _, result := range results {
		if result.Result == "win" {
			wins++
		}
	}
	return wins
}

func SimulationAverages(results []schemas.CombatResultSchema) (float32, float32) {
	var turns, hp float32
	for _, result := range results {
		turns += float32(result.Turns)
		hp += finalHP(result)
	}
	if len(results) == 0 {
		return 0, 0
	}
	return turns / float32(len(results)), hp / float32(len(results))
}

func finalHP(result schemas.CombatResultSchema) float32 {
	if len(result.CharacterResults) == 0 {
		return 0
	}
	switch value := result.CharacterResults[0]["final_hp"].(type) {
	case float64:
		return float32(value)
	case int:
		return float32(value)
	default:
		return 0
	}
}
