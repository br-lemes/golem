package best

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

func GatheringPriorities(c schemas.CharacterSchema, resource *schemas.ResourceSchema) []string {
	skill := string(resource.Skill)
	priorities := []string{skill, "prospecting"}
	level, _ := utils.GetCharacterGatheringSkillLevel(c, skill)
	if level == 50 {
		return priorities
	}
	if level-resource.Level > 10 {
		if len(resource.Drops) == 1 {
			return nil
		}
		return priorities
	}
	if len(resource.Drops) == 1 {
		return []string{"wisdom"}
	}
	if isEventResource(resource.Code) {
		return []string{skill, "prospecting", "wisdom"}
	}
	return []string{skill, "wisdom", "prospecting"}
}

func isEventResource(code string) bool {
	for _, event := range catalog.Events.All() {
		if event.Content != nil && event.Content.Code == code && event.Content.Type == "resource" {
			return true
		}
	}
	return false
}

func CraftingPriorities(c schemas.CharacterSchema, item *schemas.ItemSchema) []string {
	skillLevel, _ := utils.GetCharacterCraftingSkillLevel(c, string(*item.Craft.Skill))
	if skillLevel != 50 && skillLevel-item.Level <= 10 && skillLevel > 0 {
		return []string{"wisdom"}
	}
	return nil
}

func NormalizePriorities(priorities []string) ([]string, error) {
	validEffects := catalog.Effects().Equipments().Keys()
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
		if slices.Contains(catalog.Enums()["GatheringSkill"], effect) {
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
