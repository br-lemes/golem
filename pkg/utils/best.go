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
	Character  schemas.CharacterSchema
	ItemValues map[string]int
	OwnedItems map[string]int
	Result     map[string]bestResult
	Skill      string
	ValidItems []schemas.ItemSchema
	Weights    map[string]int
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
	equipped := []schemas.SimpleItemSchema{
		{Code: c.Character.AmuletSlot, Quantity: 1},
		{Code: c.Character.Artifact1Slot, Quantity: 1},
		{Code: c.Character.Artifact2Slot, Quantity: 1},
		{Code: c.Character.Artifact3Slot, Quantity: 1},
		{Code: c.Character.BagSlot, Quantity: 1},
		{Code: c.Character.BodyArmorSlot, Quantity: 1},
		{Code: c.Character.BootsSlot, Quantity: 1},
		{Code: c.Character.HelmetSlot, Quantity: 1},
		{Code: c.Character.LegArmorSlot, Quantity: 1},
		{Code: c.Character.Ring1Slot, Quantity: 1},
		{Code: c.Character.Ring2Slot, Quantity: 1},
		{Code: c.Character.RuneSlot, Quantity: 1},
		{Code: c.Character.ShieldSlot, Quantity: 1},
		{Code: c.Character.Utility1Slot,
			Quantity: c.Character.Utility1SlotQuantity},
		{Code: c.Character.Utility2Slot,
			Quantity: c.Character.Utility2SlotQuantity},
		{Code: c.Character.WeaponSlot, Quantity: 1},
	}
	for _, item := range equipped {
		if item.Code != "" {
			c.OwnedItems[item.Code] += item.Quantity
		}
	}
	return nil
}

func (c *bestCtx) filterAndSort() {
	allItems := database.GetItems()
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
		val := c.evaluateItem(item)
		if val == 0 {
			continue
		}
		c.ValidItems = append(c.ValidItems, item)
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
	for _, slot := range slots {
		itemType := database.EquipmentSlotToTypes[slot]
		for _, item := range c.ValidItems {
			if item.Type == itemType && c.OwnedItems[item.Code] > 0 {
				c.Result[slot] = bestResult{
					Code:  item.Code,
					Value: c.formatItem(item),
				}
				c.OwnedItems[item.Code]--
				break
			}
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
		score += effect.Value * weight
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
