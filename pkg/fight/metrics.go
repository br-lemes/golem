package fight

import (
	"math"

	"github.com/br-lemes/golem/pkg/schemas"
)

type SimulationMetrics struct {
	AverageFightCooldown float32
	XP                   int
	XPPerCycle           float32
}

func Metrics(player Fighter, level int, monster schemas.MonsterSchema, results []schemas.CombatResultSchema) SimulationMetrics {
	var turns float32
	for _, result := range results {
		turns += float32(result.Turns)
	}
	if len(results) > 0 {
		turns /= float32(len(results))
	}
	cooldown := float32(FightCooldown(int(math.Round(float64(turns))), player.Stats.Haste))
	xp := CombatXPForLevel(level, player.Stats.Wisdom, monster)
	metrics := SimulationMetrics{AverageFightCooldown: cooldown, XP: xp}
	if cooldown > 0 {
		metrics.XPPerCycle = float32(xp) * winrate(results) / 100 / cooldown
	}
	return metrics
}

func winrate(results []schemas.CombatResultSchema) float32 {
	if len(results) == 0 {
		return 0
	}
	wins := 0
	for _, result := range results {
		if result.Result == "win" {
			wins++
		}
	}
	return float32(wins) * 100 / float32(len(results))
}

// FightCooldown returns the official cooldown for a fight result. Haste is a
// percentage: one point reduces one percent of the unmodified cooldown.
func FightCooldown(turns, haste int) int {
	cooldown := float64(turns*2) * (1 - float64(haste)/100)
	if cooldown < 5 {
		return 5
	}
	return int(math.Round(cooldown))
}

// CombatXP calculates solo combat XP using the formula documented in
// artifacts-docs/concepts/stats_and_fights.
func CombatXP(character schemas.CharacterSchema, monster schemas.MonsterSchema) int {
	return CombatXPForLevel(character.Level, character.Wisdom, monster)
}

func CombatXPForLevel(level, wisdom int, monster schemas.MonsterSchema) int {
	diff := level - monster.Level
	levelPenalty := 1.0
	switch {
	case diff > 10:
		levelPenalty = 0
	case diff > 5:
		levelPenalty = 0.7
	}

	monsterMultiplier := 1.0
	switch monster.Type {
	case "elite":
		monsterMultiplier = 1.4
	case "boss", "raid_boss":
		monsterMultiplier = 2.0
	}

	base := float64(monster.Level)/float64(level)*20 + float64(monster.Hp)*0.04
	return int(math.Round(base * levelPenalty * monsterMultiplier * (1 + float64(wisdom)*0.001)))
}
