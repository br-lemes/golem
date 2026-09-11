package fight

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSimulateLocalNormalizesIterationsAndExcludesLogs(t *testing.T) {
	player := Fighter{Stats: Stats{HP: 100, AttackFire: 100, Initiative: 10}}
	monster := schemas.MonsterSchema{Name: "Dummy", Hp: 50, Initiative: 1}
	report := SimulateLocal(player, 1, monster, SimulationOptions{}, false)
	if len(report.Results) != 1 {
		t.Fatalf("result count = %d, want 1", len(report.Results))
	}
	if len(report.Results[0].Logs) != 0 {
		t.Fatalf("logs = %v, want none", report.Results[0].Logs)
	}
}

func TestSimulateLocalIncludesLogs(t *testing.T) {
	player := Fighter{Stats: Stats{HP: 100, AttackFire: 100, Initiative: 10}}
	monster := schemas.MonsterSchema{Name: "Dummy", Hp: 50, Initiative: 1}
	report := SimulateLocal(player, 1, monster, SimulationOptions{
		Iterations: 1,
		RNG:        func() float64 { return 0 },
	}, true)
	if len(report.Results) != 1 || len(report.Results[0].Logs) == 0 {
		t.Fatalf("report results/logs = %+v", report.Results)
	}
}
