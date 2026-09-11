package usage

import (
	"errors"
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestEvaluateRejectsUnknownOrNonEquipmentItems(t *testing.T) {
	for _, code := range []string{"missing_item", "minor_health_potion"} {
		_, err := evaluate(deps{}, []string{code}, false)
		if err == nil {
			t.Errorf("evaluate(%q) unexpectedly succeeded", code)
		}
	}
}

func TestEvaluateRejectsAPIErrorsAndEmptyAccounts(t *testing.T) {
	want := errors.New("characters unavailable")
	d := deps{
		accountsCharacters: func(string) ([]schemas.CharacterSchema, error) {
			return nil, want
		},
	}
	_, err := evaluate(d, nil, false)
	if !errors.Is(err, want) {
		t.Fatalf("character error = %v, want %v", err, want)
	}
	d.accountsCharacters = func(string) ([]schemas.CharacterSchema, error) {
		return []schemas.CharacterSchema{}, nil
	}
	_, err = evaluate(d, nil, false)
	if err == nil || err.Error() != "account has no characters" {
		t.Fatalf("empty account error = %v", err)
	}
}

func TestEvaluatePropagatesBankError(t *testing.T) {
	want := errors.New("bank unavailable")
	d := deps{
		accountsCharacters: func(string) ([]schemas.CharacterSchema, error) {
			return []schemas.CharacterSchema{{Level: 1}}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, want
		},
	}
	_, err := evaluate(d, nil, false)
	if !errors.Is(err, want) {
		t.Fatalf("evaluate bank error = %v, want %v", err, want)
	}
}

func TestEvaluatePropagatesNormalCombatError(t *testing.T) {
	want := errors.New("combat finder failed")
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	item := schemas.SimpleItemSchema{Code: "wooden_stick", Quantity: 1}
	d := evaluationTestDeps(1, []schemas.SimpleItemSchema{item})
	d.findFight = func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
		return best.Result{}, want
	}
	_, err = evaluate(d, []string{"wooden_stick"}, false)
	if !errors.Is(err, want) {
		t.Fatalf("evaluate combat error = %v, want %v", err, want)
	}
}

func TestEvaluatePropagatesNormalCraftingError(t *testing.T) {
	want := errors.New("crafting finder failed")
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	item := schemas.SimpleItemSchema{Code: "wooden_stick", Quantity: 1}
	d := evaluationTestDeps(1, []schemas.SimpleItemSchema{item})
	d.findEquipment = func(schemas.CharacterSchema, best.EquipmentOptions) (map[string]best.BestResult, error) {
		return nil, want
	}
	_, err = evaluate(d, []string{"wooden_stick"}, false)
	if !errors.Is(err, want) {
		t.Fatalf("evaluate crafting error = %v, want %v", err, want)
	}
}

func TestEvaluatePropagatesShortageCombatError(t *testing.T) {
	want := errors.New("shortage combat failed")
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	item := schemas.SimpleItemSchema{Code: "wooden_stick", Quantity: 1}
	d := evaluationTestDeps(1, []schemas.SimpleItemSchema{item})
	d.findFight = func(_ schemas.CharacterSchema, _ schemas.MonsterSchema, available map[string]int, _, _ bool) (best.Result, error) {
		if available["wooden_stick"] == 0 {
			return best.Result{}, want
		}
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "wooden_stick"},
		}, nil
	}
	_, err = evaluate(d, []string{"wooden_stick"}, false)
	if !errors.Is(err, want) {
		t.Fatalf("evaluate shortage combat error = %v, want %v", err, want)
	}
}

func TestEvaluatePropagatesShortageCraftingError(t *testing.T) {
	want := errors.New("shortage crafting failed")
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	item := schemas.SimpleItemSchema{Code: "wooden_stick", Quantity: 1}
	d := evaluationTestDeps(1, []schemas.SimpleItemSchema{item})
	craftingCalls := 0
	d.findEquipment = func(_ schemas.CharacterSchema, _ best.EquipmentOptions) (map[string]best.BestResult, error) {
		craftingCalls++
		if craftingCalls > 2 {
			return nil, want
		}
		return map[string]best.BestResult{}, nil
	}
	_, err = evaluate(d, []string{"wooden_stick"}, false)
	if !errors.Is(err, want) {
		t.Fatalf("evaluate shortage crafting error = %v, want %v", err, want)
	}
}

