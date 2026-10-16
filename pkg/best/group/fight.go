package group

import (
	"fmt"
	"math/rand"
	"sort"

	best "github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/catalog"
	groupfight "github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/schemas"
)

const (
	candidateIterations = 10
	finalIterations     = 50
)

var slots = []string{
	"weapon",
	"shield",
	"helmet",
	"body_armor",
	"leg_armor",
	"boots",
	"amulet",
	"rune",
	"ring1",
	"ring2",
	"artifact1",
	"artifact2",
	"artifact3",
}

type Result struct {
	Winrate        float32   `json:"winrate"`
	AverageTurns   float32   `json:"average_turns"`
	AverageFinalHP float32   `json:"average_final_hp"`
	Loadouts       []Loadout `json:"loadouts"`
}

type Loadout struct {
	Name      string             `json:"name"`
	Equipment map[string]string  `json:"equipment"`
	Utilities map[string]Utility `json:"utilities,omitempty"`
}

type Utility struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}

type score struct {
	winrate float32
	turns   float32
	hp      float32
	damage  float32
}

type allocation struct {
	equipment []map[string]string
	utilities []map[string]Utility
}

type roleChange struct {
	character int
	slot      string
	code      string
	tank      bool
}

func FindFightWithAvailable(party []schemas.CharacterSchema, monster schemas.MonsterSchema, available map[string]int) (Result, error) {
	if len(party) < 1 || len(party) > 3 {
		return Result{}, fmt.Errorf("party must contain 1 to 3 characters")
	}
	if !validNames(party) {
		return Result{}, fmt.Errorf("party character names must be distinct")
	}
	current := allocation{
		equipment: make([]map[string]string, len(party)),
		utilities: make([]map[string]Utility, len(party)),
	}
	for i, character := range party {
		current.equipment[i] = characterLoadout(character)
		current.utilities[i] = map[string]Utility{}
	}
	if !validAllocation(current, available) {
		return Result{}, fmt.Errorf("party loadouts exceed available items")
	}
	bestAllocation := current
	bestScore := score{winrate: -1}
	equipmentCandidates := []allocation{}
	for _, seed := range roleSeeds(party, current, available) {
		candidate, candidateScore := improveEquipment(party, monster, seed, available)
		equipmentCandidates = append(equipmentCandidates, candidate)
		if better(candidateScore, bestScore) {
			bestAllocation, bestScore = candidate, candidateScore
		}
	}
	if bestScore.winrate < 100 {
		for _, candidate := range uniqueSeeds(equipmentCandidates) {
			candidate, candidateScore := improveUtilities(party, monster, candidate, available)
			if better(candidateScore, bestScore) {
				bestAllocation, bestScore = candidate, candidateScore
			}
		}
	}
	final := evaluate(party, monster, bestAllocation, finalIterations)
	return result(party, bestAllocation, final), nil
}

func improveEquipment(party []schemas.CharacterSchema, monster schemas.MonsterSchema, current allocation, available map[string]int) (allocation, score) {
	bestScore := evaluate(party, monster, current, candidateIterations)
	for range 2 {
		for character := range party {
			for _, slot := range slots {
				for _, code := range equipmentOptions(party[character], slot, available) {
					candidate := copyAllocation(current)
					candidate.equipment[character][slot] = code
					if !validAllocation(candidate, available) {
						continue
					}
					candidateScore := evaluate(party, monster, candidate, candidateIterations)
					if better(candidateScore, bestScore) {
						current, bestScore = candidate, candidateScore
					}
				}
			}
		}
	}
	return current, bestScore
}

func improveUtilities(party []schemas.CharacterSchema, monster schemas.MonsterSchema, current allocation, available map[string]int) (allocation, score) {
	bestScore := evaluate(party, monster, current, candidateIterations)
	for range 2 {
		for character := range party {
			for _, slot := range []string{"utility1", "utility2"} {
				for _, utility := range utilityOptions(party[character], available) {
					candidate := copyAllocation(current)
					candidate.utilities[character][slot] = utility
					if !validAllocation(candidate, available) {
						continue
					}
					candidateScore := evaluate(party, monster, candidate, candidateIterations)
					if better(candidateScore, bestScore) {
						current, bestScore = candidate, candidateScore
					}
				}
			}
		}
	}
	return current, bestScore
}

func validNames(party []schemas.CharacterSchema) bool {
	seen := map[string]bool{}
	for _, character := range party {
		if character.Name == "" || seen[character.Name] {
			return false
		}
		seen[character.Name] = true
	}
	return true
}

func characterLoadout(character schemas.CharacterSchema) map[string]string {
	return map[string]string{
		"weapon":     character.WeaponSlot,
		"shield":     character.ShieldSlot,
		"helmet":     character.HelmetSlot,
		"body_armor": character.BodyArmorSlot,
		"leg_armor":  character.LegArmorSlot,
		"boots":      character.BootsSlot,
		"amulet":     character.AmuletSlot,
		"rune":       character.RuneSlot,
		"ring1":      character.Ring1Slot,
		"ring2":      character.Ring2Slot,
		"artifact1":  character.Artifact1Slot,
		"artifact2":  character.Artifact2Slot,
		"artifact3":  character.Artifact3Slot,
	}
}

