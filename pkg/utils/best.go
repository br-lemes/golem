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
}

func BestFinder(character schemas.CharacterSchema, skill string, weights map[string]int) (map[string]bestResult, error) {
	ctx := &bestCtx{Character: character, Skill: skill, Weights: weights}
	err := ctx.fetchItems()
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
	usedUniqueItems := make(map[string]bool)
	for _, slot := range slots {
		itemType := database.EquipmentSlotToTypes[slot]
		isRingSlot := strings.HasPrefix(slot, "ring")
		equippedCode := c.Equipped[slot]
		for _, item := range c.ValidItems {
			if item.Type != itemType {
				continue
			}
			if c.OwnedItems[item.Code] <= 0 {
				continue
			}
			if !isRingSlot && usedUniqueItems[item.Code] {
				continue
			}
			isCurrentlyEquippedElsewhere := c.AlreadyEquipped[item.Code] > 0
			if isCurrentlyEquippedElsewhere {
				if item.Code != equippedCode {
					c.AlreadyEquipped[item.Code]--
					c.OwnedItems[item.Code]--
					usedUniqueItems[item.Code] = true
					continue
				}
			}
			if item.Code != equippedCode {
				c.Result[slot] = bestResult{
					Code:  item.Code,
					Value: c.formatItem(item),
				}
			}
			c.AlreadyEquipped[item.Code]--
			c.OwnedItems[item.Code]--
			usedUniqueItems[item.Code] = true
			break
		}
	}
}

func (c *bestCtx) evaluateItem(item schemas.ItemSchema) int {
	if item.Effects == nil {
		return 0
	}
	score := 0
	for _, effect := range *item.Effects {
		weight, exists := c.Weights[effect.Code]
		if !exists {
			continue
		}
		val := effect.Value
		if effect.Code == c.Skill {
			val = -val
		}
		score += val * weight
	}
	return score
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