func TestEvaluateAcceptsDuplicateEquipmentCodesBeforeAccountLookup(t *testing.T) {
	called := false
	d := deps{
		accountsCharacters: func(string) ([]schemas.CharacterSchema, error) {
			called = true
			return nil, errors.New("stop")
		},
	}
	_, err := evaluate(d, []string{"wooden_stick", "wooden_stick"}, false)
	if !called || err == nil {
		t.Fatalf("duplicate equipment handling failed: called=%v error=%v", called, err)
	}
}

func TestGlobalOwnedCollectsBankLoadoutsAndInventory(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "wooden_stick", Quantity: 2},
	}
	characters := []schemas.CharacterSchema{
		{WeaponSlot: "iron_sword", Inventory: &inventory},
	}
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 3}}, nil
		},
	}
	owned, err := globalOwned(d, characters)
	if err != nil || owned["iron_sword"] != 4 || owned["wooden_stick"] != 2 {
		t.Fatalf("owned/error = %v/%v", owned, err)
	}
	want := errors.New("bank unavailable")
	d.myBankItems = func() ([]schemas.SimpleItemSchema, error) { return nil, want }
	_, err = globalOwned(d, characters)
	if !errors.Is(err, want) {
		t.Fatalf("globalOwned error = %v, want %v", err, want)
	}
}

func TestMarkCombatAndMergeUsage(t *testing.T) {
	monsters := []*schemas.MonsterSchema{{Code: "chicken"}, {Code: "rat"}}
	d := deps{
		findFight: func(_ schemas.CharacterSchema, monster schemas.MonsterSchema, _ map[string]int, _, _ bool) (best.Result, error) {
			if monster.Code == "rat" {
				return best.Result{Winrate: 50}, nil
			}
			return best.Result{
				Winrate:        100,
				FinalEquipment: map[string]string{"weapon": "wooden_stick"},
			}, nil
		},
	}
	loadouts, err := markCombat(d, schemas.CharacterSchema{}, monsters, nil, "normal")
	if err != nil || len(loadouts) != 1 {
		t.Fatalf("loadouts/error = %v/%v", loadouts, err)
	}
	result := map[string]Evaluation{"wooden_stick": {}}
	mergeCombatUsage(result, loadouts, true)
	if !result["wooden_stick"].Needed || len(result["wooden_stick"].Monsters) != 1 {
		t.Fatalf("merged usage = %+v", result)
	}
	result["iron_sword"] = Evaluation{}
	mergeCombatUsage(result, loadouts, true)
	if result["iron_sword"].Needed {
		t.Fatal("non-selected equipment was marked as needed")
	}
}

func TestMarkCraftingUsesSelectedEquipment(t *testing.T) {
	d := deps{
		findEquipment: func(_ schemas.CharacterSchema, options best.EquipmentOptions) (map[string]best.BestResult, error) {
			return map[string]best.BestResult{
				"weapon": {Code: options.Priorities[0]},
			}, nil
		},
	}
	result := map[string]Evaluation{"wisdom": {}, "prospecting": {}}
	err := markCrafting(d, schemas.CharacterSchema{}, nil, result, true)
	if err != nil || !result["wisdom"].Needed || !result["prospecting"].Needed {
		t.Fatalf("crafting result/error = %+v/%v", result, err)
	}
}

func TestMarkCombatPropagatesFinderError(t *testing.T) {
	want := errors.New("finder failed")
	d := deps{
		findFight: func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
			return best.Result{}, want
		},
	}
	monsters := []*schemas.MonsterSchema{{Code: "chicken"}}
	_, err := markCombat(d, schemas.CharacterSchema{}, monsters, nil, "normal")
	if !errors.Is(err, want) {
		t.Fatalf("markCombat error = %v, want %v", err, want)
	}
}

func TestMarkCraftingPropagatesFinderError(t *testing.T) {
	want := errors.New("equipment finder failed")
	d := deps{
		findEquipment: func(schemas.CharacterSchema, best.EquipmentOptions) (map[string]best.BestResult, error) {
			return nil, want
		},
	}
	err := markCrafting(d, schemas.CharacterSchema{}, nil, map[string]Evaluation{}, false)
	if !errors.Is(err, want) {
		t.Fatalf("markCrafting error = %v, want %v", err, want)
	}
}

