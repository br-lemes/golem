package damage

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
)

// Options controls which items can be selected.
type Options struct {
	// Available limits items to the given quantities. A nil map allows every
	// item in the catalog.
	Available map[string]int
	// Without excludes item codes from the selection.
	Without map[string]bool
}

// Weapon is a combat weapon that survives the dominance pruning together with
// the damage range it can reach at the character level.
type Weapon struct {
	Code         string        `json:"code"`
	Min          int           `json:"min_damage"`
	Max          int           `json:"max_damage"`
	Critical     int           `json:"critical"`
	Combinations []Combination `json:"combinations,omitempty"`
}

// Bonuses are the damage modifiers of one pruning scenario.
type Bonuses struct {
	Dmg      int `json:"dmg"`
	Fire     int `json:"fire"`
	Water    int `json:"water"`
	Earth    int `json:"earth"`
	Air      int `json:"air"`
	Critical int `json:"critical"`
}

// Combination is one bonus scenario and the damage the weapon reaches on it.
type Combination struct {
	Bonuses  Bonuses  `json:"bonuses"`
	Damage   int      `json:"damage"`
	Critical int      `json:"critical"`
	Items    []string `json:"items"`
}

type weapon struct {
	code     string
	attack   bonus
	critical int
}

// FindWeapons ranks combat weapons by the damage they reach with the best
// level-legal gear for each of them. Weapons beaten by another weapon in every
// possible loadout are removed first. A nil Available map allows every item in
// the catalog; otherwise only items with a positive quantity are considered,
// and repeatable slots never use more copies than the account owns. Items in
// Without are always excluded.
func FindWeapons(monster schemas.MonsterSchema, character schemas.CharacterSchema, options Options) []Weapon {
	pools := equipmentPools(character, options)
	corners := bonusCorners(equipmentBonuses(pools))
	survivors := pruneWeapons(findWeapons(character, options), corners, monster)
	result := make([]Weapon, 0, len(survivors))
	for _, survivor := range survivors {
		result = append(result, Weapon{
			Code:     survivor.code,
			Min:      damage(survivor.attack, bonus{}, monster),
			Max:      maxDamage(survivor.attack, pools, monster),
			Critical: survivor.critical,
		})
	}
	slices.SortFunc(result, func(left, right Weapon) int {
		if left.Max != right.Max {
			return cmp.Compare(right.Max, left.Max)
		}
		if left.Min != right.Min {
			return cmp.Compare(right.Min, left.Min)
		}
		return cmp.Compare(left.Code, right.Code)
	})
	return result
}

// FindWeapon returns every bonus scenario for one combat weapon. Weapons that
// the dominance pruning would remove are rejected, because their numbers would
// never be used by the ranking. A nil Available map allows every item in the
// catalog, and items in Without are always excluded.
func FindWeapon(monster schemas.MonsterSchema, character schemas.CharacterSchema, options Options, code string) (Weapon, error) {
	item, found := catalog.Items().Get(code)
	if !found || item.Type != "weapon" || item.Subtype != "" {
		return Weapon{}, fmt.Errorf("%s is not a combat weapon", code)
	}
	if options.Without[code] {
		return Weapon{}, fmt.Errorf("weapon %s is excluded", code)
	}
	if options.Available != nil && options.Available[code] < 1 {
		return Weapon{}, fmt.Errorf("weapon %s is not available in the account", code)
	}
	pools := equipmentPools(character, options)
	candidates := findWeapons(character, options)
	index := slices.IndexFunc(candidates, func(candidate weapon) bool {
		return candidate.code == code
	})
	if index == -1 {
		return Weapon{}, fmt.Errorf("weapon %s is not available at level %d", code, character.Level)
	}
	candidate := candidates[index]
	corners := bonusCorners(equipmentBonuses(pools))
	survivors := pruneWeapons(candidates, corners, monster)
	if !slices.ContainsFunc(survivors, func(survivor weapon) bool {
		return survivor.code == candidate.code
	}) {
		return Weapon{}, fmt.Errorf("weapon %s is dominated and pruned at level %d", code, character.Level)
	}
	return Weapon{
		Code:         candidate.code,
		Min:          damage(candidate.attack, bonus{}, monster),
		Max:          maxDamage(candidate.attack, pools, monster),
		Critical:     candidate.critical,
		Combinations: weaponCombinations(candidate, pools, monster),
	}, nil
}

// weaponCombinations lists achievable loadouts for one weapon, one per subset
// of the modifiers that change its damage. Each entry is the best gear that can
// be built while focusing on that subset, so the full subset matches the
// reported maximum and the empty subset the minimum.
func weaponCombinations(candidate weapon, pools []bonusPool, monster schemas.MonsterSchema) []Combination {
	stats := []int{bonusDmg}
	for _, element := range damageElements {
		if candidate.attack[element] > 0 {
			stats = append(stats, element)
		}
	}
	full := damageWeights(candidate.attack, monster)
	seen := make(map[string]bool, 1<<len(stats))
	result := make([]Combination, 0, 1<<len(stats))
	for mask := 0; mask < 1<<len(stats); mask++ {
		selection := make([][]int, len(pools))
		if mask > 0 {
			weights := [bonusStats]float64{}
			for bit, stat := range stats {
				if mask&(1<<bit) != 0 {
					weights[stat] = full[stat]
				}
			}
			selection = initialSelection(pools, weights)
			if mask == (1<<len(stats))-1 {
				selection = bestSelection(candidate.attack, pools, monster)
			}
		}
		gear := selectionBonus(pools, selection)
		exported := gear.export()
		items := []string{candidate.code}
		for poolIndex, pool := range pools {
			for _, itemIndex := range selection[poolIndex] {
				items = append(items, pool.items[itemIndex].code)
			}
		}
		key := fmt.Sprintf("%d,%d,%d,%d,%d,%d-%v", exported.Dmg, exported.Fire, exported.Water, exported.Earth, exported.Air, exported.Critical, items)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, Combination{
			Bonuses:  exported,
			Damage:   damage(candidate.attack, gear, monster),
			Critical: candidate.critical + gear[bonusCritical],
			Items:    items,
		})
	}
	result = pruneCombinations(result)
	slices.SortStableFunc(result, func(left, right Combination) int {
		if left.Damage != right.Damage {
			return cmp.Compare(right.Damage, left.Damage)
		}
		return cmp.Compare(right.Critical, left.Critical)
	})
	return result
}

