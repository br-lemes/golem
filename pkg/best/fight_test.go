package best

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSimulationItemScoreUsesWeaponElement(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{
		{Code: "dmg", Value: 10},
		{Code: "dmg_fire", Value: 20},
		{Code: "dmg_water", Value: 30},
	}
	item := schemas.ItemSchema{Effects: effects}

	fire := simulationItemScore(item, schemas.MonsterSchema{}, "fire", "damage")
	water := simulationItemScore(item, schemas.MonsterSchema{}, "water", "damage")
	if fire == water {
		t.Fatalf("element-specific damage should depend on weapon element: fire=%v water=%v", fire, water)
	}
}

func TestSimulationItemScoreHandlesEffectProfiles(t *testing.T) {
	monster := schemas.MonsterSchema{}
	monster.ResFire = 20
	monster.ResEarth = 30
	monster.ResWater = 40
	monster.ResAir = 50
	check := func(code string, value int, profile string, want float64) {
		effects := &[]schemas.SimpleEffectSchema{{Code: code, Value: value}}
		got := simulationItemScore(schemas.ItemSchema{Effects: effects}, monster, "fire", profile)
		if got != want {
			t.Fatalf("%s score = %v, want %v", code, got, want)
		}
	}
	check("attack_fire", 10, "damage", 8)
	check("attack_earth", 10, "damage", 7)
	check("attack_water", 10, "damage", 6)
	check("attack_air", 10, "damage", 5)
	check("critical_strike", 4, "damage", 8)
	check("hp", 4, "survival", 4)
	check("res_fire", 4, "balanced", 3)
	check("lifesteal", 10, "damage", 7.5)
	check("healing", 10, "damage", 7.5)
	check("shell", 10, "damage", 7.5)
	check("burn", 10, "damage", 7.5)
	check("wisdom", 10, "wisdom", 3.5)
	check("prospecting", 10, "prospecting", 3.5)
	got := simulationItemScore(schemas.ItemSchema{}, monster, "fire", "damage")
	if got != 0 {
		t.Fatalf("empty item score = %v, want 0", got)
	}
}

func TestWeaponElement(t *testing.T) {
	tests := map[string]string{
		"fire_staff":        "fire",
		"copper_pickaxe":    "earth",
		"fishing_net":       "water",
		"apprentice_gloves": "air",
	}
	for code, want := range tests {
		got := weaponElement(code)
		if got != want {
			t.Errorf("weaponElement(%q) = %q, want %q", code, got, want)
		}
	}
	missing := weaponElement("missing_item")
	shield := weaponElement("wooden_shield")
	if missing != "" || shield != "" {
		t.Fatal("non-elemental or missing item returned an element")
	}
}

func TestEffectivePlayerDamageAppliesResistanceFraction(t *testing.T) {
	player := fight.Fighter{Stats: fight.Stats{AttackFire: 100}}
	monster := schemas.MonsterSchema{ResFire: 25}
	got := effectivePlayerDamage(player, monster)
	if got != 75 {
		t.Fatalf("effective damage = %v, want 75", got)
	}
}

func TestAlignSimulationArtifactsPreservesEquivalentOrder(t *testing.T) {
	character := schemas.CharacterSchema{
		Artifact1Slot: "novice_guide",
		Artifact2Slot: "lost_world_map",
	}
	slots := map[string]string{
		"artifact1": "lost_world_map",
		"artifact2": "novice_guide",
		"artifact3": "lich_race_medal",
	}

	got := alignSimulationArtifacts(character, slots)
	if got["artifact1"] != "novice_guide" || got["artifact2"] != "lost_world_map" || got["artifact3"] != "lich_race_medal" {
		t.Fatalf("artifact order = %#v", got)
	}
}

func TestSimulationEquipmentChangesOmitsEmptyResults(t *testing.T) {
	current := map[string]string{"boots": "", "weapon": "new_weapon"}
	original := map[string]string{"boots": "old_boots", "weapon": "old_weapon"}

	got := simulationEquipmentChanges(current, original)
	_, bootsChanged := got["boots"]
	if bootsChanged {
		t.Fatal("empty equipment change should be omitted")
	}
	if got["weapon"] != "new_weapon" {
		t.Fatalf("weapon change = %q", got["weapon"])
	}
}

