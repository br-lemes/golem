package utils

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

type bestResult struct {
	Code  string
	Value string
}

type bestCtx struct {
	Character       schemas.CharacterSchema
	ItemValues      map[string]int
	OwnedItems      map[string]int
	Equipped        map[string]string
	AlreadyEquipped map[string]int
	Result          map[string]bestResult
	Skill           string
	ValidItems      []schemas.ItemSchema
	Weights         map[string]int
	UniqueAdeptRing bool
}

func BestFinder(character schemas.CharacterSchema, uniqueAdeptRing bool, priorities ...string) (map[string]bestResult, error) {
	priorities, err := NormalizeBestPriorities(priorities)
	if err != nil {
		return nil, err
	}
	skill := ""
	if len(priorities) > 0 && slices.Contains(database.Enum("GatheringSkill"), priorities[0]) {
		skill = priorities[0]
	}
	weights := make(map[string]int, len(priorities))
	// Effects are compared lexicographically: a later effect can only decide
	// when all preceding effects are tied.
	weight := 1
	for i := len(priorities) - 1; i >= 0; i-- {
		weights[priorities[i]] = weight
		weight *= 10000
	}
	ctx := &bestCtx{
		Character:       character,
		Skill:           skill,
		Weights:         weights,
		UniqueAdeptRing: uniqueAdeptRing,
	}
	err = ctx.fetchItems()
	if err != nil {
		return nil, err
	}
	ctx.filterAndSort()
	ctx.matchEquipment()
	return ctx.Result, nil
}

func (c *bestCtx) fetchItems() error {
	bankItems, err := api.MyBankItems()
	if err != nil {
		return err
	}
	c.OwnedItems = make(map[string]int)
	for _, item := range bankItems {
		c.OwnedItems[item.Code] += item.Quantity
	}
	if c.Character.Inventory != nil {
		for _, item := range *c.Character.Inventory {
			c.OwnedItems[item.Code] += item.Quantity
		}
	}
	c.Equipped = map[string]string{
		"amulet":     c.Character.AmuletSlot,
		"artifact1":  c.Character.Artifact1Slot,
		"artifact2":  c.Character.Artifact2Slot,
		"artifact3":  c.Character.Artifact3Slot,
		"bag":        c.Character.BagSlot,
		"body_armor": c.Character.BodyArmorSlot,
		"boots":      c.Character.BootsSlot,
		"helmet":     c.Character.HelmetSlot,
		"leg_armor":  c.Character.LegArmorSlot,
		"ring1":      c.Character.Ring1Slot,
		"ring2":      c.Character.Ring2Slot,
		"rune":       c.Character.RuneSlot,
		"shield":     c.Character.ShieldSlot,
		"utility1":   c.Character.Utility1Slot,
		"utility2":   c.Character.Utility2Slot,
		"weapon":     c.Character.WeaponSlot,
	}
	c.AlreadyEquipped = make(map[string]int)
	for slot, code := range c.Equipped {
		if code != "" {
			qty := 1
			switch slot {
			case "utility1":
				qty = c.Character.Utility1SlotQuantity
			case "utility2":
				qty = c.Character.Utility2SlotQuantity
			}
			c.OwnedItems[code] += qty
			c.AlreadyEquipped[code] += qty
		}
	}
	return nil
}

func (c *bestCtx) filterAndSort() {
	allItems := database.Items.All()
	c.ItemValues = make(map[string]int)
	for _, item := range allItems {
		if hasNegativeInventorySpace(*item) {
			continue
		}
		if c.Skill != "" && item.Type == "weapon" && item.Subtype == "tool" {
			skillLevel := 0
			switch c.Skill {
			case "mining":
				skillLevel = c.Character.MiningLevel
			case "woodcutting":
				skillLevel = c.Character.WoodcuttingLevel
			case "fishing":
				skillLevel = c.Character.FishingLevel
			case "alchemy":
				skillLevel = c.Character.AlchemyLevel
			}
			if item.Level > skillLevel {
				continue
			}
		} else {
			if item.Level > c.Character.Level {
				continue
			}
		}
		val := c.evaluateItem(*item)
		if val == 0 {
			continue
		}
		c.ValidItems = append(c.ValidItems, *item)
		c.ItemValues[item.Code] = val
	}
	slices.SortFunc(c.ValidItems, func(i, j schemas.ItemSchema) int {
		return cmp.Compare(c.ItemValues[j.Code], c.ItemValues[i.Code])
	})
}

