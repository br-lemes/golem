package routine

import (
	"slices"
	"sort"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

type Result struct {
	Target       schemas.MapSchema
	Transitions  []schemas.MapSchema
	Requirements []schemas.ConditionSchema
	Distance     int
	Potion       *schemas.ItemSchema
}

type node struct {
	point        catalog.Point
	transitions  []schemas.MapSchema
	requirements []schemas.ConditionSchema
	distance     int
}

type eventPointSet map[catalog.Point]bool

func Find(character schemas.CharacterSchema, code string, potions []schemas.SimpleItemSchema) []Result {
	//+gocover:ignore:block production wrapper over tested implementation
	return find(defaultDeps, character, code, potions, eventPoints(defaultDeps, code))
}

func find(d deps, character schemas.CharacterSchema, code string, potions []schemas.SimpleItemSchema, events eventPointSet) []Result {
	results := findFrom(d, character, code, catalog.Point{
		X:     character.X,
		Y:     character.Y,
		Layer: character.Layer,
	}, nil, d.hasAchievement, events)
	for _, simplePotion := range potions {
		if simplePotion.Quantity <= 0 {
			continue
		}
		item, exists := catalog.Items().Get(simplePotion.Code)
		if !exists || item.Effects == nil {
			continue
		}
		for _, effect := range *item.Effects {
			if effect.Code != "teleport" {
				continue
			}
			target, exists := catalog.Maps.Find(func(tile *schemas.MapSchema) bool {
				return tile.MapId == effect.Value
			})
			if !exists {
				//+gocover:ignore:block teleport destinations exist in catalog
				continue
			}
			potionResults := findFrom(d, character, code, catalog.Point{
				X:     target.X,
				Y:     target.Y,
				Layer: target.Layer,
			}, item, d.hasAchievement, events)
			results = append(results, potionResults...)
		}
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})
	return results
}

func findFrom(d deps, character schemas.CharacterSchema, code string, start catalog.Point, potion *schemas.ItemSchema, hasAchievement func(string) (bool, error), events eventPointSet) []Result {
	var bankItems []schemas.SimpleItemSchema
	bankLoaded := false
	loadBank := func() bool {
		if bankLoaded {
			return bankItems != nil
		}
		bankLoaded = true
		if d.myBankItems == nil {
			return false
		}
		var err error
		bankItems, err = d.myBankItems()
		return err == nil
	}
	conditionsSatisfied := func(conditions *[]schemas.ConditionSchema) bool {
		if conditions == nil {
			//+gocover:ignore:block catalog conditions are always lists
			return true
		}
		for _, condition := range *conditions {
			switch condition.Operator {
			case "achievement_unlocked":
				if hasAchievement == nil {
					return false
				}
				unlocked, err := hasAchievement(condition.Code)
				if err != nil || !unlocked {
					return false
				}
			case "has_item":
				if !hasItem(character, condition.Code, condition.Value, loadBank, &bankItems) {
					return false
				}
			case "cost":
				continue
			case "gt", "eq", "lt", "ne":
				level, exists := utils.GetCharacterConditionLevel(character, condition.Code)
				if !exists || !conditionSatisfied(level, condition) {
					return false
				}
			default:
				return false
			}
		}
		return true
	}
	if potion != nil && !conditionsSatisfied(potion.Conditions) {
		return nil
	}

	initialDistance := 0
	if potion != nil {
		initialDistance = 5
	}
	queue := []node{{point: start, distance: initialDistance}}
	visited := map[catalog.Point]bool{start: true}
	var results []Result
	dx := [...]int{0, 0, 1, -1}
	dy := [...]int{1, -1, 0, 0}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		tile, exists := catalog.Maps.Get(current.point)
		if !exists {
			//+gocover:ignore:block lookup only enqueues catalog tiles
			continue
		}
		target := isTarget(*tile, current.point, code, events)
		conditions := tile.Access.Conditions
		if target && conditionsSatisfied(appendItemConditions(conditions, potion)) {
			results = append(results, Result{
				Target:       *tile,
				Transitions:  current.transitions,
				Requirements: appendConditions(appendConditions(nil, &current.requirements), tile.Access.Conditions),
				Distance:     current.distance,
				Potion:       potion,
			})
		}
		transition := tile.Interactions.Transition
		if transition != nil {
			next := catalog.Point{
				X:     transition.X,
				Y:     transition.Y,
				Layer: transition.Layer,
			}
			if !visited[next] && conditionsSatisfied(transition.Conditions) {
				visited[next] = true
				transitions := appendCopy(current.transitions, *tile)
				requirements := appendConditions(current.requirements, transition.Conditions)
				queue = append(queue, node{
					point:        next,
					transitions:  transitions,
					requirements: requirements,
					distance:     current.distance + 1,
				})
			}
		}
		for i := range dx {
			next := catalog.Point{
				X:     current.point.X + dx[i],
				Y:     current.point.Y + dy[i],
				Layer: current.point.Layer,
			}
			if visited[next] {
				continue
			}
			nextTile, exists := catalog.Maps.Get(next)
			if exists && nextTile.Access.Type != "blocked" {
				visited[next] = true
				queue = append(queue, node{
					point:        next,
					transitions:  current.transitions,
					requirements: current.requirements,
					distance:     current.distance + 1,
				})
			}
		}
	}
	return results
}

