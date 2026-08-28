package surplus

import (
	"testing"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestDominatesAllElementalEffects(t *testing.T) {
	superior := itemWithEffects("superior", "attack_fire", 15, "attack_water", 8)
	inferior := itemWithEffects("inferior", "attack_fire", 10, "attack_water", 5)

	if !Dominates(superior, inferior) {
		t.Fatal("Dominates() = false, want true")
	}
	if Dominates(inferior, superior) {
		t.Fatal("Dominates() = true in reverse direction, want false")
	}
}

func TestDominatesConsidersAllEffects(t *testing.T) {
	superior := itemWithEffects("superior", "attack_fire", 15, "attack_water", 8, "wisdom", 10)
	inferior := itemWithEffects("inferior", "attack_fire", 10, "attack_water", 5, "wisdom", 10)

	if !Dominates(superior, inferior) {
		t.Fatal("Dominates() = false, want true")
	}
}

func TestDominatesTreatsGlobalDamageAsElementalDamage(t *testing.T) {
	global := itemWithEffects("global", "dmg", 10)
	fire := itemWithEffects("fire", "dmg_fire", 10)

	if !Dominates(global, fire) {
		t.Fatal("global damage should dominate equal fire damage")
	}
	if Dominates(fire, global) {
		t.Fatal("fire damage should not dominate global damage")
	}
}

func TestDominatesTreatsEquivalentGlobalAndElementalDamageAsEqual(t *testing.T) {
	global := itemWithEffects("global", "dmg", 10)
	elemental := itemWithEffects("elemental", "dmg_fire", 10, "dmg_water", 10, "dmg_earth", 10, "dmg_air", 10)

	if Dominates(global, elemental) || Dominates(elemental, global) {
		t.Fatal("equivalent global and elemental damage should not dominate each other")
	}
}

func TestDominatesRealGlobalDamageOverSingleElementalDamage(t *testing.T) {
	global := catalogItem(t, "gold_ring")
	elemental := catalogItem(t, "fire_ring")

	if !Dominates(global, elemental) {
		t.Fatal("gold_ring should dominate fire_ring through global damage")
	}
	if Dominates(elemental, global) {
		t.Fatal("fire_ring should not dominate gold_ring")
	}
}

func TestRealGlobalAndElementalDamageAreCombined(t *testing.T) {
	combined := catalogItem(t, "topaz_ring")
	values := effectValues(combined)

	if values["dmg_earth"] != 24 {
		t.Errorf("topaz_ring earth damage = %d, want 24", values["dmg_earth"])
	}
	if values["dmg_fire"] != 17 || values["dmg_water"] != 17 || values["dmg_air"] != 17 {
		t.Errorf("topaz_ring non-earth damage = %#v, want 17 for each element", values)
	}
}

func TestDominatesHandlesCooldownReduction(t *testing.T) {
	superior := itemWithEffects("superior", "woodcutting", -20)
	inferior := itemWithEffects("inferior", "woodcutting", -10)

	if !Dominates(superior, inferior) {
		t.Fatal("Dominates() = false, want true")
	}
	if Dominates(inferior, superior) {
		t.Fatal("Dominates() = true in reverse direction, want false")
	}
}

func TestDominatesRejectsWorseEffectOnlyOnSuperior(t *testing.T) {
	superior := itemWithEffects("superior", "alchemy", 1)
	inferior := itemWithEffects("inferior")

	if Dominates(superior, inferior) {
		t.Fatal("Dominates() = true for a worse superior-only effect, want false")
	}
}

func TestDominatesRejectsTradeoffs(t *testing.T) {
	first := itemWithEffects("first", "attack_fire", 15, "attack_water", 5)
	second := itemWithEffects("second", "attack_fire", 10, "attack_water", 8)

	if Dominates(first, second) || Dominates(second, first) {
		t.Fatal("Dominates() = true for a tradeoff, want false")
	}
}

func TestDominatesRealWeapons(t *testing.T) {
	goldSword := catalogItem(t, "gold_sword")
	ironSword := catalogItem(t, "iron_sword")
	fireBow := catalogItem(t, "fire_bow")
	waterBow := catalogItem(t, "water_bow")

	if !Dominates(goldSword, ironSword) {
		t.Error("gold_sword should dominate iron_sword")
	}
	if Dominates(fireBow, waterBow) || Dominates(waterBow, fireBow) {
		t.Error("fire_bow and water_bow should be a tradeoff")
	}
}

func TestDominatesRealTools(t *testing.T) {
	ironAxe := catalogItem(t, "iron_axe")
	copperAxe := catalogItem(t, "copper_axe")

	if !Dominates(ironAxe, copperAxe) {
		t.Error("iron_axe should dominate copper_axe")
	}
	if Dominates(copperAxe, ironAxe) {
		t.Error("copper_axe should not dominate iron_axe")
	}
}

func TestCanReplaceChecksTypeAndLevel(t *testing.T) {
	goldSword := catalogItem(t, "gold_sword")
	ironSword := catalogItem(t, "iron_sword")
	lowLevel := schemas.CharacterSchema{Level: 10}
	highLevel := schemas.CharacterSchema{Level: 30}

	if CanReplace(goldSword, ironSword, lowLevel) {
		t.Error("gold_sword should not replace iron_sword for a level 10 character")
	}
	if !CanReplace(goldSword, ironSword, highLevel) {
		t.Error("gold_sword should replace iron_sword for a level 30 character")
	}
}

func TestCanReplaceChecksToolSkill(t *testing.T) {
	ironAxe := catalogItem(t, "iron_axe")
	copperAxe := catalogItem(t, "copper_axe")
	character := schemas.CharacterSchema{WoodcuttingLevel: 10}

	if !CanReplace(ironAxe, copperAxe, character) {
		t.Error("iron_axe should replace copper_axe for woodcutting level 10")
	}
}

func TestCanReplaceDoesNotMixToolSkills(t *testing.T) {
	axe := catalogItem(t, "steel_axe")
	pickaxe := catalogItem(t, "steel_pickaxe")
	character := schemas.CharacterSchema{WoodcuttingLevel: 20, MiningLevel: 20}

	if CanReplace(axe, pickaxe, character) || CanReplace(pickaxe, axe, character) {
		t.Fatal("tools for different skills should not replace each other")
	}
}

func TestNonDominatedRemovesReplaceableItems(t *testing.T) {
	goldSword := catalogItem(t, "gold_sword")
	ironSword := catalogItem(t, "iron_sword")
	character := schemas.CharacterSchema{Level: 30}

	result := NonDominated([]schemas.ItemSchema{ironSword, goldSword}, character)
	if len(result) != 1 || result[0].Code != "gold_sword" {
		t.Fatalf("NonDominated() = %#v, want only gold_sword", result)
	}
}

func TestFindExcludesUtilities(t *testing.T) {
	results := Find(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "copper_axe", Quantity: 10},
			{Code: "small_health_potion", Quantity: 10},
		},
		Characters: []schemas.CharacterSchema{{Level: 10}},
	})
	if len(results) != 1 || results[0].Item.Code != "copper_axe" {
		t.Fatalf("Find() = %#v, want only copper_axe", results)
	}
}

