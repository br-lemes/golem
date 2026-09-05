package api

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/schemas"
)

func SimulationFight(request schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
	resp, err := post("/simulation/fight", request)
	if err != nil {
		return schemas.CombatSimulationDataSchema{}, err
	}
	var result schemas.CombatSimulationResponseSchema
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return schemas.CombatSimulationDataSchema{}, err
	}
	return result.Data, nil
}