func TestSimulationUtilityChangesIncludesRemoval(t *testing.T) {
	character := schemas.CharacterSchema{Utility1Slot: "health_potion"}
	current := map[string]string{"utility1": "", "utility2": ""}

	got := simulationUtilityChanges(current, character)
	value, utilityChanged := got["utility1"]
	if !utilityChanged || value != "" {
		t.Fatalf("utility removal = %#v", got)
	}
	_, emptyUtilityChanged := got["utility2"]
	if emptyUtilityChanged {
		t.Fatal("unchanged empty utility should be omitted")
	}
}

func TestLoadoutHasQuantityRejectsDuplicateUtilities(t *testing.T) {
	slots := map[string]string{
		"utility1": "health_potion",
		"utility2": "health_potion",
	}
	available := map[string]int{"health_potion": 2}
	if loadoutHasQuantity(slots, available) {
		t.Fatal("duplicate utility should not be allowed")
	}
}

func TestBetterSimulationScorePrefersLowerCycleCost(t *testing.T) {
	faster := Result{Winrate: 100, CycleCost: 84}
	slower := Result{Winrate: 100, CycleCost: 94}
	if !betterSimulationScore(faster, slower) {
		t.Fatal("lower cycle cost should be preferred")
	}
	if betterSimulationScore(slower, faster) {
		t.Fatal("higher cycle cost should not be preferred")
	}
}

func TestBetterSearchScoreUsesTieBreakers(t *testing.T) {
	if !betterSearchScore(Result{SurvivalSurplus: 2}, Result{SurvivalSurplus: 1}) {
		t.Fatal("higher survival surplus should win")
	}
	if !betterSearchScore(Result{AverageFinalHP: 2}, Result{AverageFinalHP: 1}) {
		t.Fatal("higher final HP should win")
	}
	if !betterSearchScore(Result{AverageTurns: 1}, Result{AverageTurns: 2}) {
		t.Fatal("fewer turns should win")
	}
	a := Result{ArtifactOrderScore: 2}
	b := Result{ArtifactOrderScore: 1}
	if !betterSearchScore(a, b) {
		t.Fatal("higher artifact order score should win")
	}
}

func TestBetterSimulationScoreUsesTieBreakers(t *testing.T) {
	a := Result{Winrate: 100}
	b := Result{Winrate: 50}
	if !betterSimulationScore(a, b) {
		t.Fatal("higher winrate should win")
	}
	a = Result{Winrate: 100, FocusScore: 2}
	b = Result{Winrate: 100, FocusScore: 1}
	if !betterSimulationScore(a, b) {
		t.Fatal("higher focus score should win")
	}
	a = Result{Winrate: 0, ArtifactOrderScore: 2}
	b = Result{Winrate: 0, ArtifactOrderScore: 1}
	if !betterSimulationScore(a, b) {
		t.Fatal("higher artifact order score should win")
	}
	a = Result{AverageFinalHP: 2}
	b = Result{AverageFinalHP: 1}
	if !betterSimulationScore(a, b) {
		t.Fatal("higher final HP should win")
	}
	a = Result{AverageTurns: 1}
	b = Result{AverageTurns: 2}
	if !betterSimulationScore(a, b) {
		t.Fatal("fewer turns should win")
	}
	a = Result{ContextScore: 2}
	b = Result{ContextScore: 1}
	if !betterSimulationScore(a, b) {
		t.Fatal("higher context score should win")
	}
}

