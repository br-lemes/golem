package cache

import "github.com/br-lemes/golem/pkg/models"

func GetUsageCombat(key string, version int) (models.UsageCombat, bool) {
	var result models.UsageCombat
	query := simulationDB.Where("key = ? AND version = ?", key, version).Limit(1).Find(&result)
	return result, query.Error == nil && query.RowsAffected > 0
}

func SaveUsageCombat(result models.UsageCombat) {
	simulationDB.Save(&result)
}
