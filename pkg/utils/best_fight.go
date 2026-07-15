package utils

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

// weapon must remain first to fix the element the fight scales around
var allSlots = []string{
	"weapon", "shield", "helmet", "body_armor", "leg_armor", "boots", "amulet", "rune", "bag",
	"ring1", "ring2", "artifact1", "artifact2", "artifact3", "utility1", "utility2",
}

var (
	singleSlots   = []string{"shield", "helmet", "body_armor", "leg_armor", "boots", "amulet", "rune", "bag"}
	ringSlots     = []string{"ring1", "ring2"}
	artifactSlots = []string{"artifact1", "artifact2", "artifact3"}
	utilitySlots  = []string{"utility1", "utility2"}
)

const (
	perSlotCap = 12
	weaponCap  = 10
	groupCap   = 12
)

func BestFight(character schemas.CharacterSchema, monster schemas.MonsterSchema) (map[string]bestResult, error) {
	owned, err := ownedItemCounts(character)
	if err != nil {
		return nil, err
	}
	base := baseStats(character)
	pools := buildFightPools(character, monster, owned)
	worn := equippedSlotMap(character)
	evaluate := makeFightEvaluator(base, monster, worn)

	weaponOpts := map[string]bool{"": true}
	for _, w := range pools["weapon"] {
		weaponOpts[w] = true
	}
	if worn["weapon"] != "" {
		weaponOpts[worn["weapon"]] = true
	}

	var bestSlots map[string]string
	bestScore := math.Inf(-1)

	for w := range weaponOpts {
		slots := make(map[string]string, len(allSlots))
		for _, s := range allSlots {
			slots[s] = ""
		}
		slots["weapon"] = w

		for pass := 0; pass < 3; pass++ {
			for _, s := range singleSlots {
				best := slots[s]
				_, bestSc := evaluate(slots)
				for _, opt := range append([]string{""}, pools[s]...) {
					if opt == slots[s] {
						continue
					}
					slots[s] = opt
					_, sc := evaluate(slots)
					if sc > bestSc {
						bestSc = sc
						best = opt
					}
				}
				slots[s] = best
			}
			chooseGroup(slots, ringSlots, pools["ring1"], owned, worn, evaluate)
			chooseGroup(slots, artifactSlots, pools["artifact1"], nil, worn, evaluate)
			chooseGroup(slots, utilitySlots, pools["utility1"], nil, worn, evaluate)
		}

		_, sc := evaluate(slots)
		if sc > bestScore {
			bestScore = sc
			bestSlots = make(map[string]string, len(slots))
			for k, v := range slots {
				bestSlots[k] = v
			}
		}
	}

	result := make(map[string]bestResult)
	for _, s := range allSlots {
		code := bestSlots[s]
		if code == "" || code == worn[s] {
			continue
		}
		it, ok := database.GetItem(code)
		if !ok {
			continue
		}
		result[s] = bestResult{Code: code, Value: formatAllEffects(it)}
	}
	return result, nil
}

func ownedItemCounts(character schemas.CharacterSchema) (map[string]int, error) {
	bankItems, err := api.MyBankItems()
	if err != nil {
		return nil, err
	}
	owned := make(map[string]int)
	for _, item := range bankItems {
		owned[item.Code] += item.Quantity
	}
	if character.Inventory != nil {
		for _, item := range *character.Inventory {
			owned[item.Code] += item.Quantity
		}
	}
	for slot, code := range equippedSlotMap(character) {
		if code == "" {
			continue
		}
		qty := 1
		switch slot {
		case "utility1":
			qty = character.Utility1SlotQuantity
		case "utility2":
			qty = character.Utility2SlotQuantity
		}
		owned[code] += qty
	}
	return owned, nil
}

