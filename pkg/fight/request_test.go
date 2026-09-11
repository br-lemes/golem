package fight

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestIterationChunks(t *testing.T) {
	for _, test := range []struct {
		iterations int
		want       []int
	}{
		{iterations: 1, want: []int{1}},
		{iterations: 100, want: []int{100}},
		{iterations: 201, want: []int{100, 100, 1}},
	} {
		got := IterationChunks(test.iterations)
		if len(got) != len(test.want) {
			t.Fatalf("IterationChunks(%d) = %v, want %v", test.iterations, got, test.want)
		}
		for i := range got {
			if got[i] != test.want[i] {
				t.Fatalf("IterationChunks(%d) = %v, want %v", test.iterations, got, test.want)
			}
		}
	}
}

func TestValidateRequest(t *testing.T) {
	valid := schemas.CombatSimulationRequestSchema{
		Characters: []schemas.FakeCharacterSchema{{Level: 1}},
		Monster:    "rat",
		Iterations: 1,
	}
	err := ValidateRequest(valid)
	if err != nil {
		t.Fatalf("valid request rejected: %v", err)
	}
	tests := []struct {
		name string
		edit func(*schemas.CombatSimulationRequestSchema)
	}{
		{
			name: "character count",
			edit: func(r *schemas.CombatSimulationRequestSchema) { r.Characters = nil },
		},
		{
			name: "monster",
			edit: func(r *schemas.CombatSimulationRequestSchema) { r.Monster = "" },
		},
		{
			name: "iterations",
			edit: func(r *schemas.CombatSimulationRequestSchema) { r.Iterations = 0 },
		},
		{
			name: "level low",
			edit: func(r *schemas.CombatSimulationRequestSchema) { r.Characters[0].Level = 0 },
		},
		{
			name: "level high",
			edit: func(r *schemas.CombatSimulationRequestSchema) { r.Characters[0].Level = 51 },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := valid
			test.edit(&request)
			err := ValidateRequest(request)
			if err == nil {
				t.Fatal("request unexpectedly accepted")
			}
		})
	}
}
