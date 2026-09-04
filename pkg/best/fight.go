package best

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/surplus"
)

// CACHE_VERSION / ACTION REQUIRED: increment when the fight optimizer or
// simulation can produce different results for the same inputs.
const bestFightSimulationCacheVersion = 2

type Result struct {
	Winrate               float32           `json:"winrate"`
	AverageTurns          float32           `json:"average_turns"`
	AverageFinalHP        float32           `json:"average_final_hp"`
	AverageFightCooldown  float32           `json:"-"`
	EstimatedRestCooldown float32           `json:"-"`
	CycleCost             float32           `json:"cycle_seconds"`
	DamageSurplus         float32           `json:"damage_surplus"`
	SurvivalSurplus       float32           `json:"survival_surplus"`
	XP                    int               `json:"xp"`
	XPPerSecond           float32           `json:"xp_per_second"`
	GoldPerSecond         float32           `json:"gold_per_second"`
	ProspectingEfficiency float32           `json:"prospecting_efficiency"`
	Safe                  bool              `json:"safe"`
	Wisdom                int               `json:"-"`
	Prospecting           int               `json:"-"`
	Focus                 string            `json:"focus"`
	FocusScore            float32           `json:"-"`
	ArtifactOrderScore    int               `json:"-"`
	Unowned               []string          `json:"unowned_equipment,omitempty"`
	ContextScore          int               `json:"-"`
	Equipment             map[string]string `json:"equipment"`
	Utilities             map[string]string `json:"utilities,omitempty"`
	// FinalEquipment is the complete final loadout. It is kept out of the
	// serialized result so existing commands continue to expose only changes.
	FinalEquipment map[string]string `json:"-"`
}