func (c *bestCtx) matchEquipment() {
	slots := slices.Collect(maps.Keys(database.EquipmentSlotToTypes))
	slices.Sort(slots)
	c.Result = make(map[string]bestResult)
	for i := 0; i < len(slots); {
		itemType := database.EquipmentSlotToTypes[slots[i]]
		j := i
		for j < len(slots) && database.EquipmentSlotToTypes[slots[j]] == itemType {
			j++
		}
		chosen := c.bestGroup(slots[i:j], itemType)
		for slot, code := range chosen {
			if code == "" || code == c.Equipped[slot] {
				continue
			}
			item, _ := database.Items.Get(code)
			c.Result[slot] = bestResult{Code: code, Value: c.formatItem(*item)}
		}
		i = j
	}
}

// bestGroup solves the assignment jointly because a slot-local greedy choice
// can produce unstable recommendations after equipping the suggested items.
func (c *bestCtx) bestGroup(slots []string, itemType string) map[string]string {
	available := make(map[string]int)
	for _, item := range c.ValidItems {
		if item.Type == itemType && c.OwnedItems[item.Code] > 0 {
			available[item.Code] = c.OwnedItems[item.Code]
		}
	}
	bestScore := -1
	bestMatches := -1
	best := map[string]string{}
	var visit func(int, int, map[string]string)
	visit = func(pos, score int, current map[string]string) {
		if pos == len(slots) {
			matches := 0
			for slot, code := range current {
				if code != "" && code == c.Equipped[slot] {
					matches++
				}
			}
			if score > bestScore || (score == bestScore && matches > bestMatches) {
				bestScore = score
				bestMatches = matches
				best = maps.Clone(current)
			}
			return
		}
		slot := slots[pos]
		visit(pos+1, score, current) // leave the slot unchanged/empty
		for _, item := range c.ValidItems {
			if item.Type != itemType || available[item.Code] == 0 {
				continue
			}
			if c.UniqueAdeptRing && item.Code == "ring_of_the_adept" {
				if c.AlreadyEquipped[item.Code] > 0 && c.Equipped[slot] != item.Code {
					continue
				}
				if slices.ContainsFunc(slots[:pos], func(s string) bool { return current[s] == item.Code }) {
					continue
				}
			}
			if itemType != "ring" && slices.ContainsFunc(slots[:pos], func(s string) bool { return current[s] == item.Code }) {
				continue
			}
			available[item.Code]--
			current[slot] = item.Code
			visit(pos+1, score+c.ItemValues[item.Code], current)
			delete(current, slot)
			available[item.Code]++
		}
	}
	visit(0, 0, map[string]string{})
	return best
}

func (c *bestCtx) evaluateItem(item schemas.ItemSchema) int {
	if item.Effects == nil {
		return 0
	}
	effects := make(map[string]int, len(*item.Effects))
	for _, effect := range *item.Effects {
		val := effect.Value
		if effect.Code == c.Skill {
			val = -val
		}
		effects[effect.Code] += val
	}
	score := 0
	for code, value := range effects {
		weight, exists := c.Weights[code]
		if !exists {
			continue
		}
		score += value * weight
	}
	return score
}

func hasNegativeInventorySpace(item schemas.ItemSchema) bool {
	if item.Effects == nil {
		return false
	}
	for _, effect := range *item.Effects {
		if effect.Code == "inventory_space" && effect.Value < 0 {
			return true
		}
	}
	return false
}

func (c *bestCtx) formatItem(item schemas.ItemSchema) string {
	if item.Effects == nil {
		return ""
	}
	var parts []string
	for _, effect := range *item.Effects {
		_, exists := c.Weights[effect.Code]
		if !exists {
			continue
		}
		format := "%+d %s"
		if effect.Value < 0 {
			format = "%d %s"
		}
		parts = append(parts, fmt.Sprintf(format, effect.Value, effect.Code))
	}
	return strings.Join(parts, ", ")
}
