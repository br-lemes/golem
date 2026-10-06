package potential

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestEvaluateRejectsInvalidAndExcludedItems(t *testing.T) {
	for _, test := range []struct {
		code    string
		without []string
	}{
		{code: "missing"},
		{code: "small_health_potion"},
		{code: "wooden_stick", without: []string{"wooden_stick"}},
	} {
		_, err := evaluate(deps{}, []string{test.code}, test.without, false, 0, 0)
		if err == nil {
			t.Errorf("evaluate(%q, %v) unexpectedly succeeded", test.code, test.without)
		}
	}
}

func TestEvaluateRejectsItemsThatNoCharacterCanEquip(t *testing.T) {
	d := deps{
		accountsCharacters: func(string) ([]schemas.CharacterSchema, error) {
			return []schemas.CharacterSchema{{Level: 1}}, nil
		},
	}
	_, err := evaluate(d, []string{"iron_sword"}, nil, false, 0, 0)
	if err == nil {
		t.Fatal("evaluate unexpectedly succeeded")
	}
}

func TestEvaluateReportsUsage(t *testing.T) {
	d := testDeps(t, 1)
	d.findFight = func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "wooden_stick"},
		}, nil
	}
	d.findEquipment = func(schemas.CharacterSchema, best.EquipmentOptions) (map[string]best.BestResult, error) {
		return map[string]best.BestResult{"weapon": {Code: "wooden_stick"}}, nil
	}
	results, err := evaluate(d, []string{"wooden_stick"}, nil, false, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	result := results["wooden_stick"]
	if !result.Used || len(result.Uses) != 1 || len(result.SelectedEquipment) != 0 {
		t.Fatalf("result = %+v", result)
	}
	use := result.Uses[0]
	if use.CharacterLevel != 1 || use.EquipmentLevel != 1 || len(use.Monsters) == 0 || len(use.Effects) != 2 || use.Loadouts != nil {
		t.Fatalf("result = %+v", result)
	}
}

func TestEvaluateIncludesLoadoutsWithDetails(t *testing.T) {
	d := testDeps(t, 1)
	d.findFight = func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "wooden_stick"},
		}, nil
	}
	results, err := evaluate(d, []string{"wooden_stick"}, nil, true, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	result := results["wooden_stick"]
	if len(result.Uses) != 1 || len(result.Uses[0].Loadouts) == 0 {
		t.Fatalf("result = %+v", result)
	}
	for _, loadout := range result.Uses[0].Loadouts {
		if loadout["weapon"] != "wooden_stick" {
			t.Fatalf("loadout = %+v", loadout)
		}
	}
}

func TestEvaluateReportsAlternativesWhenUnused(t *testing.T) {
	d := testDeps(t, 1)
	d.findFight = func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "iron_sword"},
		}, nil
	}
	d.findEquipment = func(schemas.CharacterSchema, best.EquipmentOptions) (map[string]best.BestResult, error) {
		return map[string]best.BestResult{"weapon": {Code: "iron_sword"}}, nil
	}
	results, err := evaluate(d, []string{"wooden_stick"}, nil, false, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	result := results["wooden_stick"]
	if result.Used || len(result.SelectedEquipment) != 1 || result.SelectedEquipment[0].Code != "iron_sword" || len(result.SelectedEquipment[0].Uses) != 1 || len(result.SelectedEquipment[0].Uses[0].Monsters) == 0 {
		t.Fatalf("result = %+v", result)
	}
}

