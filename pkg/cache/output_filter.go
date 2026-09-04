package cache

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/models"
	"gorm.io/gorm"
)

func GetOutputFilters(command string) map[string][]string {
	filters := ListOutputFilters(command)
	result := map[string][]string{}
	for _, filter := range filters {
		result[filter.Kind] = append(result[filter.Kind], filter.Pattern)
	}
	return result
}

func ListOutputFilters(command string) []models.OutputFilter {
	var filters []models.OutputFilter
	query := cache.Order("command, kind, pattern")
	if command != "" {
		query = query.Where("command = ?", command)
	}
	if query.Find(&filters).Error != nil {
		return nil
	}
	return filters
}

func AddOutputFilter(command, kind, pattern string) error {
	return cache.Create(&models.OutputFilter{
		Command: command,
		Kind:    kind,
		Pattern: pattern,
	}).Error
}

func RemoveOutputFilter(command, kind, pattern string) (bool, error) {
	result := cache.Where("command = ? AND kind = ? AND pattern = ?", command, kind, pattern).Delete(&models.OutputFilter{})
	return result.RowsAffected > 0, result.Error
}

func EditOutputFilter(command, kind, oldPattern, newPattern string) error {
	return cache.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("command = ? AND kind = ? AND pattern = ?", command, kind, oldPattern).Delete(&models.OutputFilter{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("output filter not found")
		}
		return tx.Create(&models.OutputFilter{
			Command: command,
			Kind:    kind,
			Pattern: newPattern,
		}).Error
	})
}