func TestCachedMarkCombatReadsStoredResults(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	d := deps{
		findFight: func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
			return best.Result{
				Winrate:        100,
				FinalEquipment: map[string]string{"weapon": "wooden_stick"},
			}, nil
		},
	}
	monsters := []*schemas.MonsterSchema{{Code: "chicken"}}
	first, err := cachedMarkCombat(d, schemas.CharacterSchema{}, monsters, map[string]int{}, "normal")
	if err != nil || len(first) != 1 {
		t.Fatalf("first cache result/error = %v/%v", first, err)
	}
	d.findFight = func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
		return best.Result{}, errors.New("cache miss")
	}
	second, err := cachedMarkCombat(d, schemas.CharacterSchema{}, monsters, map[string]int{}, "normal")
	if err != nil || len(second) != 1 {
		t.Fatalf("stored cache result/error = %v/%v", second, err)
	}
}

func TestEvaluateRunsCombatAndCraftingScenarios(t *testing.T) {
	d := deps{
		accountsCharacters: func(string) ([]schemas.CharacterSchema, error) {
			return []schemas.CharacterSchema{
				{Level: 10, WeaponSlot: "wooden_stick"},
				{Level: 5},
			}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			item := schemas.SimpleItemSchema{Code: "wooden_stick", Quantity: 1}
			items := []schemas.SimpleItemSchema{item}
			return items, nil
		},
		findFight: func(_ schemas.CharacterSchema, _ schemas.MonsterSchema, _ map[string]int, _, _ bool) (best.Result, error) {
			return best.Result{
				Winrate:        100,
				FinalEquipment: map[string]string{"weapon": "wooden_stick"},
			}, nil
		},
		findEquipment: func(_ schemas.CharacterSchema, options best.EquipmentOptions) (map[string]best.BestResult, error) {
			return map[string]best.BestResult{"weapon": {Code: "wooden_stick"}}, nil
		},
	}
	result, err := evaluate(d, []string{"wooden_stick"}, true)
	if err != nil {
		t.Fatalf("evaluate returned error: %v", err)
	}
	usage, ok := result["wooden_stick"]
	if !ok || !usage.Needed || len(usage.Monsters) == 0 || len(usage.Effects) == 0 {
		t.Fatalf("evaluated usage = %+v", result)
	}
}

func TestEvaluateUsesAllOwnedEquipmentWithoutExplicitCodes(t *testing.T) {
	item := schemas.SimpleItemSchema{Code: "wooden_stick", Quantity: 1}
	d := evaluationTestDeps(10, []schemas.SimpleItemSchema{item})
	result, err := evaluate(d, nil, false)
	if err != nil || result["wooden_stick"].Needed != true {
		t.Fatalf("automatic evaluation = %+v, error=%v", result, err)
	}
}

func TestEvaluateAddsExplicitUnownedEquipmentToCombatAvailability(t *testing.T) {
	d := evaluationTestDeps(1, nil)
	result, err := evaluate(d, []string{"wooden_stick"}, false)
	if err != nil {
		t.Fatalf("explicit evaluation returned error: %v", err)
	}
	_, ok := result["wooden_stick"]
	if !ok {
		t.Fatalf("explicit equipment was not evaluated: %+v", result)
	}
}