var simulationBestSlots = []string{
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

var simulationBestUtilitySlots = []string{"utility1", "utility2"}
var simulationRingSlots = []string{"ring1", "ring2"}
var simulationArtifactSlots = []string{"artifact1", "artifact2", "artifact3"}

const (
	candidateSimulationIterations = 10
	finalSimulationIterations     = 50
)

var processSimulationCache map[string]cachedSimulationResult

func ResetSimulationCache() {
	processSimulationCache = make(map[string]cachedSimulationResult)
}

type cachedSimulationResult struct {
	Winrate               float32
	AverageTurns          float32
	AverageFinalHP        float32
	AverageFightCooldown  float32
	EstimatedRestCooldown float32
	CycleCost             float32
	DamageSurplus         float32
	SurvivalSurplus       float32
	XP                    int
	XPPerSecond           float32
	GoldPerSecond         float32
	ProspectingEfficiency float32
	Safe                  bool
}

func findFight(character schemas.CharacterSchema, monster schemas.MonsterSchema, includeUnowned, allowDuplicateAdeptRing bool) (Result, error) {
	available, err := bankAvailable()
	if err != nil {
		return Result{}, err
	}
	if character.Inventory != nil {
		for _, item := range *character.Inventory {
			if item.Code != "" {
				available[item.Code] += item.Quantity
			}
		}
	}
	for _, code := range characterSlots(character) {
		if code != "" {
			available[code]++
		}
	}
	return FindFightWithAvailable(character, monster, available, includeUnowned, allowDuplicateAdeptRing)
}

func FindFightByName(name string, monster schemas.MonsterSchema, includeUnowned, allowDuplicateAdeptRing bool) (Result, error) {
	character, err := api.Characters(name)
	if err != nil {
		return Result{}, err
	}
	return findFight(character, monster, includeUnowned, allowDuplicateAdeptRing)
}

func FindFightAtLevel(level int, monster schemas.MonsterSchema, includeUnowned, allowDuplicateAdeptRing bool) (Result, error) {
	available, err := bankAvailable()
	if err != nil {
		return Result{}, err
	}
	character := schemas.CharacterSchema{Level: level}
	return FindFightWithAvailable(character, monster, available, includeUnowned, allowDuplicateAdeptRing)
}

func bankAvailable() (map[string]int, error) {
	owned, err := api.MyBankItems()
	if err != nil {
		return nil, err
	}
	available := make(map[string]int)
	for _, item := range owned {
		available[item.Code] += item.Quantity
	}
	return available, nil
}

func FindFightWithAvailable(character schemas.CharacterSchema, monster schemas.MonsterSchema, available map[string]int, includeUnowned, allowDuplicateAdeptRing bool) (Result, error) {
	if processSimulationCache == nil {
		processSimulationCache = make(map[string]cachedSimulationResult)
	}
	ownedCodes := make(map[string]bool, len(available))
	for code := range available {
		ownedCodes[code] = true
	}
	if includeUnowned {
		for _, item := range database.Items().All() {
			if item.Type == "tool" || item.Subtype == "tool" {
				continue
			}
			if available[item.Code] < 2 {
				available[item.Code] = 2
			}
		}
	}

	validItems := make([]schemas.ItemSchema, 0)
	for code, quantity := range available {
		if quantity < 1 {
			continue
		}
		item, exists := database.Items().Get(code)
		if !exists {
			continue
		}
		_, isEquipment := database.EquipmentTypeToSlots[item.Type]
		if !isEquipment || item.Type == "utility" || item.Subtype == "tool" {
			continue
		}
		if !CanEquip(character, *item) {
			continue
		}
		validItems = append(validItems, *item)
	}
	sort.Slice(validItems, func(i, j int) bool {
		return validItems[i].Code < validItems[j].Code
	})
	nonDominated := surplus.NonDominated(validItems, character)
	options := make(map[string][]string)
	for _, slot := range simulationBestSlots {
		options[slot] = []string{}
		if slot != "weapon" {
			options[slot] = append(options[slot], "")
		}
		for _, item := range nonDominated {
			if !itemFitsSlot(item, slot) {
				continue
			}
			options[slot] = append(options[slot], item.Code)
		}
		sort.Strings(options[slot])
	}

	// The initial equipment is not a search candidate. It is only used later
	// to format the recommendation and preserve equivalent artifact ordering.
	current := make(map[string]string)
	best := Result{Winrate: -1}
	var bestLoadout map[string]string
	combatCandidates := generateBeamCandidates(character, monster, options, available, allowDuplicateAdeptRing)
	finalists := combatCandidates
	preferredContext := "prospecting"
	if simulationContextUsesWisdom(character, monster) {
		preferredContext = "wisdom"
	}
	profiles := []string{"damage", "balanced", "survival"}
	profiles = append(profiles, preferredContext)
	for _, weapon := range options["weapon"] {
		for _, profile := range profiles {
			heuristic := heuristicSimulationLoadout(options, monster, available, weapon, profile)
			if !allowDuplicateAdeptRing && !adeptRingAllowed(character, heuristic) {
				if heuristic["ring1"] == "ring_of_the_adept" {
					heuristic["ring2"] = ""
				} else {
					heuristic["ring1"] = ""
				}
			}
			heuristicScore := evaluateSimulationLoadout(character, monster, heuristic)
			finalists = append(finalists, copyStringMap(heuristic))
			setSimulationContext(&heuristicScore, character, monster, heuristic)
			if betterSearchScore(heuristicScore, best) {
				current = heuristic
				best = heuristicScore
				bestLoadout = copyStringMap(heuristic)
			}
		}
	}
	for pass := 0; pass < 2; pass++ {
		for _, slot := range simulationBestSlots {
			if slot == "ring1" || slot == "ring2" || strings.HasPrefix(slot, "artifact") {
				continue
			}
			chosen := current[slot]
			for _, code := range options[slot] {
				candidate := copyStringMap(current)
				candidate[slot] = code
				if !loadoutHasQuantity(candidate, available) || (!allowDuplicateAdeptRing && !adeptRingAllowed(character, candidate)) {
					continue
				}
				score := evaluateSimulationLoadout(character, monster, candidate)
				setSimulationContext(&score, character, monster, candidate)
				if betterSearchScore(score, best) {
					best, chosen = score, code
					bestLoadout = copyStringMap(candidate)
				}
			}
			current[slot] = chosen
		}
		chooseSimulationGroup(character, monster, current, simulationRingSlots, options["ring1"], available, allowDuplicateAdeptRing, false)
		chooseSimulationGroup(character, monster, current, simulationArtifactSlots, options["artifact1"], available, allowDuplicateAdeptRing, true)
		current, best = keepBestSimulationLoadout(character, monster, current, bestLoadout, best)
		bestLoadout = copyStringMap(current)
	}
	if best.Winrate < 100 {
		current, best = refineWithUtilities(character, monster, current, bestLoadout, best, options, available, allowDuplicateAdeptRing, candidateSimulationIterations, betterSearchScore)
		bestLoadout = copyStringMap(current)
	}
	current = copyStringMap(bestLoadout)
	finalists = append(finalists, copyStringMap(current))
	finalists = uniqueSimulationLoadouts(finalists)
	bestLoadout, best = chooseFinalist(character, monster, finalists)
	current = copyStringMap(bestLoadout)
	if best.Winrate < 100 {
		current, best = refineWithUtilities(character, monster, current, bestLoadout, best, options, available, allowDuplicateAdeptRing, finalSimulationIterations, betterSimulationScore)
	}
	setSimulationContext(&best, character, monster, current)
	focus := "prospecting"
	if simulationContextUsesWisdom(character, monster) {
		focus = "wisdom"
	}
	current = alignSimulationArtifacts(character, current)
	original := characterSlots(character)
	equipment := simulationEquipmentChanges(current, original)
	utilities := simulationUtilityChanges(current, character)
	result := Result{
		Winrate:               best.Winrate,
		AverageTurns:          best.AverageTurns,
		AverageFinalHP:        best.AverageFinalHP,
		AverageFightCooldown:  best.AverageFightCooldown,
		EstimatedRestCooldown: best.EstimatedRestCooldown,
		CycleCost:             best.CycleCost,
		DamageSurplus:         best.DamageSurplus,
		SurvivalSurplus:       best.SurvivalSurplus,
		XP:                    best.XP,
		XPPerSecond:           best.XPPerSecond,
		GoldPerSecond:         best.GoldPerSecond,
		ProspectingEfficiency: best.ProspectingEfficiency,
		Safe:                  best.Safe,
		Wisdom:                best.Wisdom,
		Prospecting:           best.Prospecting,
		Focus:                 focus,
		Equipment:             equipment,
		Utilities:             utilities,
		FinalEquipment:        copyStringMap(current),
	}
	seenUnowned := make(map[string]bool)
	for _, code := range equipment {
		if code != "" && !ownedCodes[code] && !seenUnowned[code] {
			result.Unowned = append(result.Unowned, code)
			seenUnowned[code] = true
		}
	}
	for _, code := range utilities {
		if code != "" && !ownedCodes[code] && !seenUnowned[code] {
			result.Unowned = append(result.Unowned, code)
			seenUnowned[code] = true
		}
	}
	sort.Strings(result.Unowned)
	return result, nil
}

func simulationEquipmentChanges(current, original map[string]string) map[string]string {
	changes := make(map[string]string)
	for _, slot := range simulationBestSlots {
		if current[slot] != "" && current[slot] != original[slot] {
			changes[slot] = current[slot]
		}
	}
	return changes
}

func simulationUtilityChanges(current map[string]string, character schemas.CharacterSchema) map[string]string {
	original := map[string]string{
		"utility1": character.Utility1Slot,
		"utility2": character.Utility2Slot,
	}
	changes := make(map[string]string)
	for _, slot := range simulationBestUtilitySlots {
		if current[slot] != original[slot] {
			changes[slot] = current[slot]
		}
	}
	return changes
}

// alignSimulationArtifacts assigns the already-selected artifact set to
// slots while preserving the position of artifacts that are still equipped.
// Artifact slots are interchangeable for combat; this only makes the
// recommendation stable and avoids needless swaps.
func alignSimulationArtifacts(character schemas.CharacterSchema, slots map[string]string) map[string]string {
	result := copyStringMap(slots)
	original := characterSlots(character)
	artifactSlots := []string{"artifact1", "artifact2", "artifact3"}
	for _, slot := range artifactSlots {
		result[slot] = ""
	}
	remaining := make([]string, 0, len(artifactSlots))
	for _, slot := range artifactSlots {
		code := slots[slot]
		if code != "" {
			remaining = append(remaining, code)
		}
	}
	for _, slot := range artifactSlots {
		code := original[slot]
		for i, candidate := range remaining {
			if code != "" && candidate == code {
				result[slot] = candidate
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	for _, slot := range artifactSlots {
		if result[slot] != "" {
			continue
		}
		if len(remaining) > 0 {
			result[slot] = remaining[0]
			remaining = remaining[1:]
		} else {
			result[slot] = ""
		}
	}
	return result
}

func keepBestSimulationLoadout(character schemas.CharacterSchema, monster schemas.MonsterSchema, current, bestLoadout map[string]string, best Result) (map[string]string, Result) {
	currentScore := evaluateSimulationLoadout(character, monster, current)
	setSimulationContext(&currentScore, character, monster, current)
	if betterSearchScore(currentScore, best) {
		return current, currentScore
	}
	return copyStringMap(bestLoadout), best
}

func refineWithUtilities(character schemas.CharacterSchema, monster schemas.MonsterSchema, current, bestLoadout map[string]string, best Result, options map[string][]string, available map[string]int, allowDuplicateAdeptRing bool, iterations int, better func(Result, Result) bool) (map[string]string, Result) {
	for _, slot := range simulationBestUtilitySlots {
		options[slot] = []string{""}
		for code, quantity := range available {
			if quantity < 1 {
				continue
			}
			item, exists := database.Items().Get(code)
			if exists && item.Type == "utility" && CanEquip(character, *item) {
				options[slot] = append(options[slot], code)
			}
		}
		sort.Strings(options[slot])
	}
	for _, utility1 := range options["utility1"] {
		for _, utility2 := range options["utility2"] {
			candidate := copyStringMap(current)
			candidate["utility1"] = utility1
			candidate["utility2"] = utility2
			if !loadoutHasQuantity(candidate, available) || (!allowDuplicateAdeptRing && !adeptRingAllowed(character, candidate)) {
				continue
			}
			score := evaluateSimulationLoadoutIterations(character, monster, candidate, iterations)
			setSimulationContext(&score, character, monster, candidate)
			if better(score, best) {
				best = score
				bestLoadout = copyStringMap(candidate)
			}
		}
	}
	current = copyStringMap(bestLoadout)
	for pass := 0; pass < 2; pass++ {
		for _, slot := range simulationBestSlots {
			if slot == "ring1" || slot == "ring2" || strings.HasPrefix(slot, "artifact") {
				continue
			}
			chosen := current[slot]
			for _, code := range options[slot] {
				candidate := copyStringMap(current)
				candidate[slot] = code
				if !loadoutHasQuantity(candidate, available) || (!allowDuplicateAdeptRing && !adeptRingAllowed(character, candidate)) {
					continue
				}
				score := evaluateSimulationLoadoutIterations(character, monster, candidate, iterations)
				setSimulationContext(&score, character, monster, candidate)
				if better(score, best) {
					best, chosen = score, code
					bestLoadout = copyStringMap(candidate)
				}
			}
			current[slot] = chosen
		}
	}
	chooseSimulationGroup(character, monster, current, simulationRingSlots, options["ring1"], available, allowDuplicateAdeptRing, false)
	chooseSimulationGroup(character, monster, current, simulationArtifactSlots, options["artifact1"], available, allowDuplicateAdeptRing, true)
	current, best = keepBestSimulationLoadout(character, monster, current, bestLoadout, best)
	return copyStringMap(current), best
}

func chooseFinalist(character schemas.CharacterSchema, monster schemas.MonsterSchema, finalists []map[string]string) (map[string]string, Result) {
	var bestLoadout map[string]string
	var best Result
	for _, candidate := range finalists {
		score := evaluateSimulationLoadoutIterations(character, monster, candidate, finalSimulationIterations)
		setSimulationContext(&score, character, monster, candidate)
		if bestLoadout == nil || betterSimulationScore(score, best) {
			bestLoadout = copyStringMap(candidate)
			best = score
		}
	}
	return bestLoadout, best
}

func simulationLoadoutKey(slots map[string]string) string {
	keys := make([]string, 0, len(slots))
	for key := range slots {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, key+"="+slots[key])
	}
	return strings.Join(values, ",")
}

func uniqueSimulationLoadouts(loadouts []map[string]string) []map[string]string {
	seen := make(map[string]bool, len(loadouts))
	unique := make([]map[string]string, 0, len(loadouts))
	for _, loadout := range loadouts {
		key := simulationLoadoutKey(loadout)
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, loadout)
	}
	return unique
}

func generateBeamCandidates(character schemas.CharacterSchema, monster schemas.MonsterSchema, options map[string][]string, available map[string]int, allowDuplicateAdeptRing bool) []map[string]string {
	const beamWidth = 20
	searchSlots := []string{
		"shield",
		"helmet",
		"body_armor",
		"leg_armor",
		"boots",
		"amulet",
		"rune",
	}
	beam := make([]map[string]string, 0)
	for _, weapon := range options["weapon"] {
		candidate := map[string]string{"weapon": weapon}
		beam = append(beam, candidate)
	}
	for _, slot := range searchSlots {
		expanded := make([]map[string]string, 0, len(beam)*len(options[slot]))
		for _, base := range beam {
			for _, code := range options[slot] {
				candidate := copyStringMap(base)
				candidate[slot] = code
				if !loadoutHasQuantity(candidate, available) || (!allowDuplicateAdeptRing && !adeptRingAllowed(character, candidate)) {
					continue
				}
				expanded = append(expanded, candidate)
			}
		}
		sort.SliceStable(expanded, func(i, j int) bool {
			left := evaluateSimulationLoadout(character, monster, expanded[i])
			right := evaluateSimulationLoadout(character, monster, expanded[j])
			return betterSearchScore(left, right)
		})
		if len(expanded) > beamWidth {
			expanded = expanded[:beamWidth]
		}
		beam = expanded
	}
	for _, candidate := range beam {
		chooseSimulationGroup(character, monster, candidate, simulationRingSlots, options["ring1"], available, allowDuplicateAdeptRing, false)
		chooseSimulationGroup(character, monster, candidate, simulationArtifactSlots, options["artifact1"], available, allowDuplicateAdeptRing, true)
	}
	return beam
}

func heuristicSimulationLoadout(options map[string][]string, monster schemas.MonsterSchema, available map[string]int, forcedWeapon, profile string) map[string]string {
	loadout := make(map[string]string, len(simulationBestSlots))
	used := make(map[string]int)
	for _, slot := range simulationBestSlots {
		if slot == "weapon" && forcedWeapon != "" {
			loadout[slot] = forcedWeapon
			used[forcedWeapon]++
			continue
		}
		bestCode := ""
		bestScore := math.Inf(-1)
		selectedWeaponElement := weaponElement(loadout["weapon"])
		for _, code := range options[slot] {
			if code != "" && used[code] >= available[code] {
				continue
			}
			if strings.HasPrefix(slot, "artifact") && code != "" {
				duplicate := false
				for _, artifactSlot := range []string{
					"artifact1",
					"artifact2",
					"artifact3",
				} {
					if loadout[artifactSlot] == code {
						duplicate = true
					}
				}
				if duplicate {
					continue
				}
			}
			item, ok := database.Items().Get(code)
			if !ok {
				continue
			}
			score := simulationItemScore(*item, monster, selectedWeaponElement, profile)
			if score > bestScore {
				bestCode, bestScore = code, score
			}
		}
		loadout[slot] = bestCode
		if bestCode != "" {
			used[bestCode]++
		}
	}
	return loadout
}

func simulationItemScore(item schemas.ItemSchema, monster schemas.MonsterSchema, selectedWeaponElement, profile string) float64 {
	if item.Effects == nil {
		return 0
	}
	score := 0.0
	for _, effect := range *item.Effects {
		value := float64(effect.Value)
		switch effect.Code {
		case "attack_fire":
			score += value * (1 - float64(monster.ResFire)/100)
		case "attack_earth":
			score += value * (1 - float64(monster.ResEarth)/100)
		case "attack_water":
			score += value * (1 - float64(monster.ResWater)/100)
		case "attack_air":
			score += value * (1 - float64(monster.ResAir)/100)
		case "dmg":
			score += value * 1.5
		case "dmg_fire", "dmg_earth", "dmg_water", "dmg_air":
			if effect.Code == "dmg_"+selectedWeaponElement {
				score += value * 1.5
			}
		case "critical_strike":
			score += value * 2
		case "hp":
			weight := 0.25
			switch profile {
			case "survival":
				weight = 1
			case "balanced":
				weight = 0.5
			}
			score += value * weight
		case "res_fire", "res_earth", "res_water", "res_air":
			weight := 0.5
			switch profile {
			case "survival":
				weight = 1.5
			case "balanced":
				weight = 0.75
			}
			score += value * weight
		case "lifesteal", "healing", "shell", "burn":
			score += value * 0.75
		case "wisdom", "prospecting":
			if effect.Code == profile {
				score += value * 0.35
			}
		}
	}
	return score
}

func weaponElement(code string) string {
	item, ok := database.Items().Get(code)
	if !ok || item.Effects == nil {
		return ""
	}
	for _, effect := range *item.Effects {
		switch effect.Code {
		case "attack_fire":
			return "fire"
		case "attack_earth":
			return "earth"
		case "attack_water":
			return "water"
		case "attack_air":
			return "air"
		}
	}
	return ""
}

func evaluateSimulationLoadout(character schemas.CharacterSchema, monster schemas.MonsterSchema, slots map[string]string) Result {
	return evaluateSimulationLoadoutIterations(character, monster, slots, candidateSimulationIterations)
}

func evaluateSimulationLoadoutIterations(character schemas.CharacterSchema, monster schemas.MonsterSchema, slots map[string]string, iterations int) Result {
	cacheKey := bestFightSimulationCacheKey(character, monster, slots, iterations)
	cached, memoryHit := processSimulationCache[cacheKey]
	if memoryHit {
		return simulationResultFromCache(cached)
	}
	stored, ok := cache.GetFightSimulation(cacheKey, bestFightSimulationCacheVersion)
	if ok {
		cachedResult := cachedSimulationResultFromModel(stored)
		processSimulationCache[cacheKey] = cachedResult
		return simulationResultFromCache(cachedResult)
	}
	gear := copyStringMap(slots)
	utilities := make(map[string]int)
	for _, slot := range simulationBestUtilitySlots {
		code := gear[slot]
		if code != "" {
			utilities[code] = 100
		}
		delete(gear, slot)
	}
	// Rebuild the fighter from level and the complete candidate loadout, just
	// like simulation local does. Character attributes are entirely derived
	// from these values, so there is no need to subtract the character's
	// currently equipped items from API-calculated stats.
	fighter := fight.FromLoadout(character.Level, gear, utilities)
	// Use common random numbers: every candidate sees the same critical-hit
	// sequence, so tiny score differences reflect the loadout rather than luck.
	summary := fight.SimulateMany(fighter, monster, fight.SimulationOptions{
		Iterations: iterations,
		RNG:        rand.New(rand.NewSource(1)).Float64,
	})
	playerCritical, monsterCritical := 0, 100
	criticalOptions := fight.SimulationOptions{
		Critical: fight.CriticalOptions{
			PlayerChance:  &playerCritical,
			MonsterChance: &monsterCritical,
		},
	}
	criticalOptions.Iterations = iterations
	criticalOptions.RNG = rand.New(rand.NewSource(1)).Float64
	pessimistic := fight.SimulateMany(fighter, monster, criticalOptions)
	var turns, hp float32
	for _, result := range summary.Results {
		turns += float32(result.Turns)
		if len(result.CharacterResults) > 0 {
			switch value := result.CharacterResults[0]["final_hp"].(type) {
			case int:
				hp += float32(value)
			case float64:
				hp += float32(value)
			}
		}
	}
	if len(summary.Results) > 0 {
		turns /= float32(len(summary.Results))
		hp /= float32(len(summary.Results))
	}
	metrics := fight.Metrics(fighter, character.Level, monster, summary.Results)
	damagePerTurn := effectivePlayerDamage(fighter, monster)
	damageSurplus := damagePerTurn*turns - float32(monster.Hp)
	survivalSurplus := hp / float32(fighter.Stats.HP)
	missingHP := float32(fighter.Stats.HP) - hp
	restCooldown := float32(0)
	if missingHP > 0 {
		restCooldown = float32(math.Ceil(float64(missingHP) / float64(fighter.Stats.HP) * 100))
		if restCooldown < 3 {
			restCooldown = 3
		}
	}
	cycleCost := metrics.AverageFightCooldown + restCooldown
	result := Result{
		Winrate:               summary.Winrate,
		AverageTurns:          turns,
		AverageFinalHP:        hp,
		AverageFightCooldown:  metrics.AverageFightCooldown,
		EstimatedRestCooldown: restCooldown,
		CycleCost:             cycleCost,
		DamageSurplus:         damageSurplus,
		SurvivalSurplus:       survivalSurplus,
		XP:                    metrics.XP,
		XPPerSecond:           float32(metrics.XP) / cycleCost,
		GoldPerSecond:         float32(monster.MinGold+monster.MaxGold) / 2 / cycleCost,
		Safe:                  pessimistic.Winrate == 100,
		ArtifactOrderScore:    simulationArtifactOrderScore(character, slots),
	}
	result.Wisdom, result.Prospecting = simulationLoadoutMeta(character, slots)
	result.ProspectingEfficiency = (1 + float32(result.Prospecting)/100) / cycleCost
	result.ContextScore = simulationContextScore(character, monster, result.Wisdom, result.Prospecting)
	cachedResult := cachedSimulationResult{
		Winrate:               result.Winrate,
		AverageTurns:          result.AverageTurns,
		AverageFinalHP:        result.AverageFinalHP,
		AverageFightCooldown:  result.AverageFightCooldown,
		EstimatedRestCooldown: result.EstimatedRestCooldown,
		CycleCost:             result.CycleCost,
		DamageSurplus:         result.DamageSurplus,
		SurvivalSurplus:       result.SurvivalSurplus,
		XP:                    result.XP,
		XPPerSecond:           result.XPPerSecond,
		GoldPerSecond:         result.GoldPerSecond,
		ProspectingEfficiency: result.ProspectingEfficiency,
		Safe:                  result.Safe,
	}
	processSimulationCache[cacheKey] = cachedResult
	cache.SaveFightSimulation(fightSimulationFromResult(cacheKey, cachedResult))
	return result
}

func cachedSimulationResultFromModel(stored models.FightSimulation) cachedSimulationResult {
	return cachedSimulationResult{
		Winrate:               stored.Winrate,
		AverageTurns:          stored.AverageTurns,
		AverageFinalHP:        stored.AverageFinalHP,
		AverageFightCooldown:  stored.AverageFightCooldown,
		EstimatedRestCooldown: stored.EstimatedRestCooldown,
		CycleCost:             stored.CycleCost,
		DamageSurplus:         stored.DamageSurplus,
		SurvivalSurplus:       stored.SurvivalSurplus,
		XP:                    stored.XP,
		XPPerSecond:           stored.XPPerSecond,
		GoldPerSecond:         stored.GoldPerSecond,
		ProspectingEfficiency: stored.ProspectingEfficiency,
		Safe:                  stored.Safe,
	}
}

func fightSimulationFromResult(key string, result cachedSimulationResult) models.FightSimulation {
	return models.FightSimulation{
		Key:                   key,
		Version:               bestFightSimulationCacheVersion,
		Winrate:               result.Winrate,
		AverageTurns:          result.AverageTurns,
		AverageFinalHP:        result.AverageFinalHP,
		AverageFightCooldown:  result.AverageFightCooldown,
		EstimatedRestCooldown: result.EstimatedRestCooldown,
		CycleCost:             result.CycleCost,
		DamageSurplus:         result.DamageSurplus,
		SurvivalSurplus:       result.SurvivalSurplus,
		XP:                    result.XP,
		XPPerSecond:           result.XPPerSecond,
		GoldPerSecond:         result.GoldPerSecond,
		ProspectingEfficiency: result.ProspectingEfficiency,
		Safe:                  result.Safe,
	}
}

func simulationResultFromCache(cached cachedSimulationResult) Result {
	return Result{
		Winrate:               cached.Winrate,
		AverageTurns:          cached.AverageTurns,
		AverageFinalHP:        cached.AverageFinalHP,
		AverageFightCooldown:  cached.AverageFightCooldown,
		EstimatedRestCooldown: cached.EstimatedRestCooldown,
		CycleCost:             cached.CycleCost,
		DamageSurplus:         cached.DamageSurplus,
		SurvivalSurplus:       cached.SurvivalSurplus,
		XP:                    cached.XP,
		XPPerSecond:           cached.XPPerSecond,
		GoldPerSecond:         cached.GoldPerSecond,
		ProspectingEfficiency: cached.ProspectingEfficiency,
		Safe:                  cached.Safe,
	}
}

func bestFightSimulationCacheKey(character schemas.CharacterSchema, monster schemas.MonsterSchema, slots map[string]string, iterations int) string {
	value := fmt.Sprintf("%d|%d|%s|%d|%s", bestFightSimulationCacheVersion, character.Level, monster.Code, iterations, simulationLoadoutKey(slots))
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func effectivePlayerDamage(player fight.Fighter, monster schemas.MonsterSchema) float32 {
	resistances := []int{
		monster.ResFire,
		monster.ResEarth,
		monster.ResWater,
		monster.ResAir,
	}
	attacks := []int{
		player.Stats.AttackFire,
		player.Stats.AttackEarth,
		player.Stats.AttackWater,
		player.Stats.AttackAir,
	}
	bonuses := []int{
		player.Stats.DmgFire,
		player.Stats.DmgEarth,
		player.Stats.DmgWater,
		player.Stats.DmgAir,
	}
	var total float32
	for i := range attacks {
		multiplier := 1 - float32(resistances[i])/100
		if multiplier < 0 {
			multiplier = 0
		}
		total += float32(attacks[i]) * (1 + float32(player.Stats.Dmg+bonuses[i])/100) * multiplier
	}
	return total
}

func chooseSimulationGroup(character schemas.CharacterSchema, monster schemas.MonsterSchema, current map[string]string, slots, options []string, available map[string]int, allowDuplicateAdeptRing, unique bool) {
	if len(options) > 13 {
		ranked := append([]string{}, options...)
		sort.SliceStable(ranked, func(i, j int) bool {
			if ranked[i] == "" {
				return true
			}
			if ranked[j] == "" {
				return false
			}
			ii, _ := database.Items().Get(ranked[i])
			jj, _ := database.Items().Get(ranked[j])
			return simulationItemScore(*ii, monster, weaponElement(current["weapon"]), "damage") > simulationItemScore(*jj, monster, weaponElement(current["weapon"]), "damage")
		})
		options = append([]string{""}, ranked[:12]...)
	}
	best := copyStringMap(current)
	bestScore := evaluateSimulationLoadout(character, monster, best)
	setSimulationContext(&bestScore, character, monster, best)
	bestKeep := simulationGroupKeepCount(current, best, slots)
	var visit func(int, []string)
	visit = func(pos int, chosen []string) {
		if pos == len(slots) {
			candidate := copyStringMap(current)
			for i, slot := range slots {
				candidate[slot] = chosen[i]
			}
			if !loadoutHasQuantity(candidate, available) || (!allowDuplicateAdeptRing && !adeptRingAllowed(character, candidate)) {
				return
			}
			if unique {
				seen := map[string]bool{}
				for _, slot := range slots {
					code := candidate[slot]
					if code != "" {
						if seen[code] {
							return
						}
						seen[code] = true
					}
				}
			}
			score := evaluateSimulationLoadout(character, monster, candidate)
			setSimulationContext(&score, character, monster, candidate)
			keep := simulationGroupKeepCount(current, candidate, slots)
			if betterSearchScore(score, bestScore) || (!betterSearchScore(bestScore, score) && keep > bestKeep) {
				best, bestScore, bestKeep = candidate, score, keep
			}
			return
		}
		for _, code := range options {
			visit(pos+1, append(chosen, code))
		}
	}
	visit(0, nil)
	for _, slot := range slots {
		current[slot] = best[slot]
	}
}

func simulationGroupKeepCount(original, candidate map[string]string, slots []string) int {
	kept := 0
	for _, slot := range slots {
		if original[slot] != "" && original[slot] == candidate[slot] {
			kept++
		}
	}
	return kept
}

func betterSearchScore(a, b Result) bool {
	if a.Winrate != b.Winrate {
		return a.Winrate > b.Winrate
	}
	if a.Winrate > 0 && b.Winrate > 0 && a.CycleCost != b.CycleCost {
		return a.CycleCost < b.CycleCost
	}
	if a.DamageSurplus != b.DamageSurplus {
		return a.DamageSurplus > b.DamageSurplus
	}
	if a.SurvivalSurplus != b.SurvivalSurplus {
		return a.SurvivalSurplus > b.SurvivalSurplus
	}
	if a.AverageFinalHP != b.AverageFinalHP {
		return a.AverageFinalHP > b.AverageFinalHP
	}
	if a.AverageTurns != b.AverageTurns {
		return a.AverageTurns < b.AverageTurns
	}
	return a.ArtifactOrderScore > b.ArtifactOrderScore
}

func betterSimulationScore(a, b Result) bool {
	if a.Winrate != b.Winrate {
		return a.Winrate > b.Winrate
	}
	if a.Winrate > 0 && b.Winrate > 0 && a.CycleCost != b.CycleCost {
		return a.CycleCost < b.CycleCost
	}
	if a.Winrate > 0 && b.Winrate > 0 && a.FocusScore != b.FocusScore {
		return a.FocusScore > b.FocusScore
	}
	if a.ArtifactOrderScore != b.ArtifactOrderScore {
		return a.ArtifactOrderScore > b.ArtifactOrderScore
	}
	if a.AverageFinalHP != b.AverageFinalHP {
		return a.AverageFinalHP > b.AverageFinalHP
	}
	if a.AverageTurns != b.AverageTurns {
		return a.AverageTurns < b.AverageTurns
	}
	return a.ContextScore > b.ContextScore
}

func adeptRingAllowed(character schemas.CharacterSchema, slots map[string]string) bool {
	original := characterSlots(character)
	originalCount := 0
	for _, slot := range []string{"ring1", "ring2"} {
		if original[slot] == "ring_of_the_adept" {
			originalCount++
		}
	}
	limit := originalCount
	if limit < 1 {
		limit = 1
	}
	count := 0
	for _, slot := range []string{"ring1", "ring2"} {
		if slots[slot] == "ring_of_the_adept" {
			count++
		}
	}
	return count <= limit
}

func simulationArtifactOrderScore(character schemas.CharacterSchema, slots map[string]string) int {
	original := characterSlots(character)
	score := 0
	for _, slot := range []string{"artifact1", "artifact2", "artifact3"} {
		if original[slot] != "" && original[slot] == slots[slot] {
			score++
		}
	}
	return score
}

func simulationContextScore(character schemas.CharacterSchema, monster schemas.MonsterSchema, wisdom, prospecting int) int {
	if simulationContextUsesWisdom(character, monster) {
		return wisdom
	}
	return prospecting
}

func setSimulationContext(result *Result, character schemas.CharacterSchema, monster schemas.MonsterSchema, slots map[string]string) {
	result.Wisdom, result.Prospecting = simulationLoadoutMeta(character, slots)
	result.ContextScore = simulationContextScore(character, monster, result.Wisdom, result.Prospecting)
	if simulationContextUsesWisdom(character, monster) {
		result.Focus = "wisdom"
		result.FocusScore = result.XPPerSecond
		return
	}
	result.Focus = "prospecting"
	result.FocusScore = result.ProspectingEfficiency
}

func simulationContextUsesWisdom(character schemas.CharacterSchema, monster schemas.MonsterSchema) bool {
	difference := character.Level - monster.Level
	if difference < 0 {
		difference = -difference
	}
	return difference <= 10
}

func simulationLoadoutMeta(character schemas.CharacterSchema, slots map[string]string) (int, int) {
	wisdom, prospecting := character.Wisdom, character.Prospecting
	for _, code := range characterSlots(character) {
		wisdom -= itemMeta(code, "wisdom")
		prospecting -= itemMeta(code, "prospecting")
	}
	for _, code := range slots {
		wisdom += itemMeta(code, "wisdom")
		prospecting += itemMeta(code, "prospecting")
	}
	return wisdom, prospecting
}

func itemMeta(code, wanted string) int {
	item, ok := database.Items().Get(code)
	if !ok || item.Effects == nil {
		return 0
	}
	value := 0
	for _, effect := range *item.Effects {
		if effect.Code == wanted {
			value += effect.Value
		}
	}
	return value
}

func characterSlots(c schemas.CharacterSchema) map[string]string {
	return map[string]string{
		"weapon":     c.WeaponSlot,
		"shield":     c.ShieldSlot,
		"helmet":     c.HelmetSlot,
		"body_armor": c.BodyArmorSlot,
		"leg_armor":  c.LegArmorSlot,
		"boots":      c.BootsSlot,
		"amulet":     c.AmuletSlot,
		"rune":       c.RuneSlot,
		"ring1":      c.Ring1Slot,
		"ring2":      c.Ring2Slot,
		"artifact1":  c.Artifact1Slot,
		"artifact2":  c.Artifact2Slot,
		"artifact3":  c.Artifact3Slot,
	}
}

func itemFitsSlot(item schemas.ItemSchema, slot string) bool {
	for _, candidate := range database.EquipmentTypeToSlots[item.Type] {
		if candidate == slot {
			return true
		}
	}
	return false
}

func CanEquip(c schemas.CharacterSchema, item schemas.ItemSchema) bool {
	if item.Conditions == nil {
		return true
	}
	for _, condition := range *item.Conditions {
		var value int
		switch condition.Code {
		case "level":
			value = c.Level
		default:
			continue
		}
		switch condition.Operator {
		case schemas.Gt:
			if value <= condition.Value {
				return false
			}
		case schemas.Lt:
			if value >= condition.Value {
				return false
			}
		case schemas.Eq:
			if value != condition.Value {
				return false
			}
		case schemas.Ne:
			if value == condition.Value {
				return false
			}
		}
	}
	return true
}

func copyStringMap(source map[string]string) map[string]string {
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func loadoutHasQuantity(slots map[string]string, available map[string]int) bool {
	used := make(map[string]int)
	for _, code := range slots {
		if code == "" {
			continue
		}
		used[code]++
		if used[code] > available[code] {
			return false
		}
	}
	artifacts := map[string]bool{}
	for _, slot := range []string{"artifact1", "artifact2", "artifact3"} {
		code := slots[slot]
		if code == "" {
			continue
		}
		if artifacts[code] {
			return false
		}
		artifacts[code] = true
	}
	utilities := map[string]bool{}
	for _, slot := range simulationBestUtilitySlots {
		code := slots[slot]
		if code == "" {
			continue
		}
		if utilities[code] {
			return false
		}
		utilities[code] = true
	}
	return true
}
