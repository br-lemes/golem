package fight

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

type deps struct {
	characters      func(string) (schemas.CharacterSchema, error)
	simulationFight func(schemas.CombatSimulationRequestSchema) (schemas.CombatSimulationDataSchema, error)
}

var defaultDeps = deps{
	characters:      api.Characters,
	simulationFight: api.SimulationFight,
}
