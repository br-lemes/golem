package surplus

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFindSurplusTools(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{WoodcuttingLevel: 1},
		{WoodcuttingLevel: 10},
		{WoodcuttingLevel: 10},
		{WoodcuttingLevel: 10},
		{WoodcuttingLevel: 10},
	}
	got := Find(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "copper_axe", Quantity: 5},
			{Code: "iron_axe", Quantity: 2},
		},
		Characters: characters,
	})
	results := resultsByCode(got)
	if results["copper_axe"].Total != 5 || results["copper_axe"].Surplus != 2 {
		t.Fatalf("Find(copper_axe) returned total %d and surplus %d, want 5 and 2", results["copper_axe"].Total, results["copper_axe"].Surplus)
	}
}

func TestFindAllSkillsAndLocations(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{
			AlchemyLevel:     1,
			FishingLevel:     1,
			Inventory:        inventory("iron_axe"),
			MiningLevel:      1,
			WoodcuttingLevel: 1,
		},
		{
			AlchemyLevel:     10,
			FishingLevel:     10,
			Inventory:        inventory("iron_pickaxe"),
			MiningLevel:      10,
			WoodcuttingLevel: 10,
		},
		{
			AlchemyLevel:     20,
			FishingLevel:     20,
			Inventory:        inventory("spruce_fishing_rod"),
			MiningLevel:      20,
			WeaponSlot:       "gold_fishing_rod",
			WoodcuttingLevel: 20,
		},
		{
			AlchemyLevel:     30,
			FishingLevel:     30,
			Inventory:        inventory("leather_gloves"),
			MiningLevel:      30,
			WeaponSlot:       "golden_gloves",
			WoodcuttingLevel: 30,
		},
		{
			AlchemyLevel:     40,
			FishingLevel:     40,
			MiningLevel:      40,
			WeaponSlot:       "gold_axe",
			WoodcuttingLevel: 40,
		},
	}
	got := Find(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "copper_axe", Quantity: 5},
			{Code: "copper_pickaxe", Quantity: 5},
			{Code: "fishing_net", Quantity: 5},
			{Code: "apprentice_gloves", Quantity: 5},
		},
		Characters: characters,
	})
	want := map[string]struct{ total, surplus int }{
		"copper_axe":        {5, 2},
		"copper_pickaxe":    {5, 1},
		"fishing_net":       {5, 2},
		"apprentice_gloves": {5, 2},
	}
	assertResults(t, got, want)
}

func TestFindSurplusExcessOfHighestTool(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
	}
	got := Find(Input{
		BankItems:  []schemas.SimpleItemSchema{{Code: "steel_axe", Quantity: 6}},
		Characters: characters,
	})
	if len(got) != 1 || got[0].Total != 6 || got[0].Surplus != 1 {
		t.Fatalf("Find() returned %#v, want one steel_axe with total 6 and surplus 1", got)
	}
}

func TestFindKeepsToolsAboveCharacterLevels(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
		{WoodcuttingLevel: 20},
	}
	got := Find(Input{
		BankItems:  []schemas.SimpleItemSchema{{Code: "gold_axe", Quantity: 6}},
		Characters: characters,
	})
	if len(got) != 1 || got[0].Total != 6 || got[0].Surplus != 1 {
		t.Fatalf("Find() returned %#v, want one gold_axe with total 6 and surplus 1", got)
	}
}

func TestEvaluateUnownedTool(t *testing.T) {

	result := Evaluate(Input{
		BankItems:  []schemas.SimpleItemSchema{{Code: "copper_axe", Quantity: 1}},
		Characters: []schemas.CharacterSchema{{WoodcuttingLevel: 10}},
	}, "iron_axe")
	if result.Status != "not_dominated" {
		t.Fatalf("Evaluate() status = %q, want not_dominated", result.Status)
	}
	if len(result.ComparedTo) != 1 || result.ComparedTo[0].Code != "copper_axe" {
		t.Fatalf("Evaluate() comparisons = %#v, want copper_axe", result.ComparedTo)
	}
}

func TestExplainNotOwned(t *testing.T) {
	result := Explain(Input{}, "iron_axe")

	if result.Status != "not_owned" {
		t.Fatalf("Explain() status = %q, want not_owned", result.Status)
	}
}

func TestExplainFutureItem(t *testing.T) {
	result := Explain(Input{
		BankItems:  []schemas.SimpleItemSchema{{Code: "iron_axe", Quantity: 1}},
		Characters: []schemas.CharacterSchema{{WoodcuttingLevel: 1}},
	}, "iron_axe")

	if result.Status != "future" {
		t.Fatalf("Explain() status = %q, want future", result.Status)
	}
}

