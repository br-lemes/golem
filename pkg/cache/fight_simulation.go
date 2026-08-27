package cache

import "github.com/br-lemes/golem/pkg/models"

func GetFightSimulation(key string, version int) (models.FightSimulation, bool) {
	var simulation models.FightSimulation
	result := cache.Where("key = ? AND version = ?", key, version).Limit(1).Find(&simulation)
	return simulation, result.Error == nil && result.RowsAffected > 0
}

func SaveFightSimulation(simulation models.FightSimulation) {
	cache.Save(&simulation)
}