func TestFindDominatedUsesTwoRingSlots(t *testing.T) {
	characters := []schemas.CharacterSchema{
		{Level: 20},
		{Level: 20},
		{Level: 20},
		{Level: 20},
		{Level: 20},
	}
	results := Find(Input{
		BankItems: []schemas.SimpleItemSchema{
			{Code: "iron_ring", Quantity: 10},
			{Code: "steel_ring", Quantity: 1},
		},
		Characters: characters,
	})
	found := false
	for _, result := range results {
		if result.Item.Code == "iron_ring" && result.Surplus != 1 {
			t.Fatalf("Find(iron_ring) = %d surplus, want 1", result.Surplus)
		} else if result.Item.Code == "iron_ring" {
			found = true
		}
	}
	if !found {
		t.Fatal("Find() did not return iron_ring")
	}
}

func TestFindUsesTwoPotentialSlotsForUnavailableRing(t *testing.T) {
	results := Find(Input{
		BankItems:  []schemas.SimpleItemSchema{{Code: "iron_ring", Quantity: 3}},
		Characters: []schemas.CharacterSchema{{Level: 1}},
	})

	if len(results) != 1 || results[0].Item.Code != "iron_ring" {
		t.Fatalf("Find() = %#v, want one iron_ring result", results)
	}
	if results[0].Surplus != 1 {
		t.Fatalf("Find(iron_ring) = %d surplus, want 1", results[0].Surplus)
	}
}

func catalogItem(t *testing.T, code string) schemas.ItemSchema {
	t.Helper()
	item, ok := database.Items().Get(code)
	if !ok {
		t.Fatalf("database.Items().Get(%q) did not find an item", code)
	}
	return *item
}

func itemWithEffects(code string, effects ...any) schemas.ItemSchema {
	list := make([]schemas.SimpleEffectSchema, 0, len(effects)/2)
	for i := 0; i < len(effects); i += 2 {
		list = append(list, schemas.SimpleEffectSchema{
			Code:  effects[i].(string),
			Value: effects[i+1].(int),
		})
	}
	return schemas.ItemSchema{Code: code, Effects: &list}
}
