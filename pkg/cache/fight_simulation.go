package cache

import "github.com/br-lemes/golem/pkg/models"

var pendingFightSimulations []models.FightSimulation
var batchFightSimulations bool

func BeginFightSimulationBatch() {
	pendingFightSimulations = nil
	batchFightSimulations = true
}

func FlushFightSimulationBatch() {
	if !batchFightSimulations {
		return
	}
	batchFightSimulations = false
	if len(pendingFightSimulations) > 0 {
		simulationDB.CreateInBatches(pendingFightSimulations, 500)
	}
	pendingFightSimulations = nil
}

func GetFightSimulation(key string, version int) (models.FightSimulation, bool) {
	var simulation models.FightSimulation
	result := simulationDB.Where("key = ? AND version = ?", key, version).Limit(1).Find(&simulation)
	return simulation, result.Error == nil && result.RowsAffected > 0
}

func SaveFightSimulation(simulation models.FightSimulation) {
	if batchFightSimulations {
		pendingFightSimulations = append(pendingFightSimulations, simulation)
		return
	}
	simulationDB.Save(&simulation)
}
