package group

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestCombatRunPartyDefeatsMonster(t *testing.T) {
	fighters := []Fighter{
		{Name: "one", Stats: Stats{HP: 100, Initiative: 100, AttackEarth: 60}},
		{Name: "two", Stats: Stats{HP: 100, Initiative: 90, AttackEarth: 60}},
	}
	monster := Monster{
		Name:  "monster",
		Stats: Stats{HP: 100},
		HP:    100,
		MaxHP: 100,
	}
	combat := NewCombat(fighters, monster, func() float64 { return 0.5 })
	result := combat.Run(10)
	if !result.Win || result.Turns != 2 {
		t.Fatalf("result = %+v, want one turn per character", result)
	}
	if len(result.FinalHP) != 2 {
		t.Fatalf("final HP = %v, want one result per player", result.FinalHP)
	}
}

func TestCombatRunMonsterDefeatsParty(t *testing.T) {
	fighters := []Fighter{{Name: "one", Stats: Stats{HP: 10, Initiative: 10}}}
	monster := Monster{
		Name:  "monster",
		Stats: Stats{AttackFire: 20, Initiative: 100},
		HP:    100,
		MaxHP: 100,
	}
	combat := NewCombat(fighters, monster, func() float64 { return 0.5 })
	result := combat.Run(10)
	if result.Win || result.FinalHP[0] != 0 {
		t.Fatalf("result = %+v, want party loss", result)
	}
}

func TestNewCombatCopiesUtilities(t *testing.T) {
	fighters := []Fighter{{
		Name:      "one",
		Stats:     Stats{HP: 100, Initiative: 10},
		Utilities: []Utility{{Code: "health_potion", Restore: 50, Quantity: 1}},
	}}
	monster := Monster{
		Name:  "monster",
		Stats: Stats{AttackFire: 60, Initiative: 100},
		HP:    100,
		MaxHP: 100,
	}

	first := NewCombat(fighters, monster, func() float64 { return 0.5 })
	first.Run(1)
	if first.Players[0].Fighter.Utilities[0].Quantity != 0 {
		t.Fatalf("first combat utility quantity = %d, want 0", first.Players[0].Fighter.Utilities[0].Quantity)
	}
	if fighters[0].Utilities[0].Quantity != 1 {
		t.Fatalf("fighter utility quantity = %d, want 1", fighters[0].Utilities[0].Quantity)
	}

	second := NewCombat(fighters, monster, func() float64 { return 0.5 })
	if second.Players[0].Fighter.Utilities[0].Quantity != 1 {
		t.Fatalf("second combat utility quantity = %d, want 1", second.Players[0].Fighter.Utilities[0].Quantity)
	}
}

func TestSimulateManyUsesInjectedRNG(t *testing.T) {
	fighters := []Fighter{{
		Name: "one",
		Stats: Stats{
			HP:             100,
			Initiative:     100,
			AttackFire:     30,
			CriticalStrike: 50,
		},
	}}
	monster := Monster{
		Name: "monster",
		Stats: Stats{
			HP:             100,
			Initiative:     90,
			AttackFire:     10,
			CriticalStrike: 50,
		},
		HP:    100,
		MaxHP: 100,
	}
	first := SimulateMany(fighters, monster, SimulationOptions{
		Iterations: 5,
		RNG:        rand.New(rand.NewSource(1)).Float64,
	})
	second := SimulateMany(fighters, monster, SimulationOptions{
		Iterations: 5,
		RNG:        rand.New(rand.NewSource(1)).Float64,
	})
	if !reflect.DeepEqual(first, second) {
		t.Fatal("simulations with the same RNG seed differ")
	}
}

func TestMetrics(t *testing.T) {
	metrics := Metrics([]Result{
		{Turns: 3, FinalHP: []float64{20, 30}},
		{Turns: 5, FinalHP: []float64{0, 10}},
	})
	if metrics.AverageTurns != 4 {
		t.Fatalf("average turns = %v, want 4", metrics.AverageTurns)
	}
	if metrics.AverageFinalHP != 30 {
		t.Fatalf("average final HP = %v, want 30", metrics.AverageFinalHP)
	}
}