func currentTarget(d deps, character schemas.CharacterSchema, code string) (eventPointSet, bool) {
	point := catalog.Point{
		X:     character.X,
		Y:     character.Y,
		Layer: character.Layer,
	}
	tile, exists := catalog.Maps.Get(point)
	if exists && tile.Interactions.Content != nil && tile.Interactions.Content.Code == code {
		return nil, true
	}
	events := eventPoints(d, code)
	return events, events[point]
}

func conditionSatisfied(level int, condition schemas.ConditionSchema) bool {
	switch condition.Operator {
	case "gt":
		return level > condition.Value
	case "eq":
		return level == condition.Value
	case "lt":
		return level < condition.Value
	case "ne":
		return level != condition.Value
	default:
		return false
	}
}

func appendItemConditions(conditions *[]schemas.ConditionSchema, item *schemas.ItemSchema) *[]schemas.ConditionSchema {
	if item == nil || item.Conditions == nil {
		return conditions
	}
	result := appendConditions(nil, conditions)
	result = append(result, *item.Conditions...)
	return &result
}

func isTarget(tile schemas.MapSchema, point catalog.Point, code string, events eventPointSet) bool {
	contentTarget := tile.Interactions.Content != nil && tile.Interactions.Content.Code == code
	return contentTarget || events[point]
}

func appendCopy(values []schemas.MapSchema, value schemas.MapSchema) []schemas.MapSchema {
	result := make([]schemas.MapSchema, len(values), len(values)+1)
	copy(result, values)
	return append(result, value)
}

func appendConditions(values []schemas.ConditionSchema, conditions *[]schemas.ConditionSchema) []schemas.ConditionSchema {
	if conditions == nil {
		//+gocover:ignore:block catalog conditions are always lists
		return values
	}
	result := make([]schemas.ConditionSchema, len(values), len(values)+len(*conditions))
	copy(result, values)
	return append(result, *conditions...)
}

func eventPoints(d deps, code string) eventPointSet {
	if !slices.Contains(catalog.EventContentCodes(), code) {
		return eventPointSet{}
	}
	if d.eventsActive == nil {
		return eventPointSet{}
	}
	events, err := d.eventsActive()
	if err != nil {
		return eventPointSet{}
	}
	result := eventPointSet{}
	for _, event := range events {
		if event.Map.Interactions.Content != nil && event.Map.Interactions.Content.Code == code {
			result[catalog.Point{
				X:     event.Map.X,
				Y:     event.Map.Y,
				Layer: event.Map.Layer,
			}] = true
		}
	}
	return result
}

func hasItem(character schemas.CharacterSchema, code string, quantity int, loadBank func() bool, bankItems *[]schemas.SimpleItemSchema) bool {
	if itemQuantity(character, code) >= quantity {
		return true
	}
	if !loadBank() {
		return false
	}
	for _, item := range *bankItems {
		if item.Code == code && item.Quantity+itemQuantity(character, code) >= quantity {
			return true
		}
	}
	return false
}

func itemQuantity(character schemas.CharacterSchema, code string) int {
	quantity := 0
	if character.Inventory != nil {
		for _, item := range *character.Inventory {
			if item.Code == code {
				quantity += item.Quantity
			}
		}
	}
	for _, equipped := range []string{
		character.WeaponSlot,
		character.ShieldSlot,
		character.HelmetSlot,
		character.BodyArmorSlot,
		character.BootsSlot,
		character.LegArmorSlot,
		character.Ring1Slot,
		character.Ring2Slot,
		character.AmuletSlot,
		character.BagSlot,
		character.RuneSlot,
		character.Artifact1Slot,
		character.Artifact2Slot,
		character.Artifact3Slot,
	} {
		if equipped == code {
			quantity++
		}
	}
	quantity += equippedQuantity(character.Utility1Slot, character.Utility1SlotQuantity, code)
	quantity += equippedQuantity(character.Utility2Slot, character.Utility2SlotQuantity, code)
	return quantity
}

func equippedQuantity(slot string, quantity int, code string) int {
	if slot != code {
		return 0
	}
	if quantity == 0 {
		return 1
	}
	return quantity
}