func TestEvaluateExcludesUnavailableItems(t *testing.T) {
	d := testDeps(t, 10)
	d.findFight = func(_ schemas.CharacterSchema, _ schemas.MonsterSchema, available map[string]int, _ bool, _ bool) (best.Result, error) {
		if available["wooden_stick"] != 0 {
			return best.Result{}, errors.New("excluded item is available")
		}
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "iron_sword"},
		}, nil
	}
	d.findEquipment = func(_ schemas.CharacterSchema, options best.EquipmentOptions) (map[string]best.BestResult, error) {
		if options.Owned["wooden_stick"] != 0 {
			return nil, errors.New("excluded item is available")
		}
		return map[string]best.BestResult{}, nil
	}
	_, err := evaluate(d, []string{"iron_sword"}, []string{"wooden_stick"}, false, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEvaluateSeparatesCharacterAndGearLevels(t *testing.T) {
	d := testDeps(t, 10)
	d.findFight = func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "water_bow"},
		}, nil
	}
	results, err := evaluate(d, []string{"water_bow"}, nil, false, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	result := results["water_bow"]
	if len(result.Uses) != 1 || result.Uses[0].CharacterLevel != 10 || result.Uses[0].EquipmentLevel != 5 {
		t.Fatalf("result = %+v", result)
	}
}

func TestEvaluateIncludesTargetAboveGearLevel(t *testing.T) {
	d := testDeps(t, 10)
	d.findFight = func(_ schemas.CharacterSchema, _ schemas.MonsterSchema, available map[string]int, _ bool, _ bool) (best.Result, error) {
		if available["iron_sword"] == 0 {
			return best.Result{}, errors.New("evaluated item is unavailable")
		}
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "iron_sword"},
		}, nil
	}
	results, err := evaluate(d, []string{"iron_sword"}, nil, false, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	result := results["iron_sword"]
	if !result.Used || result.Uses[0].EquipmentLevel != 5 {
		t.Fatalf("result = %+v", result)
	}
}

func TestEvaluateMakesAllTargetItemsAvailable(t *testing.T) {
	d := testDeps(t, 10)
	d.findFight = func(_ schemas.CharacterSchema, _ schemas.MonsterSchema, available map[string]int, _ bool, _ bool) (best.Result, error) {
		if available["iron_sword"] == 0 || available["water_bow"] == 0 {
			return best.Result{}, errors.New("evaluated items are unavailable")
		}
		return best.Result{
			Winrate:        100,
			FinalEquipment: map[string]string{"weapon": "water_bow"},
		}, nil
	}
	results, err := evaluate(d, []string{"iron_sword", "water_bow"}, nil, false, 0, 5)
	if err != nil {
		t.Fatal(err)
	}
	if !results["water_bow"].Used || results["iron_sword"].Used {
		t.Fatalf("results = %+v", results)
	}
}

func TestCachedCombatReusesStoredLoadouts(t *testing.T) {
	_ = testDeps(t, 1)
	calls := 0
	d := deps{
		findFight: func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
			calls++
			return best.Result{
				Winrate:        100,
				FinalEquipment: map[string]string{"weapon": "wooden_stick"},
			}, nil
		},
	}
	character := schemas.CharacterSchema{Level: 1}
	monsters := []*schemas.MonsterSchema{{Code: "chicken"}}
	available := map[string]int{"wooden_stick": 5}
	_, err := cachedCombat(d, character, monsters, available)
	if err != nil || calls != 1 {
		t.Fatalf("first cachedCombat calls/error = %d/%v", calls, err)
	}
	d.findFight = func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
		return best.Result{}, errors.New("cache miss")
	}
	_, err = cachedCombat(d, character, monsters, available)
	if err != nil || calls != 1 {
		t.Fatalf("second cachedCombat calls/error = %d/%v", calls, err)
	}
}

func testDeps(t *testing.T, level int) deps {
	t.Helper()
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	return deps{
		accountsCharacters: func(string) ([]schemas.CharacterSchema, error) {
			return []schemas.CharacterSchema{{Level: level}}, nil
		},
		findEquipment: func(schemas.CharacterSchema, best.EquipmentOptions) (map[string]best.BestResult, error) {
			return map[string]best.BestResult{}, nil
		},
		findFight: func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error) {
			return best.Result{}, nil
		},
	}
}