func TestEvaluateSkipsEquipmentUnavailableAtCharacterLevel(t *testing.T) {
	d := evaluationTestDeps(1, nil)
	result, err := evaluate(d, []string{"iron_sword"}, false)
	if err != nil {
		t.Fatalf("level-filtered evaluation returned error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("unavailable equipment was evaluated: %+v", result)
	}
}

func evaluationTestDeps(level int, bank []schemas.SimpleItemSchema) deps {
	return deps{
		accountsCharacters: func(string) ([]schemas.CharacterSchema, error) {
			return []schemas.CharacterSchema{{Level: level}}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) { return bank, nil },
		findFight: func(_ schemas.CharacterSchema, _ schemas.MonsterSchema, _ map[string]int, _, _ bool) (best.Result, error) {
			return best.Result{
				Winrate:        100,
				FinalEquipment: map[string]string{"weapon": "wooden_stick"},
			}, nil
		},
		findEquipment: func(_ schemas.CharacterSchema, _ best.EquipmentOptions) (map[string]best.BestResult, error) {
			return map[string]best.BestResult{}, nil
		},
	}
}

func TestCanonicalAvailableNormalizesQuantitiesAndSorts(t *testing.T) {
	available := map[string]int{
		"iron_sword":          20,
		"ring_of_the_adept":   20,
		"small_health_potion": 3,
		"missing_item":        4,
		"empty":               0,
	}
	got := canonicalAvailable(available)
	if !strings.Contains(got, "iron_sword=5") {
		t.Fatalf("canonical available = %q", got)
	}
	if !strings.Contains(got, "ring_of_the_adept=6") {
		t.Fatalf("canonical available = %q", got)
	}
	if !strings.Contains(got, "small_health_potion=1") {
		t.Fatalf("canonical available = %q", got)
	}
	if !strings.Contains(got, "missing_item=4") {
		t.Fatalf("canonical available = %q", got)
	}
	if strings.Contains(got, "empty") || !strings.HasPrefix(got, "iron_sword") {
		t.Fatalf("canonical available is not sorted/filtered: %q", got)
	}
}

func TestEquipmentQuantityLimit(t *testing.T) {
	ring := schemas.ItemSchema{Type: "ring"}
	weapon := schemas.ItemSchema{Type: "weapon"}
	if equipmentQuantityLimit("ring_of_the_adept", ring) != 6 {
		t.Fatal("adept ring limit is not six")
	}
	if equipmentQuantityLimit("ring_of_chance", ring) != 10 {
		t.Fatal("ring limit is not ten")
	}
	if equipmentQuantityLimit("iron_sword", weapon) != 5 {
		t.Fatal("equipment limit is not five")
	}
}

func TestCombatCacheKeyIsStableForMapAndMonsterOrder(t *testing.T) {
	first := []*schemas.MonsterSchema{{Code: "z"}, {Code: "a"}}
	second := []*schemas.MonsterSchema{{Code: "a"}, {Code: "z"}}
	available := map[string]int{"iron_sword": 1, "ring_of_chance": 2}
	if combatCacheKey(10, first, available) != combatCacheKey(10, second, available) {
		t.Fatal("cache key depends on monster order")
	}
}

func TestEquipmentAndCombatItemFilters(t *testing.T) {
	character := schemas.CharacterSchema{Level: 20}
	owned := map[string]int{
		"wooden_stick":        1,
		"minor_health_potion": 1,
		"iron_pickaxe":        1,
		"missing_item":        1,
		"iron_shield":         0,
	}
	codes := equipmentCodes(character, owned)
	if len(codes) != 2 || !contains(codes, "wooden_stick") {
		t.Fatalf("equipment codes = %v", codes)
	}
	combat := combatItems(character, owned)
	if combat["wooden_stick"] != 1 || combat["minor_health_potion"] != 1 {
		t.Fatalf("combat items = %v", combat)
	}
	_, ok := combat["iron_pickaxe"]
	if ok {
		t.Fatal("tool was included in combat items")
	}
	if !isEquipment(*mustItem(t, "wooden_stick")) || isEquipment(*mustItem(t, "minor_health_potion")) {
		t.Fatal("equipment classification is incorrect")
	}
}

func TestLoadoutAndCountHelpers(t *testing.T) {
	character := schemas.CharacterSchema{
		WeaponSlot:   "iron_sword",
		Utility1Slot: "small_health_potion",
	}
	loadout := characterLoadout(character)
	if !loadoutContains(loadout, "iron_sword") || loadoutContains(loadout, "missing") {
		t.Fatal("loadout containment is incorrect")
	}
	if !contains([]string{"a", "b"}, "b") || contains([]string{"a"}, "c") {
		t.Fatal("slice containment is incorrect")
	}
	cloned := cloneCounts(map[string]int{"iron_sword": 2})
	cloned["iron_sword"] = 3
	if cloned["iron_sword"] != 3 {
		t.Fatal("clone counts was not writable")
	}
}

func TestInsufficientCodesUsesEquipmentLimits(t *testing.T) {
	character := schemas.CharacterSchema{Level: 50}
	owned := map[string]int{
		"iron_sword":        1,
		"ring_of_chance":    1,
		"ring_of_the_adept": 1,
	}
	codes := insufficientCodes(character, owned)
	if len(codes) != 3 {
		t.Fatalf("insufficient codes = %v, want three codes", codes)
	}
}

func mustItem(t *testing.T, code string) *schemas.ItemSchema {
	t.Helper()
	item, ok := database.Items().Get(code)
	if !ok {
		t.Fatalf("catalog item %q is missing", code)
	}
	return item
}
