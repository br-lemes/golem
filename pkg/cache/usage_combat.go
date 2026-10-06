package cache

import "github.com/br-lemes/golem/pkg/models"

// UsageCombatVersion must increase when aggregated combat loadout
// selection can produce different results for the same inputs. Increment it
// too whenever the individual fight simulation cache version changes.
const UsageCombatVersion = 2

const (
	UsageCombatType     = "usage"
	PotentialCombatType = "potential"
)

func GetUsageCombat(cacheType, key string) (models.UsageCombat, bool) {
	var result models.UsageCombat
	query := simulationDB.Where("type = ? AND key = ? AND version = ?", cacheType, key, UsageCombatVersion).Limit(1).Find(&result)
	return result, query.Error == nil && query.RowsAffected > 0
}

func SaveUsageCombat(result models.UsageCombat) {
	simulationDB.Save(&result)
}