func buildFightPools(character schemas.CharacterSchema, monster schemas.MonsterSchema, owned map[string]int) map[string][]string {
	pools := make(map[string][]string, len(allSlots))
	seen := make(map[string]map[string]bool, len(allSlots))
	for _, s := range allSlots {
		seen[s] = make(map[string]bool)
	}

	for code, qty := range owned {
		if qty <= 0 {
			continue
		}
		it, ok := database.GetItem(code)
		if !ok {
			continue
		}
		slots, ok := database.EquipmentTypeToSlots[it.Type]
		if !ok {
			continue
		}
		if !canEquipItem(character, it) {
			continue
		}
		for _, s := range slots {
			if !seen[s][code] {
				seen[s][code] = true
				pools[s] = append(pools[s], code)
			}
		}
	}

	for _, s := range allSlots {
		limit := perSlotCap
		if s == "weapon" {
			limit = weaponCap
		}
		pools[s] = capPool(pools[s], limit, monster)
	}
	return pools
}

func capPool(codes []string, cap int, monster schemas.MonsterSchema) []string {
	if len(codes) <= cap {
		out := append([]string{}, codes...)
		sort.Strings(out)
		return out
	}
	ranked := append([]string{}, codes...)
	sort.Slice(ranked, func(i, j int) bool {
		ii, _ := database.GetItem(ranked[i])
		jj, _ := database.GetItem(ranked[j])
		hi, hj := fightHeuristic(ii, monster), fightHeuristic(jj, monster)
		if hi != hj {
			return hi > hj
		}
		return ranked[i] < ranked[j]
	})
	out := append([]string{}, ranked[:cap]...)
	sort.Strings(out)
	return out
}

func canEquipItem(c schemas.CharacterSchema, it schemas.ItemSchema) bool {
	if it.Conditions == nil {
		return true
	}
	for _, cond := range *it.Conditions {
		stat := characterSkillLevel(c, cond.Code)
		ok := true
		switch cond.Operator {
		case schemas.Gt:
			ok = stat > cond.Value
		case schemas.Lt:
			ok = stat < cond.Value
		case schemas.Eq:
			ok = stat == cond.Value
		case schemas.Ne:
			ok = stat != cond.Value
		}
		if !ok {
			return false
		}
	}
	return true
}

func characterSkillLevel(c schemas.CharacterSchema, code string) int {
	switch code {
	case "level":
		return c.Level
	case "mining":
		return c.MiningLevel
	case "woodcutting":
		return c.WoodcuttingLevel
	case "fishing":
		return c.FishingLevel
	case "weaponcrafting":
		return c.WeaponcraftingLevel
	case "gearcrafting":
		return c.GearcraftingLevel
	case "jewelrycrafting":
		return c.JewelrycraftingLevel
	case "cooking":
		return c.CookingLevel
	case "alchemy":
		return c.AlchemyLevel
	default:
		return 0
	}
}

func fightHeuristic(it schemas.ItemSchema, monster schemas.MonsterSchema) float64 {
	if it.Effects == nil {
		return 0
	}
	v := 0.0
	for _, e := range *it.Effects {
		val := float64(e.Value)
		switch e.Code {
		case "attack_fire":
			v += val * 2 * math.Max(0, 1-float64(monster.ResFire)/100)
		case "attack_earth":
			v += val * 2 * math.Max(0, 1-float64(monster.ResEarth)/100)
		case "attack_water":
			v += val * 2 * math.Max(0, 1-float64(monster.ResWater)/100)
		case "attack_air":
			v += val * 2 * math.Max(0, 1-float64(monster.ResAir)/100)
		case "dmg", "dmg_fire", "dmg_earth", "dmg_water", "dmg_air":
			v += val
		case "critical_strike":
			v += val * 0.5
		case "hp", "boost_hp":
			v += val * 0.3
		case "res_fire", "res_earth", "res_water", "res_air":
			v += val * 0.5
		case "boost_dmg_fire", "boost_dmg_earth", "boost_dmg_water", "boost_dmg_air":
			v += val
		case "restore":
			v += val * 0.3
		}
	}
	return v
}

type fightEvaluator func(slots map[string]string) (FightForecast, float64)