func equipmentOptions(character schemas.CharacterSchema, slot string, available map[string]int) []string {
	result := []string{}
	if slot != "weapon" {
		result = append(result, "")
	}
	for code, quantity := range available {
		item, found := catalog.Items().Get(code)
		if quantity < 1 || !found || item.Type == "utility" || !best.CanEquip(character, *item) {
			continue
		}
		if item.Type == catalog.EquipmentSlotToTypes[slot] {
			result = append(result, code)
		}
	}
	sort.Strings(result)
	return result
}

func roleSeeds(party []schemas.CharacterSchema, current allocation, available map[string]int) []allocation {
	changes := []roleChange{}
	for character := range party {
		for _, slot := range slots {
			for _, code := range equipmentOptions(party[character], slot, available) {
				item, found := catalog.Items().Get(code)
				if !found || item.Effects == nil {
					continue
				}
				change := roleChange{
					character: character,
					slot:      slot,
					code:      code,
				}
				for _, effect := range *item.Effects {
					if effect.Code == "threat" && effect.Value > 0 {
						change.tank = true
					}
					if slot == "rune" && isSupportEffect(effect.Code) {
						changes = append(changes, change)
						break
					}
				}
				if change.tank {
					changes = append(changes, change)
				}
			}
		}
	}
	seeds := []allocation{current}
	for _, change := range changes {
		seed := applyRoleChanges(current, change)
		if validAllocation(seed, available) {
			seeds = append(seeds, seed)
		}
	}
	for _, tank := range changes {
		if !tank.tank {
			continue
		}
		for _, support := range changes {
			if support.tank || tank.character == support.character {
				continue
			}
			seed := applyRoleChanges(current, tank, support)
			if validAllocation(seed, available) {
				seeds = append(seeds, seed)
			}
		}
	}
	return uniqueSeeds(seeds)
}

func isSupportEffect(code string) bool {
	return code == "guard" || code == "healing_aura" || code == "vampiric_strike"
}

func applyRoleChanges(current allocation, changes ...roleChange) allocation {
	result := copyAllocation(current)
	for _, change := range changes {
		result.equipment[change.character][change.slot] = change.code
	}
	return result
}

func uniqueSeeds(seeds []allocation) []allocation {
	seen := map[string]bool{}
	result := []allocation{}
	for _, seed := range seeds {
		key := allocationKey(seed)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, seed)
	}
	return result
}

func allocationKey(current allocation) string {
	result := ""
	for i := range current.equipment {
		for _, slot := range slots {
			result += fmt.Sprintf("%d:%s=%s,", i, slot, current.equipment[i][slot])
		}
		for _, slot := range []string{"utility1", "utility2"} {
			utility := current.utilities[i][slot]
			result += fmt.Sprintf("%d:%s=%s:%d,", i, slot, utility.Code, utility.Quantity)
		}
	}
	return result
}

func utilityOptions(character schemas.CharacterSchema, available map[string]int) []Utility {
	result := []Utility{{}}
	codes := make([]string, 0)
	for code, quantity := range available {
		item, found := catalog.Items().Get(code)
		if quantity > 0 && found && item.Type == "utility" && best.CanEquip(character, *item) {
			codes = append(codes, code)
		}
	}
	sort.Strings(codes)
	for _, code := range codes {
		result = append(result, Utility{Code: code, Quantity: 100})
	}
	return result
}

func validLoadouts(party []schemas.CharacterSchema, loadouts []map[string]string, available map[string]int) bool {
	utilities := make([]map[string]Utility, len(party))
	for i, character := range party {
		utilities[i] = utilityLoadout(character)
	}
	return validAllocation(allocation{equipment: loadouts, utilities: utilities}, available)
}

func validAllocation(current allocation, available map[string]int) bool {
	used := map[string]int{}
	for i, loadout := range current.equipment {
		artifacts := map[string]bool{}
		adept := 0
		for _, slot := range slots {
			code := loadout[slot]
			if code == "" {
				continue
			}
			used[code]++
			if used[code] > available[code] {
				return false
			}
			if slot == "artifact1" || slot == "artifact2" || slot == "artifact3" {
				if artifacts[code] {
					return false
				}
				artifacts[code] = true
			}
			if code == "ring_of_the_adept" {
				adept++
			}
		}
		if adept > 1 {
			return false
		}
		seenUtilities := map[string]bool{}
		for _, utility := range current.utilities[i] {
			if utility.Code == "" {
				continue
			}
			if seenUtilities[utility.Code] {
				return false
			}
			seenUtilities[utility.Code] = true
		}
	}
	return true
}