// pruneCombinations removes loadouts dominated in both damage and critical.
// Loadouts that are better in only one of them are kept as valid trade-offs.
func pruneCombinations(combinations []Combination) []Combination {
	result := make([]Combination, 0, len(combinations))
	for index, combination := range combinations {
		if !dominatedCombination(combinations, index) {
			result = append(result, combination)
		}
	}
	return result
}

func dominatedCombination(combinations []Combination, index int) bool {
	for other := range combinations {
		if other == index {
			continue
		}
		if dominatesCombination(combinations[other], combinations[index]) {
			return true
		}
	}
	return false
}

func dominatesCombination(better, worse Combination) bool {
	if better.Damage < worse.Damage || better.Critical < worse.Critical {
		return false
	}
	return better.Damage > worse.Damage || better.Critical > worse.Critical
}

func findWeapons(character schemas.CharacterSchema, options Options) []weapon {
	items := catalog.Items().Filter(func(item *schemas.ItemSchema) bool {
		if item.Type != "weapon" || item.Subtype != "" {
			return false
		}
		if options.Without[item.Code] {
			return false
		}
		if options.Available != nil && options.Available[item.Code] < 1 {
			return false
		}
		return utils.MeetsItemConditions(character, *item)
	})
	result := make([]weapon, 0, len(items))
	for _, item := range items {
		result = append(result, newWeapon(item))
	}
	return result
}

func newWeapon(item *schemas.ItemSchema) weapon {
	result := weapon{code: item.Code}
	if item.Effects == nil {
		return result
	}
	for _, effect := range *item.Effects {
		switch effect.Code {
		case "attack_fire":
			result.attack[bonusFire] += effect.Value
		case "attack_water":
			result.attack[bonusWater] += effect.Value
		case "attack_earth":
			result.attack[bonusEarth] += effect.Value
		case "attack_air":
			result.attack[bonusAir] += effect.Value
		case "critical_strike":
			result.critical += effect.Value
		}
	}
	return result
}

// bonusCorners returns every point of the bonus box where each modifier is
// either zero or its achievable maximum.
func bonusCorners(maximum bonus) []bonus {
	stats := make([]int, bonusStats)
	for stat := range stats {
		stats[stat] = stat
	}
	return statCorners(maximum, stats)
}

func statCorners(maximum bonus, stats []int) []bonus {
	corners := []bonus{{}}
	for _, stat := range stats {
		next := make([]bonus, 0, len(corners)*2)
		for _, corner := range corners {
			next = append(next, corner)
			high := corner
			high[stat] = maximum[stat]
			next = append(next, high)
		}
		corners = next
	}
	return corners
}

func (b bonus) export() Bonuses {
	return Bonuses{
		Dmg:      b[bonusDmg],
		Fire:     b[bonusFire],
		Water:    b[bonusWater],
		Earth:    b[bonusEarth],
		Air:      b[bonusAir],
		Critical: b[bonusCritical],
	}
}

type weaponAttributes struct {
	damage   int
	critical int
}

func pruneWeapons(weapons []weapon, corners []bonus, monster schemas.MonsterSchema) []weapon {
	values := make([][]weaponAttributes, len(weapons))
	for index, candidate := range weapons {
		values[index] = make([]weaponAttributes, len(corners))
		for cornerIndex, corner := range corners {
			values[index][cornerIndex] = weaponValue(candidate, corner, monster)
		}
	}
	result := make([]weapon, 0, len(weapons))
	for index, candidate := range weapons {
		if !dominatedWeapon(values, index) {
			result = append(result, candidate)
		}
	}
	return result
}

func dominatedWeapon(values [][]weaponAttributes, index int) bool {
	for other := range values {
		if other == index {
			continue
		}
		if dominatesWeapon(values, other, index) {
			return true
		}
	}
	return false
}

func dominatesWeapon(values [][]weaponAttributes, better, worse int) bool {
	strict := false
	for corner := range values[worse] {
		best := values[better][corner]
		worst := values[worse][corner]
		if best.damage < worst.damage || best.critical < worst.critical {
			return false
		}
		if best.damage > worst.damage || best.critical > worst.critical {
			strict = true
		}
	}
	return strict
}

func weaponValue(candidate weapon, gear bonus, monster schemas.MonsterSchema) weaponAttributes {
	return weaponAttributes{
		damage:   damage(candidate.attack, gear, monster),
		critical: candidate.critical + gear[bonusCritical],
	}
}
