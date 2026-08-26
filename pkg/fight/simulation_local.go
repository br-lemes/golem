package fight

import (
	"math/rand"
	"time"

	"github.com/br-lemes/golem/pkg/schemas"
)

func SimulateLocal(player Fighter, level int, monster schemas.MonsterSchema, options SimulationOptions, includeLogs bool) SimulationReport {
	if options.Iterations < 1 {
		options.Iterations = 1
	}
	if options.RNG == nil {
		options.RNG = defaultRNG()
	}
	summary := SimulateMany(player, monster, options)
	report := Report(player, level, monster, summary.Results)
	if !includeLogs {
		for i := range report.Results {
			report.Results[i].Logs = []string{}
		}
	}
	return report
}

func defaultRNG() RNG {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Float64
}
