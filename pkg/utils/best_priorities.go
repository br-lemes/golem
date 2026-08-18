package utils

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/database"
)

func NormalizeBestPriorities(priorities []string) ([]string, error) {
	validEffects := database.Effects().Equipments().Keys()
	gatheringSkills := database.Enum("GatheringSkill")
	seen := make(map[string]bool, len(priorities))
	result := make([]string, 0, len(priorities)+1)
	skill := ""

	for _, effect := range priorities {
		if !slices.Contains(validEffects, effect) {
			return nil, fmt.Errorf("invalid effect %q: allowed values are %v", effect, validEffects)
		}
		if seen[effect] {
			return nil, fmt.Errorf("effect specified more than once: %s", effect)
		}
		seen[effect] = true
		if slices.Contains(gatheringSkills, effect) {
			if skill != "" {
				return nil, fmt.Errorf("multiple gathering skills specified: %s and %s", skill, effect)
			}
			skill = effect
			continue
		}
		result = append(result, effect)
	}

	if skill != "" {
		result = append([]string{skill}, result...)
	}
	if !seen["inventory_space"] {
		result = append(result, "inventory_space")
	}
	return result, nil
}
