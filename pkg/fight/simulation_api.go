package fight

import (
	"time"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

func SimulateAPI(request schemas.CombatSimulationRequestSchema) ([]schemas.CombatResultSchema, error) {
	results := []schemas.CombatResultSchema{}
	var lastRequest time.Time
	for _, iterations := range IterationChunks(request.Iterations) {
		request.Iterations = iterations
		if !lastRequest.IsZero() {
			wait := time.Second - time.Since(lastRequest)
			if wait > 0 {
				time.Sleep(wait)
			}
		}
		lastRequest = time.Now()
		response, err := api.SimulationFight(request)
		if err != nil {
			return nil, err
		}
		results = append(results, response.Results...)
	}
	return results, nil
}
