package cache

import (
	"github.com/br-lemes/golem/pkg/models"
	"gorm.io/gorm"
)

func ListReactions(names, codes []string) []models.Reaction {
	var reactions []models.Reaction
	query := cache.Order("name, priority, code")
	if len(names) > 0 {
		query = query.Where("name IN ?", names)
	}
	if len(codes) > 0 {
		query = query.Where("code IN ?", codes)
	}
	if query.Find(&reactions).Error != nil {
		return nil
	}
	return reactions
}

func GetReaction(name, code string) (models.Reaction, bool) {
	var reaction models.Reaction
	result := cache.Where("name = ? AND code = ?", name, code).Limit(1).Find(&reaction)
	return reaction, result.Error == nil && result.RowsAffected > 0
}

func SetReaction(name, code string, priority *int, enabled *bool) error {
	var reaction models.Reaction
	result := cache.Where("name = ? AND code = ?", name, code).Limit(1).Find(&reaction)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		reaction.Name = name
		reaction.Code = code
		reaction.Enabled = true
		if priority == nil {
			var nextPriority int
			result = cache.Model(&models.Reaction{}).Where("name = ?", name).Select("COALESCE(MAX(priority), -1) + 1").Scan(&nextPriority)
			if result.Error != nil {
				return result.Error
			}
			reaction.Priority = nextPriority
		} else {
			reaction.Priority = *priority
		}
	} else if priority != nil {
		reaction.Priority = *priority
	}
	if enabled != nil {
		reaction.Enabled = *enabled
	}
	return cache.Save(&reaction).Error
}

func SetReactionsEnabled(code string, names []string, enabled bool) error {
	query := cache.Model(&models.Reaction{}).Where("code = ?", code)
	if len(names) > 0 {
		query = query.Where("name IN ?", names)
	}
	return query.Update("enabled", enabled).Error
}

func RemoveReactions(code string, names []string) error {
	query := cache.Where("code = ?", code)
	if len(names) > 0 {
		query = query.Where("name IN ?", names)
	}
	return query.Delete(&models.Reaction{}).Error
}

func CleanReactions() {
	cache.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Reaction{})
}
