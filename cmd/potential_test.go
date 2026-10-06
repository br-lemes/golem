package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/potential"
)

func TestPotentialSimple(t *testing.T) {
	result := potentialSimple(map[string]potential.Result{
		"wratharmor": {
			Used: false,
			SelectedEquipment: []potential.SelectedEquipment{
				{Code: "darkforged_plate"},
				{Code: "diamond_armor"},
			},
		},
	})
	evaluation := result["wratharmor"]
	if evaluation.Used || len(evaluation.SelectedEquipment) != 2 || evaluation.SelectedEquipment[0] != "darkforged_plate" {
		t.Fatalf("simple result = %+v", result)
	}
}