func evaluate(party []schemas.CharacterSchema, monster schemas.MonsterSchema, current allocation, iterations int) score {
	fighters := make([]groupfight.Fighter, len(party))
	damage := float32(0)
	for i, character := range party {
		fighters[i] = groupfight.FighterFromFakeCharacter(i, fakeCharacter(character, current.equipment[i], current.utilities[i]))
		damage += effectiveDamage(fighters[i], monster)
	}
	summary := groupfight.SimulateMany(fighters, groupfight.MonsterFromSchema(monster), groupfight.SimulationOptions{
		Iterations: iterations,
		RNG:        rand.New(rand.NewSource(1)).Float64,
	})
	metrics := groupfight.Metrics(summary.Results)
	return score{
		winrate: summary.Winrate,
		turns:   metrics.AverageTurns,
		hp:      metrics.AverageFinalHP,
		damage:  damage,
	}
}

func fakeCharacter(character schemas.CharacterSchema, loadout map[string]string, utilities map[string]Utility) schemas.FakeCharacterSchema {
	result := schemas.FakeCharacterSchema{Level: character.Level}
	result.WeaponSlot = stringPointer(loadout["weapon"])
	result.ShieldSlot = stringPointer(loadout["shield"])
	result.HelmetSlot = stringPointer(loadout["helmet"])
	result.BodyArmorSlot = stringPointer(loadout["body_armor"])
	result.LegArmorSlot = stringPointer(loadout["leg_armor"])
	result.BootsSlot = stringPointer(loadout["boots"])
	result.AmuletSlot = stringPointer(loadout["amulet"])
	result.RuneSlot = stringPointer(loadout["rune"])
	result.Ring1Slot = stringPointer(loadout["ring1"])
	result.Ring2Slot = stringPointer(loadout["ring2"])
	result.Artifact1Slot = stringPointer(loadout["artifact1"])
	result.Artifact2Slot = stringPointer(loadout["artifact2"])
	result.Artifact3Slot = stringPointer(loadout["artifact3"])
	result.Utility1Slot, result.Utility1SlotQuantity = utilityPointers(utilities["utility1"])
	result.Utility2Slot, result.Utility2SlotQuantity = utilityPointers(utilities["utility2"])
	return result
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func utilityPointers(utility Utility) (*string, *int) {
	if utility.Code == "" || utility.Quantity < 1 {
		return nil, nil
	}
	return &utility.Code, &utility.Quantity
}

func utilityLoadout(character schemas.CharacterSchema) map[string]Utility {
	result := map[string]Utility{}
	if character.Utility1Slot != "" && character.Utility1SlotQuantity > 0 {
		result["utility1"] = Utility{
			Code:     character.Utility1Slot,
			Quantity: character.Utility1SlotQuantity,
		}
	}
	if character.Utility2Slot != "" && character.Utility2SlotQuantity > 0 {
		result["utility2"] = Utility{
			Code:     character.Utility2Slot,
			Quantity: character.Utility2SlotQuantity,
		}
	}
	return result
}

func copyAllocation(source allocation) allocation {
	result := allocation{
		equipment: make([]map[string]string, len(source.equipment)),
		utilities: make([]map[string]Utility, len(source.utilities)),
	}
	for i, loadout := range source.equipment {
		result.equipment[i] = map[string]string{}
		for slot, code := range loadout {
			result.equipment[i][slot] = code
		}
		result.utilities[i] = map[string]Utility{}
		for slot, utility := range source.utilities[i] {
			result.utilities[i][slot] = utility
		}
	}
	return result
}

func better(a, b score) bool {
	if a.winrate != b.winrate {
		return a.winrate > b.winrate
	}
	if a.winrate > 0 && a.turns != b.turns {
		return a.turns < b.turns
	}
	if a.winrate == 0 && a.damage != b.damage {
		return a.damage > b.damage
	}
	return a.hp > b.hp
}

func effectiveDamage(fighter groupfight.Fighter, monster schemas.MonsterSchema) float32 {
	attacks := []int{
		fighter.Stats.AttackFire,
		fighter.Stats.AttackEarth,
		fighter.Stats.AttackWater,
		fighter.Stats.AttackAir,
	}
	bonuses := []int{
		fighter.Stats.DmgFire,
		fighter.Stats.DmgEarth,
		fighter.Stats.DmgWater,
		fighter.Stats.DmgAir,
	}
	resistances := []int{
		monster.ResFire,
		monster.ResEarth,
		monster.ResWater,
		monster.ResAir,
	}
	result := float32(0)
	for i, attack := range attacks {
		multiplier := 1 - float32(resistances[i])/100
		if multiplier < 0 {
			multiplier = 0
		}
		result += float32(attack) * (1 + float32(fighter.Stats.Dmg+bonuses[i])/100) * multiplier
	}
	return result
}

func result(party []schemas.CharacterSchema, current allocation, value score) Result {
	loadouts := make([]Loadout, len(party))
	for i, character := range party {
		loadouts[i] = Loadout{
			Name:      character.Name,
			Equipment: current.equipment[i],
			Utilities: current.utilities[i],
		}
	}
	return Result{
		Winrate:        value.winrate,
		AverageTurns:   value.turns,
		AverageFinalHP: value.hp,
		Loadouts:       loadouts,
	}
}