func TestExplainSurplusItem(t *testing.T) {
	result := Explain(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "copper_axe", Quantity: 5},
			{Code: "iron_axe", Quantity: 2},
		},
		Characters: []schemas.CharacterSchema{
			{WoodcuttingLevel: 1},
			{WoodcuttingLevel: 10},
			{WoodcuttingLevel: 10},
			{WoodcuttingLevel: 10},
			{WoodcuttingLevel: 10},
		},
	}, "copper_axe")

	if result.Status != "surplus" || result.Surplus != 2 {
		t.Fatalf("Explain() = status %q, surplus %d; want surplus, 2", result.Status, result.Surplus)
	}
	if len(result.DominatedBy) != 1 || result.DominatedBy[0].Code != "iron_axe" {
		t.Fatalf("Explain() dominated by %#v, want iron_axe", result.DominatedBy)
	}
}

func TestExplainNotDominatedWithoutCompatibleItems(t *testing.T) {
	result := Explain(Input{
		BankItems:  []schemas.SimpleItemSchema{{Code: "copper_axe", Quantity: 1}},
		Characters: []schemas.CharacterSchema{{WoodcuttingLevel: 10}},
	}, "copper_axe")

	if result.Status != "not_dominated" || result.Reason != "no other owned item has the same type and subtype" {
		t.Fatalf("Explain() = status %q, reason %q; want not_dominated with no-compatible reason", result.Status, result.Reason)
	}
}

func TestEvaluateUnknownItem(t *testing.T) {
	result := Evaluate(Input{}, "unknown_item")

	if result.Status != "unknown" {
		t.Fatalf("Evaluate() status = %q, want unknown", result.Status)
	}
}

func TestExplainComparisonEffectsUseNumericDirection(t *testing.T) {
	result := Explain(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "life_ring", Quantity: 1},
			{Code: "earth_ring", Quantity: 1},
		},
		Characters: []schemas.CharacterSchema{{Level: 20}},
	}, "life_ring")
	if len(result.ComparedTo) == 0 {
		t.Fatal("Explain() returned no comparisons")
	}
	effects := result.ComparedTo[0].Effects
	if effects["dmg_earth"] != "0 < 8 (worse)" {
		t.Fatalf("dmg_earth comparison = %q, want 0 < 8 (worse)", effects["dmg_earth"])
	}
	if effects["hp"] != "25 > 0 (better)" || effects["wisdom"] != "20 > 0 (better)" {
		t.Fatalf("positive effects = %#v, want hp 25 > 0 (better) and wisdom 20 > 0 (better)", effects)
	}
}

func TestExplainComparisonMarksBetterNegativeCooldown(t *testing.T) {
	result := Explain(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "copper_axe", Quantity: 1},
			{Code: "iron_axe", Quantity: 1},
		},
		Characters: []schemas.CharacterSchema{
			{WoodcuttingLevel: 20},
			{WoodcuttingLevel: 20},
		},
	}, "iron_axe")
	if len(result.ComparedTo) == 0 {
		t.Fatal("Explain() returned no comparisons")
	}
	got := result.ComparedTo[0].Effects["woodcutting"]
	if got != "-20 < -10 (better)" {
		t.Fatalf("woodcutting comparison = %q, want -20 < -10 (better)", got)
	}
}

func TestExplainSortsMultipleComparisons(t *testing.T) {
	result := Explain(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "life_ring", Quantity: 1},
			{Code: "earth_ring", Quantity: 1},
			{Code: "steel_ring", Quantity: 1},
		},
		Characters: []schemas.CharacterSchema{{Level: 20}},
	}, "life_ring")

	if len(result.ComparedTo) != 2 {
		t.Fatalf("Explain() returned %d comparisons, want 2", len(result.ComparedTo))
	}
	if result.ComparedTo[0].Code != "earth_ring" || result.ComparedTo[1].Code != "steel_ring" {
		t.Fatalf("Explain() comparisons = %#v, want earth_ring, steel_ring", result.ComparedTo)
	}
}

func resultsByCode(results []Result) map[string]Result {
	byCode := make(map[string]Result, len(results))
	for _, result := range results {
		byCode[result.Item.Code] = result
	}
	return byCode
}

func assertResults(t *testing.T, got []Result, want map[string]struct{ total, surplus int }) {
	t.Helper()
	results := resultsByCode(got)
	if len(results) != len(want) {
		t.Fatalf("returned %d items, want %d: %#v", len(results), len(want), results)
	}
	for code, expected := range want {
		actual, ok := results[code]
		if !ok {
			t.Fatalf("did not return %q", code)
		}
		if actual.Total != expected.total || actual.Surplus != expected.surplus {
			t.Errorf("%s = total %d, surplus %d; want total %d, surplus %d", code, actual.Total, actual.Surplus, expected.total, expected.surplus)
		}
	}
}

func inventory(code string) *[]schemas.InventorySlotSchema {
	items := []schemas.InventorySlotSchema{{Code: code, Quantity: 1}}
	return &items
}