func makeFightEvaluator(base EffectiveStats, monster schemas.MonsterSchema, worn map[string]string) fightEvaluator {
	return func(slots map[string]string) (FightForecast, float64) {
		var codes []string
		for _, s := range allSlots {
			if slots[s] != "" {
				codes = append(codes, slots[s])
			}
		}
		fighter := applyGear(base, codes)
		f := Simulate(fighter, monster, SimOptions{})

		// Prefer filling slots, then prefer keeping currently worn items to avoid lateral swaps
		filled, keep := 0, 0
		for _, s := range allSlots {
			if slots[s] == "" {
				continue
			}
			filled++
			if slots[s] == worn[s] {
				keep++
			}
		}
		tiebreak := float64(filled)*1e-3 + float64(keep)*1e-4

		if !f.Win || f.TimedOut {
			return f, -1e12 + float64(f.Turns)*1000 + float64(f.HpRemaining) + tiebreak
		}
		worst := Simulate(fighter, monster, SimOptions{Pessimistic: true})
		score := float64(f.HpRemaining)*1000 - float64(f.Turns) + tiebreak
		if worst.Win {
			score += 1e9
		}
		return f, score
	}
}

func alignGroupToWorn(group []string, chosen []string, worn map[string]string) []string {
	var remaining []string
	for _, c := range chosen {
		if c != "" {
			remaining = append(remaining, c)
		}
	}
	result := make([]string, len(group))
	for i, slot := range group {
		w := worn[slot]
		if w == "" {
			continue
		}
		for j, c := range remaining {
			if c == w {
				result[i] = c
				remaining = append(remaining[:j], remaining[j+1:]...)
				break
			}
		}
	}
	ri := 0
	for i := range result {
		if result[i] != "" {
			continue
		}
		if ri < len(remaining) {
			result[i] = remaining[ri]
			ri++
		}
	}
	return result
}

func chooseGroup(slots map[string]string, group, pool []string, ownedQty map[string]int, worn map[string]string, evaluate fightEvaluator) {
	n := len(group)
	capped := pool
	if len(capped) > groupCap {
		capped = capped[:groupCap]
	}

	countOf := func(code string) int {
		if ownedQty == nil {
			return 1
		}
		q, ok := ownedQty[code]
		if ok {
			return q
		}
		return 3
	}

	best := make([]string, n)
	for i, g := range group {
		best[i] = slots[g]
	}
	_, bestScore := evaluate(slots)

	var combos [][]string
	var gen func(start int, chosen []string)
	gen = func(start int, chosen []string) {
		combo := append([]string{}, chosen...)
		combos = append(combos, combo)
		if len(chosen) == n {
			return
		}
		for i := start; i < len(capped); i++ {
			max := countOf(capped[i])
			if max > n {
				max = n
			}
			used := 0
			for _, c := range chosen {
				if c == capped[i] {
					used++
				}
			}
			if used >= max {
				continue
			}
			next := i
			if used+1 >= max {
				next = i + 1
			}
			gen(next, append(chosen, capped[i]))
		}
	}
	gen(0, nil)

	for _, combo := range combos {
		for i := 0; i < n; i++ {
			v := ""
			if i < len(combo) {
				v = combo[i]
			}
			slots[group[i]] = v
		}
		_, sc := evaluate(slots)
		if sc > bestScore {
			bestScore = sc
			best = make([]string, n)
			copy(best, combo)
		}
	}
	aligned := alignGroupToWorn(group, best, worn)
	for i := 0; i < n; i++ {
		slots[group[i]] = aligned[i]
	}
}

func formatAllEffects(it schemas.ItemSchema) string {
	if it.Effects == nil {
		return ""
	}
	var parts []string
	for _, e := range *it.Effects {
		format := "%+d %s"
		if e.Value < 0 {
			format = "%d %s"
		}
		parts = append(parts, fmt.Sprintf(format, e.Value, e.Code))
	}
	return strings.Join(parts, ", ")
}
