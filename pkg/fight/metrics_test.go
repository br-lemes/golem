package fight

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFightCooldownHasMinimumAndAppliesHaste(t *testing.T) {
	got := FightCooldown(10, 0)
	if got != 20 {
		t.Fatalf("FightCooldown(10, 0) = %d, want 20", got)
	}
	got = FightCooldown(10, 50)
	if got != 10 {
		t.Fatalf("FightCooldown(10, 50) = %d, want 10", got)
	}
	got = FightCooldown(1, 100)
	if got != 5 {
		t.Fatalf("minimum cooldown = %d, want 5", got)
	}
}

func TestCombatXPUsesMonsterTypeWisdomAndLevelPenalty(t *testing.T) {
	monster := schemas.MonsterSchema{Level: 10, Hp: 100, Type: "boss"}
	got := CombatXP(10, 0, monster)
	if got != 48 {
		t.Fatalf("boss XP = %d, want 48", got)
	}
	got = CombatXP(21, 0, monster)
	if got != 0 {
		t.Fatalf("overleveled XP = %d, want 0", got)
	}
	got = CombatXP(10, 100, monster)
	if got != 53 {
		t.Fatalf("wisdom XP = %d, want 53", got)
	}
}

func TestSimulationAveragesAndWinCount(t *testing.T) {
	results := []schemas.CombatResultSchema{
		{
			Result:           "win",
			Turns:            2,
			CharacterResults: []map[string]interface{}{{"final_hp": 80}},
		},
		{
			Result:           "loss",
			Turns:            4,
			CharacterResults: []map[string]interface{}{{"final_hp": float64(0)}},
		},
	}
	turns, hp := SimulationAverages(results)
	if turns != 3 || hp != 40 || CountSimulationWins(results) != 1 {
		t.Fatalf("averages/count = %v, %v, %d", turns, hp, CountSimulationWins(results))
	}
}

func TestSimulationAveragesReturnsZerosForEmptyResults(t *testing.T) {
	turns, hp := SimulationAverages(nil)
	if turns != 0 || hp != 0 {
		t.Fatalf("empty averages = %v, %v, want 0, 0", turns, hp)
	}
}

func TestFinalHPIgnoresUnsupportedTypes(t *testing.T) {
	result := schemas.CombatResultSchema{
		CharacterResults: []map[string]interface{}{{"final_hp": "80"}},
	}
	got := finalHP(result)
	if got != 0 {
		t.Fatalf("unsupported final HP type = %v, want 0", got)
	}
}

func TestWinrateReturnsZeroForEmptyResults(t *testing.T) {
	got := winrate(nil)
	if got != 0 {
		t.Fatalf("empty winrate = %v, want 0", got)
	}
}

func TestCombatXPAppliesLevelPenaltyAndEliteMultiplier(t *testing.T) {
	monster := schemas.MonsterSchema{Level: 10, Hp: 100, Type: "elite"}
	got := CombatXP(16, 0, monster)
	if got != 16 {
		t.Fatalf("elite XP with level penalty = %d, want 16", got)
	}
}
