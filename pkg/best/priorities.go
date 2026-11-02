package best

import (
	"fmt"
	"slices"
	"strings"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

func GatheringPriorities(c schemas.CharacterSchema, resource *schemas.ResourceSchema) []string {
	skill := resource.Skill
	priorities := []string{skill + "_gathering", "prospecting"}
	level, _ := utils.GetCharacterGatheringSkillLevel(c, skill)
	if level == 50 {
		return priorities
	}
	if level-resource.Level > 10 {
		if len(resource.Drops) == 1 {
			return []string{skill + "_gathering"}
		}
		return priorities
	}
	if len(resource.Drops) == 1 {
		return []string{skill + "_gathering", "wisdom"}
	}
	if isEventResource(resource.Code) {
		return []string{skill + "_gathering", "prospecting", "wisdom"}
	}
	return []string{skill + "_gathering", "wisdom", "prospecting"}
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
	skill := *item.Craft.Skill
	skillLevel, _ := utils.GetCharacterCraftingSkillLevel(c, skill)
	if skillLevel != 50 && skillLevel-item.Level <= 10 && skillLevel > 0 {
		return []string{skill + "_crafting", "wisdom"}
	}
	return []string{skill + "_crafting"}
}

func NormalizePriorities(priorities []string) ([]string, error) {
	validEffects := catalog.Effects().Equipments().Keys()
	seen := make(map[string]bool, len(priorities))
	result := make([]string, 0, len(priorities)+1)
	gatheringEffect := ""
	craftingEffect := ""

	for _, effect := range priorities {
		if !slices.Contains(validEffects, effect) {
			return nil, fmt.Errorf("invalid effect %q: allowed values are %v", effect, validEffects)
		}
		if seen[effect] {
			return nil, fmt.Errorf("effect specified more than once: %s", effect)
		}
		seen[effect] = true
		if isGatheringEffect(effect) {
			if craftingEffect != "" {
				return nil, fmt.Errorf("gathering and crafting effects cannot be combined: %s and %s", effect, craftingEffect)
			}
			if gatheringEffect != "" {
				return nil, fmt.Errorf("multiple gathering skills specified: %s and %s", gatheringEffect, effect)
			}
			gatheringEffect = effect
			continue
		}
		if isCraftingEffect(effect) {
			if gatheringEffect != "" {
				return nil, fmt.Errorf("gathering and crafting effects cannot be combined: %s and %s", gatheringEffect, effect)
			}
			if craftingEffect != "" {
				return nil, fmt.Errorf("multiple crafting skills specified: %s and %s", craftingEffect, effect)
			}
			craftingEffect = effect
			continue
		}
		result = append(result, effect)
	}

	skillEffect := gatheringEffect
	if skillEffect == "" {
		skillEffect = craftingEffect
	}
	if skillEffect != "" {
		result = append([]string{skillEffect}, result...)
	}
	if !seen["inventory_space"] {
		result = append(result, "inventory_space")
	}
	return result, nil
}

func isGatheringEffect(effect string) bool {
	return strings.HasSuffix(effect, "_gathering")
}

func isCraftingEffect(effect string) bool {
	return strings.HasSuffix(effect, "_crafting")
}

func isSkillEffect(effect string) bool {
	return isGatheringEffect(effect) || isCraftingEffect(effect)
}
