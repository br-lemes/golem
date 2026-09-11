package fight

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCompareSimulationsBuildsComparisonWithLogsAndDifferences(t *testing.T) {
	d := deps{
		simulationFight: func(schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			result := schemas.CombatResultSchema{
				Result:           "win",
				Turns:            1,
				Logs:             []string{"remote log"},
				CharacterResults: []map[string]interface{}{{"final_hp": 120}},
			}
			return schemas.CombatSimulationDataSchema{
				Results: []schemas.CombatResultSchema{result},
			}, nil
		},
	}
	request := schemas.CombatSimulationRequestSchema{
		Characters: []schemas.FakeCharacterSchema{{Level: 1}},
		Iterations: 1,
	}
	monster := schemas.MonsterSchema{Name: "Dummy", Hp: 1, Level: 1}
	comparison, err := compareSimulations(d, request, monster, SimulationOptions{
		RNG: func() float64 { return 0 },
	}, true)
	if err != nil {
		t.Fatalf("compareSimulations returned error: %v", err)
	}
	if comparison.API.Wins != 1 || comparison.Local.Losses != 1 {
		t.Fatalf("comparison stats = %+v", comparison)
	}
	if len(comparison.Differences) != 2 {
		t.Fatalf("differences = %+v, want wins and losses", comparison.Differences)
	}
	if len(comparison.API.Logs) != 1 || len(comparison.Local.Logs) == 0 {
		t.Fatalf("comparison logs = API %v, local %v", comparison.API.Logs, comparison.Local.Logs)
	}
}

func TestCompareSimulationsPropagatesAPIError(t *testing.T) {
	want := errors.New("simulation unavailable")
	d := deps{
		simulationFight: func(schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			return schemas.CombatSimulationDataSchema{}, want
		},
	}
	request := schemas.CombatSimulationRequestSchema{
		Characters: []schemas.FakeCharacterSchema{{Level: 1}},
		Iterations: 1,
	}
	_, err := compareSimulations(d, request, schemas.MonsterSchema{}, SimulationOptions{}, false)
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

func TestCompareSimulationsInitializesDefaultRNG(t *testing.T) {
	d := deps{
		simulationFight: func(schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			return schemas.CombatSimulationDataSchema{
				Results: []schemas.CombatResultSchema{{Result: "win", Turns: 1}},
			}, nil
		},
	}
	request := schemas.CombatSimulationRequestSchema{
		Characters: []schemas.FakeCharacterSchema{
			{Level: 1, WeaponSlot: stringPointer("wooden_stick")},
		},
		Iterations: 1,
	}
	monster := schemas.MonsterSchema{Name: "Dummy", Hp: 1, Level: 1}
	comparison, err := compareSimulations(d, request, monster, SimulationOptions{}, false)
	if err != nil {
		t.Fatalf("compareSimulations returned error: %v", err)
	}
	if comparison.Local.Wins != 1 {
		t.Fatalf("local wins = %d, want 1", comparison.Local.Wins)
	}
}