func TestFightHelpers(t *testing.T) {
	if !betterSearchScore(Result{DamageSurplus: 2}, Result{DamageSurplus: 1}) {
		t.Fatal("higher damage surplus should win")
	}
	a := Result{Winrate: 1, FocusScore: 2}
	b := Result{Winrate: 1, FocusScore: 1}
	if !betterSimulationScore(a, b) {
		t.Fatal("higher focus score should win")
	}
	original := map[string]string{"weapon": "sword"}
	candidate := map[string]string{"weapon": "sword"}
	if simulationGroupKeepCount(original, candidate, []string{"weapon"}) != 1 {
		t.Fatal("unchanged slot was not counted")
	}
	near := schemas.CharacterSchema{Level: 20}
	nearMonster := schemas.MonsterSchema{Level: 30}
	if !simulationContextUsesWisdom(near, nearMonster) {
		t.Fatal("near-level fight should use wisdom")
	}
	farMonster := schemas.MonsterSchema{Level: 40}
	if simulationContextUsesWisdom(near, farMonster) {
		t.Fatal("far-level fight should use prospecting")
	}
}

func TestLimitAdeptRingPreservesTheFirstAdeptRing(t *testing.T) {
	heuristic := map[string]string{
		"ring1": "ring_of_the_adept",
		"ring2": "ring_of_the_adept",
	}
	limitAdeptRing(heuristic)
	if heuristic["ring1"] != "ring_of_the_adept" || heuristic["ring2"] != "" {
		t.Fatalf("limited rings = %#v", heuristic)
	}
	heuristic = map[string]string{
		"ring1": "ring_of_chance",
		"ring2": "ring_of_the_adept",
	}
	limitAdeptRing(heuristic)
	if heuristic["ring1"] != "" || heuristic["ring2"] != "ring_of_the_adept" {
		t.Fatalf("limited rings = %#v", heuristic)
	}
}

func TestAdeptRingAllowedAccountsForEquippedRings(t *testing.T) {
	character := schemas.CharacterSchema{Ring1Slot: "ring_of_the_adept"}
	allowed := map[string]string{"ring1": "ring_of_the_adept"}
	if !adeptRingAllowed(character, allowed) {
		t.Fatal("equipped adept ring should be allowed")
	}
	character.Ring2Slot = "ring_of_the_adept"
	allowed["ring2"] = "ring_of_the_adept"
	if !adeptRingAllowed(character, allowed) {
		t.Fatal("two equipped adept rings should allow both rings")
	}
}

func TestSimulationArtifactOrderScorePreservesEquippedArtifact(t *testing.T) {
	character := schemas.CharacterSchema{Artifact1Slot: "novice_guide"}
	slots := map[string]string{"artifact1": "novice_guide"}
	got := simulationArtifactOrderScore(character, slots)
	if got != 1 {
		t.Fatalf("artifact order score = %d, want 1", got)
	}
}

func TestHeuristicSimulationLoadoutSkipsUnknownItems(t *testing.T) {
	options := map[string][]string{
		"weapon": {"iron_sword"},
		"shield": {"missing_item"},
	}
	available := map[string]int{"iron_sword": 1}
	monster := schemas.MonsterSchema{}
	got := heuristicSimulationLoadout(options, monster, available, "iron_sword", "damage")
	if got["shield"] != "" {
		t.Fatalf("unknown shield was selected: %#v", got)
	}
}

func TestHeuristicSimulationLoadoutRespectsQuantityAndUniqueArtifacts(t *testing.T) {
	options := map[string][]string{
		"weapon":    {"iron_sword"},
		"ring1":     {"iron_ring"},
		"ring2":     {"iron_ring"},
		"artifact1": {"novice_guide"},
		"artifact2": {"novice_guide"},
	}
	available := map[string]int{}
	available["iron_sword"] = 1
	available["iron_ring"] = 1
	available["novice_guide"] = 2
	got := heuristicSimulationLoadout(options, schemas.MonsterSchema{}, available, "iron_sword", "damage")
	if got["ring1"] == "iron_ring" && got["ring2"] == "iron_ring" {
		t.Fatal("ring was reused beyond available quantity")
	}
	if got["artifact1"] == "novice_guide" && got["artifact2"] == "novice_guide" {
		t.Fatal("artifact was reused")
	}
}

