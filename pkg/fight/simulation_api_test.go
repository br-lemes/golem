package fight

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSimulateAPIAggregatesChunkedResults(t *testing.T) {
	var calls []int
	d := deps{
		simulationFight: func(request schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			calls = append(calls, request.Iterations)
			return schemas.CombatSimulationDataSchema{
				Results: []schemas.CombatResultSchema{
					{Turns: request.Iterations},
				},
			}, nil
		},
	}
	results, err := simulateAPI(d, request101())
	if err != nil {
		t.Fatalf("simulateAPI returned error: %v", err)
	}
	if len(results) != 2 || len(calls) != 2 || calls[0] != 100 || calls[1] != 1 {
		t.Fatalf("calls/results = %v/%v, want [100 1] and two results", calls, results)
	}
}

func TestSimulateAPIPropagatesSimulationError(t *testing.T) {
	want := errors.New("simulation unavailable")
	d := deps{
		simulationFight: func(schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
			return schemas.CombatSimulationDataSchema{}, want
		},
	}
	request := schemas.CombatSimulationRequestSchema{Iterations: 1}
	results, err := simulateAPI(d, request)
	if !errors.Is(err, want) || results != nil {
		t.Fatalf("error/results = %v/%v, want propagated error and nil results", err, results)
	}
}

func request101() schemas.CombatSimulationRequestSchema {
	return schemas.CombatSimulationRequestSchema{Iterations: 101}
}
