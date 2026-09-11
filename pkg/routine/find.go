package routine

import (
	"slices"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

type Result struct {
	Target      schemas.MapSchema
	Transitions []schemas.MapSchema
	Costs       []schemas.ConditionSchema
	Distance    int
	Potion      *schemas.ItemSchema
}

type node struct {
	point       catalog.Point
	transitions []schemas.MapSchema
	costs       []schemas.ConditionSchema
	distance    int
}

type eventPointSet map[catalog.Point]bool

func Find(character schemas.CharacterSchema, code string, potions []schemas.SimpleItemSchema) []Result {
	//+gocover:ignore:block production wrapper over tested implementation
	return find(defaultDeps, character, code, potions)
}

func find(d deps, character schemas.CharacterSchema, code string, potions []schemas.SimpleItemSchema) []Result {
	results := findFrom(d, character, code, catalog.Point{
		X:     character.X,
		Y:     character.Y,
		Layer: character.Layer,
	}, nil)
	for _, simplePotion := range potions {
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
			}, item)
			results = append(results, potionResults...)
		}
	}
	return results
}

func findFrom(d deps, character schemas.CharacterSchema, code string, start catalog.Point, potion *schemas.ItemSchema) []Result {
	events := eventPoints(d, code)
	var achievementValues []schemas.AccountAchievementSchema
	loaded := false
	load := func() bool {
		if loaded {
			//+gocover:ignore:block no repeated achievements in catalog paths
			return achievementValues != nil
		}
		loaded = true
		if d.accountsAchievements == nil {
			return false
		}
		var err error
		achievementValues, err = d.accountsAchievements(character.Account)
		return err == nil
	}
	hasAchievement := func(code string) bool {
		if !load() {
			return false
		}
		for _, achievement := range achievementValues {
			if achievement.Code == code && achievement.CompletedAt != nil {
				return true
			}
		}
		return false
	}
	conditionsSatisfied := func(conditions *[]schemas.ConditionSchema) bool {
		if conditions == nil {
			//+gocover:ignore:block catalog conditions are always lists
			return true
		}
		for _, condition := range *conditions {
			if condition.Operator == "achievement_unlocked" && !hasAchievement(condition.Code) {
				return false
			}
		}
		return true
	}

	queue := []node{{point: start}}
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
				Target:      *tile,
				Transitions: current.transitions,
				Costs:       current.costs,
				Distance:    current.distance,
				Potion:      potion,
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
				costs := appendConditions(current.costs, transition.Conditions)
				queue = append(queue, node{
					point:       next,
					transitions: transitions,
					costs:       costs,
					distance:    current.distance + 1,
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
					point:       next,
					transitions: current.transitions,
					costs:       current.costs,
					distance:    current.distance + 1,
				})
			}
		}
	}
	return results
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
