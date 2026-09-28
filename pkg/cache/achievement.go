package cache

import (
	"encoding/json"
	"slices"

	"github.com/br-lemes/golem/pkg/models"
)

func HasAchievement(code string) bool {
	var achievementCache models.Cache
	if !findByName(&achievementCache, "achievements") {
		return false
	}
	var codes []string
	err := json.Unmarshal([]byte(achievementCache.Data), &codes)
	if err != nil {
		return false
	}
	return slices.Contains(codes, code)
}

func CleanAchievements() {
	cache.Where("name = ?", "achievements").Delete(&models.Cache{})
}

func SaveAchievements(codes []string) {
	data, err := json.Marshal(codes)
	if err != nil {
		//+gocover:ignore:block achievement codes are always JSON serializable
		return
	}
	cache.Save(&models.Cache{Name: "achievements", Data: string(data)})
}
