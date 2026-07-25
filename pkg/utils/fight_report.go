package utils

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

// 2000 yields ±1pp error at ~50% win rate; sufficient precision and sub-second.
const monteCarloIterations = 2000

type FightReport struct {
	Win            bool
	Safe           bool
	SuccessRatePct float64
	Turns          int
	HpRemaining    int
	MaxHp          int
	HpLostPct      float64
	Margin         float64
	PlayerFirst    bool
	TimedOut       bool
	XP             int
	Equipment      map[string]string
}

func FightSimulate(character schemas.CharacterSchema, monster schemas.MonsterSchema, overrides map[string]string) (FightReport, error) {
	worn := equippedSlotMap(character)
	resolved := make(map[string]string, len(allSlots))
	for _, s := range allSlots {
		resolved[s] = worn[s]
	}
	for slot, code := range overrides {
		if code == "" {
			resolved[slot] = ""
			continue
		}
		if _, ok := database.GetItem(code); !ok {
			return FightReport{}, fmt.Errorf("invalid item code: %s", code)
		}
		resolved[slot] = code
	}

	base := baseStats(character)
	var codes []string
	equipment := make(map[string]string)
	for _, s := range allSlots {
		if resolved[s] != "" {
			codes = append(codes, resolved[s])
			equipment[s] = resolved[s]
		}
	}
	fighter := applyGear(base, codes)

	ev := Simulate(fighter, monster, SimOptions{})
	worst := Simulate(fighter, monster, SimOptions{Pessimistic: true})

	rng := rand.New(rand.NewSource(1))
	wins := 0
	for range monteCarloIterations {
		r := Simulate(fighter, monster, SimOptions{Rng: rng.Float64})
		if r.Win {
			wins++
		}
	}

	return FightReport{
		Win:            ev.Win,
		Safe:           worst.Win,
		SuccessRatePct: float64(wins) / float64(monteCarloIterations),
		Turns:          ev.Turns,
		HpRemaining:    ev.HpRemaining,
		MaxHp:          ev.MaxHp,
		HpLostPct:      ev.HpLostPct,
		Margin:         ev.Margin,
		PlayerFirst:    ev.PlayerFirst,
		TimedOut:       ev.TimedOut,
		XP:             combatXP(character, monster),
		Equipment:      equipment,
	}, nil
}

// Solo XP formula (docs/concepts/stats_and_fights). Strict '>' used for level
// diff because tier 0 starts at level 11+ in play. raid_boss defaults to 2.0
// (same as boss), pending confirmation.
func combatXP(character schemas.CharacterSchema, monster schemas.MonsterSchema) int {
	diff := character.Level - monster.Level
	levelPenalty := 1.0
	switch {
	case diff > 10:
		levelPenalty = 0
	case diff > 5:
		levelPenalty = 0.7
	}

	monsterMultiplier := 1.0
	switch string(monster.Type) {
	case "elite":
		monsterMultiplier = 1.4
	case "boss", "raid_boss":
		monsterMultiplier = 2.0
	}

	wisdomBonus := 1 + float64(character.Wisdom)*0.001
	base := (float64(monster.Level)/float64(character.Level))*20 +
		float64(monster.Hp)*0.04
	return int(math.Round(
		base * levelPenalty * monsterMultiplier * wisdomBonus))
}
