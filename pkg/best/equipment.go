package best

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

type BestResult struct {
	Code  string
	Value string
}

type bestCtx struct {
	AlreadyEquipped map[string]int
	Character       schemas.CharacterSchema
	Equipped        map[string]string
	Full            bool
	ItemValues      map[string]int
	OwnedItems      map[string]int
	Result          map[string]BestResult
	Skill           string
	UniqueAdeptRing bool
	ValidItems      []schemas.ItemSchema
	Weights         map[string]int
}

type EquipmentOptions struct {
	All             bool
	Full            bool
	Owned           map[string]int
	Priorities      []string
	UniqueAdeptRing bool
	With            []string
	Without         []string
}

func FindEquipment(character schemas.CharacterSchema, options EquipmentOptions) (map[string]BestResult, error) {
	return findEquipment(defaultDeps, character, options)
}

func findEquipment(d deps, character schemas.CharacterSchema, options EquipmentOptions) (map[string]BestResult, error) {
	err := validateEquipmentOptions(options)
	if err != nil {
		return nil, err
	}
	priorities, err := NormalizePriorities(options.Priorities)
	if err != nil {
		return nil, err
	}
	skill := ""
	if len(priorities) > 0 && isSkillEffect(priorities[0]) {
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
		UniqueAdeptRing: options.UniqueAdeptRing,
		Full:            options.Full,
	}
	err = ctx.fetchItems(d, options.Owned)
	if err != nil {
		return nil, err
	}
	ctx.applyAvailability(options)
	ctx.filterAndSort()
	ctx.matchEquipment()
	return ctx.Result, nil
}

func validateEquipmentOptions(options EquipmentOptions) error {
	forbidden := make(map[string]bool, len(options.Without))
	for _, code := range options.Without {
		_, exists := catalog.Items().Get(code)
		if !exists {
			return fmt.Errorf("invalid item: %s", code)
		}
		forbidden[code] = true
	}
	for _, code := range options.With {
		item, exists := catalog.Items().Get(code)
		if !exists {
			return fmt.Errorf("invalid item: %s", code)
		}
		_, equippable := catalog.EquipmentTypeToSlots[item.Type]
		if !equippable {
			return fmt.Errorf("item cannot be equipped: %s", code)
		}
		if forbidden[code] {
			return fmt.Errorf("item cannot be forbidden and included: %s", code)
		}
	}
	return nil
}

func (c *bestCtx) applyAvailability(options EquipmentOptions) {
	if options.All {
		for _, item := range catalog.Items().Equipments().All() {
			if c.OwnedItems[item.Code] < 2 {
				c.OwnedItems[item.Code] = 2
			}
		}
	}
	for _, code := range options.With {
		if c.OwnedItems[code] < 2 {
			c.OwnedItems[code] = 2
		}
	}
	for _, code := range options.Without {
		delete(c.OwnedItems, code)
	}
}

func FindEquipmentSchemas(character schemas.CharacterSchema, options EquipmentOptions) ([]schemas.EquipSchema, error) {
	return findEquipmentSchemas(defaultDeps, character, options)
}

func findEquipmentSchemas(d deps, character schemas.CharacterSchema, options EquipmentOptions) ([]schemas.EquipSchema, error) {
	result, err := findEquipment(d, character, options)
	if err != nil {
		return nil, err
	}
	equipments := make([]schemas.EquipSchema, 0, len(result))
	for slot, item := range result {
		equipments = append(equipments, schemas.EquipSchema{
			Code: item.Code,
			Slot: slot,
		})
	}
	slices.SortFunc(equipments, func(i, j schemas.EquipSchema) int {
		return cmp.Compare(i.Slot, j.Slot)
	})
	return equipments, nil
}

func (c *bestCtx) fetchItems(d deps, owned map[string]int) error {
	suppliedOwned := owned != nil
	if owned == nil {
		bankItems, err := d.myBankItems()
		if err != nil {
			return err
		}
		owned = make(map[string]int)
		for _, item := range bankItems {
			owned[item.Code] += item.Quantity
		}
	}
	c.OwnedItems = maps.Clone(owned)
	if !suppliedOwned && c.Character.Inventory != nil {
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
			if !suppliedOwned {
				c.OwnedItems[code] += qty
			}
			c.AlreadyEquipped[code] += qty
		}
	}
	return nil
}

func (c *bestCtx) filterAndSort() {
	allItems := catalog.Items().All()
	c.ItemValues = make(map[string]int)
	for _, item := range allItems {
		if hasNegativeInventorySpace(*item) {
			continue
		}
		if !utils.MeetsItemConditions(c.Character, *item) {
			continue
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
	slots := slices.Collect(maps.Keys(catalog.EquipmentSlotToTypes))
	slices.Sort(slots)
	c.Result = make(map[string]BestResult)
	for i := 0; i < len(slots); {
		itemType := catalog.EquipmentSlotToTypes[slots[i]]
		j := i
		for j < len(slots) && catalog.EquipmentSlotToTypes[slots[j]] == itemType {
			j++
		}
		chosen := c.bestGroup(slots[i:j], itemType)
		for _, slot := range slots[i:j] {
			code := chosen[slot]
			if code == "" || (!c.Full && code == c.Equipped[slot]) {
				continue
			}
			item, _ := catalog.Items().Get(code)
			c.Result[slot] = BestResult{Code: code, Value: c.formatItem(*item)}
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
