package api

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/schemas"
)

func SimulationFight(request schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.SimulationFight(request)
}

func (c *Client) SimulationFight(request schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error) {
	resp, err := c.Post("/simulation/fight", request)
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