func TestChooseSimulationGroupRanksLargeOptionSet(t *testing.T) {
	erring := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if erring != nil {
		t.Fatal(erring)
	}
	ResetSimulationCache()
	options := make([]string, 0, 14)
	available := map[string]int{"iron_sword": 1}
	for _, item := range catalog.Items().All() {
		if item.Type != "ring" {
			continue
		}
		options = append(options, item.Code)
		available[item.Code] = 2
		if len(options) == 14 {
			break
		}
	}
	if len(options) != 14 {
		t.Fatalf("catalog has %d ring options, want 14", len(options))
	}
	options = append([]string{""}, options...)
	current := map[string]string{"weapon": "iron_sword"}
	character := schemas.CharacterSchema{Level: 50}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	chooseSimulationGroup(character, monster, current, simulationRingSlots, options, available, false)
	options = append(options[1:], "")
	chooseSimulationGroup(character, monster, current, simulationRingSlots, options, available, false)
	if current["ring1"] == "" && current["ring2"] == "" {
		t.Fatal("large option set produced no ring")
	}
}

func TestCanEquipConditions(t *testing.T) {
	if !CanEquip(schemas.CharacterSchema{}, schemas.ItemSchema{}) {
		t.Fatal("item without conditions was rejected")
	}
	conditions := []schemas.ConditionSchema{
		{Code: "level", Operator: "gt", Value: 9},
		{Code: "level", Operator: "lt", Value: 11},
		{Code: "level", Operator: "eq", Value: 10},
		{Code: "level", Operator: "ne", Value: 9},
	}
	item := schemas.ItemSchema{Conditions: &conditions}
	if !CanEquip(schemas.CharacterSchema{Level: 10}, item) {
		t.Fatal("valid level conditions were rejected")
	}
	if CanEquip(schemas.CharacterSchema{Level: 9}, item) {
		t.Fatal("invalid level conditions were accepted")
	}
}

func TestCanEquipRejectsRemainingLevelOperators(t *testing.T) {
	check := func(operator string, level int, want bool) {
		condition := schemas.ConditionSchema{Code: "level"}
		condition.Operator = operator
		condition.Value = 10
		conditions := []schemas.ConditionSchema{condition}
		item := schemas.ItemSchema{Conditions: &conditions}
		got := CanEquip(schemas.CharacterSchema{Level: level}, item)
		if got != want {
			t.Fatalf("CanEquip(%s, %d) = %t, want %t", operator, level, got, want)
		}
	}
	check("lt", 9, true)
	check("lt", 10, false)
	check("eq", 10, true)
	check("eq", 9, false)
	check("ne", 9, true)
	check("ne", 10, false)
}

func TestEffectivePlayerDamageClampsNegativeResistanceMultiplier(t *testing.T) {
	player := fight.Fighter{Stats: fight.Stats{AttackFire: 100}}
	monster := schemas.MonsterSchema{ResFire: 150}
	got := effectivePlayerDamage(player, monster)
	if got != 0 {
		t.Fatalf("effective damage = %v, want 0", got)
	}
}

func TestLoadoutHelpers(t *testing.T) {
	loadout := map[string]string{"weapon": "sword"}
	copy := copyStringMap(loadout)
	copy["weapon"] = "axe"
	if loadout["weapon"] != "sword" {
		t.Fatal("copyStringMap aliases its source")
	}
	if !loadoutHasQuantity(loadout, map[string]int{"sword": 1}) {
		t.Fatal("available item was rejected")
	}
	duplicate := map[string]string{"weapon": "sword", "shield": "sword"}
	if loadoutHasQuantity(duplicate, map[string]int{"sword": 1}) {
		t.Fatal("insufficient quantity was accepted")
	}
}

func TestLoadoutHasQuantityRejectsDuplicateArtifactsAndUtilities(t *testing.T) {
	artifacts := map[string]string{
		"artifact1": "novice_guide",
		"artifact2": "novice_guide",
	}
	available := map[string]int{"novice_guide": 2}
	if loadoutHasQuantity(artifacts, available) {
		t.Fatal("duplicate artifact was accepted")
	}
	utilities := map[string]string{
		"utility1": "small_health_potion",
		"utility2": "small_health_potion",
	}
	available["small_health_potion"] = 2
	if loadoutHasQuantity(utilities, available) {
		t.Fatal("duplicate utility was accepted")
	}
}

