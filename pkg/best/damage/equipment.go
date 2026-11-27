package damage

import (
	"cmp"
	"slices"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

const (
	bonusDmg = iota
	bonusFire
	bonusWater
	bonusEarth
	bonusAir
	bonusCritical
	bonusStats
)

// bonus holds the damage modifiers contributed by equipment. Weapons only
// use the elemental entries; gear uses all of them.
type bonus [bonusStats]int

type slotRule struct {
	kind     string
	count    int
	distinct bool
}

type bonusItem struct {
	code  string
	bonus bonus
}

type bonusPool struct {
	rule  slotRule
	items []bonusItem
}

var combatSlots = []slotRule{
	{kind: "amulet", count: 1},
	{kind: "artifact", count: 3, distinct: true},
	{kind: "body_armor", count: 1},
	{kind: "boots", count: 1},
	{kind: "helmet", count: 1},
	{kind: "leg_armor", count: 1},
	{kind: "ring", count: 2},
}

var damageElements = [4]int{bonusFire, bonusWater, bonusEarth, bonusAir}

func equipmentPools(character schemas.CharacterSchema, options Options) []bonusPool {
	items := catalog.Items().Filter(func(item *schemas.ItemSchema) bool {
		if options.Without[item.Code] {
			return false
		}
		if options.Available != nil && options.Available[item.Code] < 1 {
			return false
		}
		return utils.MeetsItemConditions(character, *item)
	})
	pools := make([]bonusPool, len(combatSlots))
	for index, rule := range combatSlots {
		pools[index].rule = rule
		for _, item := range items {
			if item.Type != rule.kind {
				continue
			}
			entry := bonusItem{code: item.Code, bonus: itemBonus(item)}
			for range poolCopies(rule, options.Available, item.Code) {
				pools[index].items = append(pools[index].items, entry)
			}
		}
	}
	return pools
}

// poolCopies is how many copies of an item the pool may hold. Distinct slots
// hold one copy per code, while repeatable slots hold at most their slot count
// and never more than the account owns.
func poolCopies(rule slotRule, available map[string]int, code string) int {
	if rule.distinct {
		return 1
	}
	if available == nil {
		return rule.count
	}
	return min(rule.count, available[code])
}

func itemBonus(item *schemas.ItemSchema) bonus {
	result := bonus{}
	if item.Effects == nil {
		return result
	}
	for _, effect := range *item.Effects {
		index, ok := bonusIndex(effect.Code)
		if !ok {
			continue
		}
		result[index] += effect.Value
	}
	return result
}

func bonusIndex(code string) (int, bool) {
	switch code {
	case "dmg":
		return bonusDmg, true
	case "dmg_fire":
		return bonusFire, true
	case "dmg_water":
		return bonusWater, true
	case "dmg_earth":
		return bonusEarth, true
	case "dmg_air":
		return bonusAir, true
	case "critical_strike":
		return bonusCritical, true
	}
	return 0, false
}

// equipmentBonuses returns the independent maximum of each modifier. It is a
// superset of every loadout the character can build, used only for pruning.
func equipmentBonuses(pools []bonusPool) bonus {
	result := bonus{}
	for _, pool := range pools {
		for stat := 0; stat < bonusStats; stat++ {
			result[stat] += pool.topSum(stat, pool.rule.count)
		}
	}
	return result
}

func (pool bonusPool) topSum(stat, count int) int {
	values := make([]int, 0, len(pool.items))
	for _, item := range pool.items {
		values = append(values, item.bonus[stat])
	}
	slices.SortFunc(values, func(left, right int) int {
		return cmp.Compare(right, left)
	})
	total := 0
	for index := 0; index < count && index < len(values); index++ {
		total += values[index]
	}
	return total
}

// maxDamage finds the highest damage the weapon can deal with the best
// level-legal gear for it. The linear weights come from the damage formula
// without rounding; the refinement pass fixes the rounding leftovers.
func maxDamage(attack bonus, pools []bonusPool, monster schemas.MonsterSchema) int {
	return selectionDamage(attack, pools, monster, bestSelection(attack, pools, monster))
}

func bestSelection(attack bonus, pools []bonusPool, monster schemas.MonsterSchema) [][]int {
	selection := initialSelection(pools, damageWeights(attack, monster))
	refineSelection(attack, pools, monster, selection)
	return selection
}

func damageWeights(attack bonus, monster schemas.MonsterSchema) [bonusStats]float64 {
	resistances := [bonusStats]int{}
	resistances[bonusFire] = monster.ResFire
	resistances[bonusWater] = monster.ResWater
	resistances[bonusEarth] = monster.ResEarth
	resistances[bonusAir] = monster.ResAir
	weights := [bonusStats]float64{}
	total := 0.0
	for _, element := range damageElements {
		weighted := float64(attack[element]) * resistanceFactor(resistances[element])
		total += weighted
		weights[element] = weighted / 100
	}
	weights[bonusDmg] = total / 100
	return weights
}

func initialSelection(pools []bonusPool, weights [bonusStats]float64) [][]int {
	selection := make([][]int, len(pools))
	for index, pool := range pools {
		selection[index] = chooseItems(pool, weights)
	}
	return selection
}

func chooseItems(pool bonusPool, weights [bonusStats]float64) []int {
	if len(pool.items) == 0 {
		return nil
	}
	ranked := make([]int, len(pool.items))
	for index := range ranked {
		ranked[index] = index
	}
	slices.SortFunc(ranked, func(left, right int) int {
		leftScore := weightedScore(pool.items[left].bonus, weights)
		rightScore := weightedScore(pool.items[right].bonus, weights)
		return cmp.Compare(rightScore, leftScore)
	})
	return slices.Clone(ranked[:min(pool.rule.count, len(ranked))])
}

func weightedScore(item bonus, weights [bonusStats]float64) float64 {
	total := 0.0
	for stat := 0; stat < bonusStats; stat++ {
		total += float64(item[stat]) * weights[stat]
	}
	return total
}

func refineSelection(attack bonus, pools []bonusPool, monster schemas.MonsterSchema, selection [][]int) int {
	best := selectionDamage(attack, pools, monster, selection)
	for pass := 0; pass < 2; pass++ {
		improved := false
		for poolIndex, pool := range pools {
			for slot := 0; slot < len(selection[poolIndex]); slot++ {
				for itemIndex := range pool.items {
					if containsItem(selection[poolIndex], slot, itemIndex) {
						continue
					}
					previous := selection[poolIndex][slot]
					selection[poolIndex][slot] = itemIndex
					value := selectionDamage(attack, pools, monster, selection)
					if value > best {
						best = value
						improved = true
					} else {
						selection[poolIndex][slot] = previous
					}
				}
			}
		}
		if !improved {
			break
		}
	}
	return best
}

func containsItem(selection []int, skip, value int) bool {
	for index, item := range selection {
		if index != skip && item == value {
			return true
		}
	}
	return false
}

func selectionDamage(attack bonus, pools []bonusPool, monster schemas.MonsterSchema, selection [][]int) int {
	return damage(attack, selectionBonus(pools, selection), monster)
}

// selectionBonus sums the bonuses of the items picked by a selection.
func selectionBonus(pools []bonusPool, selection [][]int) bonus {
	gear := bonus{}
	for poolIndex, pool := range pools {
		for _, itemIndex := range selection[poolIndex] {
			for stat := 0; stat < bonusStats; stat++ {
				gear[stat] += pool.items[itemIndex].bonus[stat]
			}
		}
	}
	return gear
}
