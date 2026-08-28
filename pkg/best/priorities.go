package best

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

func GatheringPriorities(c schemas.CharacterSchema, resource *schemas.ResourceSchema) []string {
	skill := string(resource.Skill)
	priorities := []string{skill, "prospecting"}
	level, _ := utils.GetCharacterGatheringSkillLevel(c, skill)
	if level-resource.Level <= 10 {
		priorities = []string{skill, "wisdom", "prospecting"}
	}
	return priorities
}

func CraftingPriorities(c schemas.CharacterSchema, item *schemas.ItemSchema) []string {
	skillLevel, _ := utils.GetCharacterCraftingSkillLevel(c, string(*item.Craft.Skill))
	priorities := []string{"inventory_space"}
	if skillLevel-item.Level <= 10 && skillLevel > 0 {
		priorities = []string{"wisdom", "inventory_space"}
	}
	return priorities
}

func NormalizePriorities(priorities []string) ([]string, error) {
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
