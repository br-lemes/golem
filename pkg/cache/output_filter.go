package cache

import "github.com/br-lemes/golem/pkg/models"

func GetOutputFilters(command string) map[string][]string {
	var filters []models.OutputFilter
	result := map[string][]string{}
	query := cache.Where("command = ?", command).Find(&filters)
	if query.Error != nil {
		return result
	}
	for _, filter := range filters {
		result[filter.Kind] = append(result[filter.Kind], filter.Pattern)
	}
	return result
}