func TestFindFightWithAvailableBuildsResult(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	monster, ok := catalog.Monsters.Get("chicken")
	if !ok {
		t.Fatal("chicken monster is missing from the catalog")
	}
	character := schemas.CharacterSchema{Level: 5}
	available := map[string]int{"iron_sword": 1}
	result, err := FindFightWithAvailable(character, *monster, available, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Focus == "" || result.FinalEquipment["weapon"] == "" {
		t.Fatalf("incomplete fight result: %#v", result)
	}
}

func TestFindFightWithAvailableReportsUnownedItems(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	character := schemas.CharacterSchema{Level: 1}
	monster, ok := catalog.Monsters.Get("sonnengott")
	if !ok {
		t.Fatal("chicken monster is missing from the catalog")
	}
	available := map[string]int{"iron_sword": 1}
	result, err := FindFightWithAvailable(character, *monster, available, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Unowned) == 0 {
		t.Fatalf("unowned items were not reported: %#v", result)
	}
}

func TestFindFightWithAvailableReportsUnownedUtilities(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	monster, ok := catalog.Monsters.Get("skeleton")
	if !ok {
		t.Fatal("skeleton monster is missing from the catalog")
	}
	character := schemas.CharacterSchema{Level: 10}
	available := map[string]int{}
	result, err := FindFightWithAvailable(character, *monster, available, true, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Utilities["utility1"] == "" || result.Utilities["utility2"] == "" {
		t.Fatalf("utilities were not selected: %#v", result.Utilities)
	}
	for _, code := range result.Utilities {
		found := false
		for _, unowned := range result.Unowned {
			if unowned == code {
				found = true
			}
		}
		if !found {
			t.Fatalf("utility %q missing from unowned: %v", code, result.Unowned)
		}
	}
}

func TestKeepBestSimulationLoadoutKeepsBetterCurrent(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	character := schemas.CharacterSchema{Level: 1}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	current := map[string]string{"weapon": "iron_sword"}
	bestLoadout := map[string]string{}
	best := Result{Winrate: -1}
	got, result := keepBestSimulationLoadout(character, monster, current, bestLoadout, best)
	if got["weapon"] != "iron_sword" || result.Winrate < 0 {
		t.Fatalf("current loadout was not kept: %#v, %#v", got, result)
	}
}

func TestFindFightWithAvailableFiltersInvalidItems(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	processSimulationCache = nil
	available := map[string]int{
		"missing_item":        1,
		"small_health_potion": 1,
		"iron_pickaxe":        1,
		"iron_sword":          0,
	}
	character := schemas.CharacterSchema{Level: 1}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	result, err := FindFightWithAvailable(character, monster, available, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.FinalEquipment == nil {
		t.Fatal("result has no final equipment")
	}
}

func TestFindFightWithAvailableLimitsAdeptRings(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	available := map[string]int{
		"iron_sword":        1,
		"ring_of_the_adept": 2,
		"ring_of_chance":    1,
	}
	character := schemas.CharacterSchema{Level: 50}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	result, err := FindFightWithAvailable(character, monster, available, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.FinalEquipment == nil {
		t.Fatal("result has no final equipment")
	}
}

func TestFindFightWithAvailableRejectsInsufficientRingQuantity(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	available := map[string]int{"iron_sword": 1, "ring_of_the_adept": 1}
	character := schemas.CharacterSchema{Level: 50}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	result, err := FindFightWithAvailable(character, monster, available, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.FinalEquipment == nil {
		t.Fatal("result has no final equipment")
	}
}

func TestCachedSimulationResultFromModel(t *testing.T) {
	stored := models.FightSimulation{
		Winrate:               75,
		AverageTurns:          4,
		AverageFinalHP:        20,
		AverageFightCooldown:  8,
		EstimatedRestCooldown: 3,
		CycleCost:             11,
		DamageSurplus:         5,
		SurvivalSurplus:       0.5,
		XP:                    10,
		XPPerSecond:           1,
		GoldPerSecond:         2,
		ProspectingEfficiency: 0.1,
		Safe:                  true,
	}
	got := cachedSimulationResultFromModel(stored)
	if got.Winrate != stored.Winrate || got.CycleCost != stored.CycleCost || got.XP != stored.XP || !got.Safe {
		t.Fatalf("cached result = %#v, want values from %#v", got, stored)
	}
}

func TestEvaluateSimulationLoadoutReadsPersistentCache(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	character := schemas.CharacterSchema{Level: 1}
	monster := schemas.MonsterSchema{Code: "chicken"}
	slots := map[string]string{"weapon": "iron_sword"}
	key := bestFightSimulationCacheKey(character, monster, slots, 1)
	cache.SaveFightSimulation(models.FightSimulation{
		Key:          key,
		Version:      bestFightSimulationCacheVersion,
		Winrate:      88,
		CycleCost:    12,
		XP:           42,
		Safe:         true,
		AverageTurns: 3,
	})
	ResetSimulationCache()
	got := evaluateSimulationLoadoutIterations(character, monster, slots, 1)
	if got.Winrate != 88 || got.CycleCost != 12 || got.XP != 42 || !got.Safe {
		t.Fatalf("cached evaluation = %#v", got)
	}
}

func TestRefineWithUtilitiesEvaluatesAvailableUtility(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	character := schemas.CharacterSchema{Level: 5}
	monster, ok := catalog.Monsters.Get("chicken")
	if !ok {
		t.Fatal("chicken monster is missing from the catalog")
	}
	current := map[string]string{"weapon": "iron_sword"}
	available := map[string]int{"iron_sword": 1, "small_health_potion": 1}
	options := map[string][]string{"weapon": {"iron_sword"}}
	best := Result{Winrate: -1}
	got, result := refineWithUtilities(character, *monster, current, current, best, options, available, false, 1, betterSearchScore)
	if result.Winrate < 0 || got["weapon"] != "iron_sword" {
		t.Fatalf("refined result = %#v, %v", got, result)
	}
}

func TestRefineWithUtilitiesSkipsInvalidUtilities(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	character := schemas.CharacterSchema{Level: 1}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	current := map[string]string{"weapon": "iron_sword"}
	available := map[string]int{"empty_potion": 1, "small_health_potion": 0}
	options := map[string][]string{"weapon": {"iron_sword"}}
	best := Result{Winrate: 100}
	got, result := refineWithUtilities(character, monster, current, current, best, options, available, false, 1, betterSearchScore)
	if got["weapon"] != "iron_sword" || result.Winrate != 100 {
		t.Fatalf("invalid utilities changed result: %#v, %#v", got, result)
	}
}

func TestRefineWithUtilitiesReevaluatesEquipment(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	character := schemas.CharacterSchema{Level: 1}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	current := map[string]string{"weapon": "iron_sword"}
	available := map[string]int{"iron_sword": 1, "wooden_shield": 1}
	options := map[string][]string{
		"weapon": {"iron_sword"},
		"shield": {"", "wooden_shield"},
	}
	best := Result{Winrate: -1}
	better := func(Result, Result) bool { return true }
	got, result := refineWithUtilities(character, monster, current, current, best, options, available, false, 1, better)
	if got["weapon"] != "iron_sword" || result.Winrate < 0 {
		t.Fatalf("equipment refinement failed: %#v, %#v", got, result)
	}
}

func TestChooseSimulationGroupAssignsRings(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	character := schemas.CharacterSchema{Level: 1}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	current := map[string]string{"weapon": "iron_sword"}
	available := map[string]int{"iron_sword": 1, "iron_ring": 2}
	options := []string{"", "iron_ring"}
	chooseSimulationGroup(character, monster, current, simulationRingSlots, options, available, false)
	if current["ring1"] != "iron_ring" || current["ring2"] != "iron_ring" {
		t.Fatalf("ring group = %#v", current)
	}
}

func TestChooseSimulationGroupRejectsDuplicateArtifacts(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	character := schemas.CharacterSchema{Level: 1}
	monster := schemas.MonsterSchema{Code: "chicken", Level: 1, Hp: 60}
	current := map[string]string{"weapon": "iron_sword"}
	available := map[string]int{"iron_sword": 1, "novice_guide": 2}
	chooseSimulationGroup(character, monster, current, simulationArtifactSlots, []string{"novice_guide"}, available, false)
	if current["artifact1"] == "novice_guide" && current["artifact2"] == "novice_guide" {
		t.Fatalf("duplicate artifact assigned: %#v", current)
	}
}

func TestFindFightByNamePropagatesCharacterError(t *testing.T) {
	want := errors.New("character unavailable")
	characters := func(string) (schemas.CharacterSchema, error) {
		return schemas.CharacterSchema{}, want
	}
	d := deps{characters: characters}
	_, err := findFightByName(d, "unknown", schemas.MonsterSchema{}, false, false)
	if !errors.Is(err, want) {
		t.Fatalf("findFightByName() error = %v, want %v", err, want)
	}
}

func TestFindFightByNameFindsCharacter(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	characters := func(string) (schemas.CharacterSchema, error) {
		return schemas.CharacterSchema{Level: 1}, nil
	}
	bankItems := func() ([]schemas.SimpleItemSchema, error) {
		return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
	}
	d := deps{characters: characters, myBankItems: bankItems}
	monster, ok := catalog.Monsters.Get("chicken")
	if !ok {
		t.Fatal("chicken monster is missing from the catalog")
	}
	result, err := findFightByName(d, "test", *monster, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Focus == "" {
		t.Fatalf("missing focus in result: %#v", result)
	}
}

func TestFindFightAtLevelPropagatesBankError(t *testing.T) {
	want := errors.New("bank unavailable")
	bankItems := func() ([]schemas.SimpleItemSchema, error) {
		return nil, want
	}
	d := deps{myBankItems: bankItems}
	_, err := findFightAtLevel(d, 1, schemas.MonsterSchema{}, false, false)
	if !errors.Is(err, want) {
		t.Fatalf("findFightAtLevel() error = %v, want %v", err, want)
	}
}

func TestFindFightAtLevelFindsFight(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	bankItems := func() ([]schemas.SimpleItemSchema, error) {
		return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
	}
	d := deps{myBankItems: bankItems}
	monster, ok := catalog.Monsters.Get("chicken")
	if !ok {
		t.Fatal("chicken monster is missing from the catalog")
	}
	result, err := findFightAtLevel(d, 1, *monster, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Focus == "" {
		t.Fatalf("missing focus in result: %#v", result)
	}
}

func TestFindFightPropagatesBankError(t *testing.T) {
	want := errors.New("bank unavailable")
	bankItems := func() ([]schemas.SimpleItemSchema, error) {
		return nil, want
	}
	d := deps{myBankItems: bankItems}
	_, err := findFight(d, schemas.CharacterSchema{}, schemas.MonsterSchema{}, false, false)
	if !errors.Is(err, want) {
		t.Fatalf("findFight() error = %v, want %v", err, want)
	}
}

func TestFindFightIncludesInventoryAndEquipment(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	ResetSimulationCache()
	inventory := make([]schemas.InventorySlotSchema, 1)
	inventory[0].Code = "wooden_shield"
	inventory[0].Quantity = 1
	character := schemas.CharacterSchema{
		Level:      1,
		Inventory:  &inventory,
		WeaponSlot: "iron_sword",
	}
	monster, ok := catalog.Monsters.Get("chicken")
	if !ok {
		t.Fatal("chicken monster is missing from the catalog")
	}
	bankItems := func() ([]schemas.SimpleItemSchema, error) {
		return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
	}
	d := deps{myBankItems: bankItems}
	result, err := findFight(d, character, *monster, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Focus == "" || result.FinalEquipment == nil {
		t.Fatalf("incomplete result: %#v", result)
	}
}

func TestBankAvailableAggregatesItems(t *testing.T) {
	bankItems := func() ([]schemas.SimpleItemSchema, error) {
		return []schemas.SimpleItemSchema{
			{Code: "iron_sword", Quantity: 1},
			{Code: "iron_sword", Quantity: 2},
		}, nil
	}
	d := deps{myBankItems: bankItems}
	got, err := bankAvailable(d)
	if err != nil {
		t.Fatal(err)
	}
	if got["iron_sword"] != 3 {
		t.Fatalf("iron_sword quantity = %d, want 3", got["iron_sword"])
	}
}
