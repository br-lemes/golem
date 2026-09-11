package fight

import (
	"errors"
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCompareCriticalBuildsComparisonAndIncludesLogs(t *testing.T) {
	var receivedIterations int
	d := deps{
		simulationFight: func(request schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			receivedIterations = request.Iterations
			logs := []string{
				"Turn 1: Character_1 used earth attack and dealt 10 damage.",
			}
			result := schemas.CombatResultSchema{
				Result: "win",
				Turns:  1,
				Logs:   logs,
			}
			return schemas.CombatSimulationDataSchema{
				Results: []schemas.CombatResultSchema{result},
			}, nil
		},
	}
	request := schemas.CombatSimulationRequestSchema{
		Characters: []schemas.FakeCharacterSchema{
			{Level: 1, WeaponSlot: stringPointer("wooden_stick")},
		},
		Iterations: 2,
	}
	monster := schemas.MonsterSchema{Name: "Dummy", Hp: 1, Level: 1}
	comparison, err := compareCritical(d, request, monster, SimulationOptions{}, true)
	if err != nil {
		t.Fatalf("compareCritical returned error: %v", err)
	}
	if receivedIterations != 1 {
		t.Fatalf("API iterations = %d, want 1", receivedIterations)
	}
	if comparison.API.Wins != 1 || comparison.Local.Wins != 2 {
		t.Fatalf("comparison stats = %+v", comparison)
	}
	if len(comparison.API.Logs) != 1 || len(comparison.Local.Logs) == 0 {
		t.Fatalf("comparison logs = API %v, local %v", comparison.API.Logs, comparison.Local.Logs)
	}
	if !strings.Contains(strings.Join(comparison.Local.Logs, "\n"), "Fight result: win") {
		t.Fatal("local logs do not contain the fight result")
	}
}

func TestCompareCriticalReturnsAPIError(t *testing.T) {
	want := errors.New("simulation unavailable")
	d := deps{
		simulationFight: func(schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			return schemas.CombatSimulationDataSchema{}, want
		},
	}
	request := schemas.CombatSimulationRequestSchema{Iterations: 1}
	_, err := compareCritical(d, request, schemas.MonsterSchema{}, SimulationOptions{}, false)
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

func TestCompareCriticalRejectsEmptyAPIResults(t *testing.T) {
	d := deps{
		simulationFight: func(schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			return schemas.CombatSimulationDataSchema{}, nil
		},
	}
	request := schemas.CombatSimulationRequestSchema{Iterations: 1}
	_, err := compareCritical(d, request, schemas.MonsterSchema{}, SimulationOptions{}, false)
	if err == nil || err.Error() != "API returned no combat result" {
		t.Fatalf("error = %v, want empty result error", err)
	}
}
