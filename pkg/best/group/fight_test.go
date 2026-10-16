package group

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFindFightWithAvailableRejectsInvalidParty(t *testing.T) {
	_, err := FindFightWithAvailable(nil, schemas.MonsterSchema{}, map[string]int{})
	if err == nil {
		t.Fatal("empty party accepted")
	}
	party := make([]schemas.CharacterSchema, 4)
	for i := range party {
		party[i].Name = string(rune('a' + i))
	}
	_, err = FindFightWithAvailable(party, schemas.MonsterSchema{}, map[string]int{})
	if err == nil {
		t.Fatal("four-character party accepted")
	}
}

func TestValidLoadoutsUsesGlobalQuantities(t *testing.T) {
	party := []schemas.CharacterSchema{{Name: "one"}, {Name: "two"}}
	loadouts := []map[string]string{
		{"weapon": "wooden_stick"},
		{"weapon": "wooden_stick"},
	}
	if validLoadouts(party, loadouts, map[string]int{"wooden_stick": 1}) {
		t.Fatal("one item was assigned to two characters")
	}
	if !validLoadouts(party, loadouts, map[string]int{"wooden_stick": 2}) {
		t.Fatal("two items were not assigned across two characters")
	}
}

func TestValidLoadoutsRejectsDuplicateArtifactsAndAdeptRings(t *testing.T) {
	party := []schemas.CharacterSchema{{Name: "one"}}
	loadouts := []map[string]string{{
		"artifact1": "artifact",
		"artifact2": "artifact",
		"ring1":     "ring_of_the_adept",
		"ring2":     "ring_of_the_adept",
	}}
	available := map[string]int{"artifact": 2, "ring_of_the_adept": 2}
	if validLoadouts(party, loadouts, available) {
		t.Fatal("duplicate artifact and adept rings were accepted")
	}
}

func TestUtilityOptionsAlwaysUseMaximumQuantity(t *testing.T) {
	options := utilityOptions(schemas.CharacterSchema{Level: 50}, map[string]int{
		"health_potion": 1,
	})
	if len(options) != 2 {
		t.Fatalf("utility options = %#v, want empty and health potion", options)
	}
	if options[1] != (Utility{Code: "health_potion", Quantity: 100}) {
		t.Fatalf("utility option = %#v, want 100 health potions", options[1])
	}
}

func TestValidAllocationAllowsSharedUtilitiesAcrossParty(t *testing.T) {
	current := allocation{
		equipment: []map[string]string{{}, {}},
		utilities: []map[string]Utility{
			{"utility1": {Code: "health_potion", Quantity: 100}},
			{"utility1": {Code: "health_potion", Quantity: 100}},
		},
	}
	if !validAllocation(current, map[string]int{"health_potion": 1}) {
		t.Fatal("shared utility was rejected across party members")
	}
}

func TestValidAllocationRejectsDuplicateUtilityForCharacter(t *testing.T) {
	current := allocation{
		equipment: []map[string]string{{}},
		utilities: []map[string]Utility{{
			"utility1": {Code: "health_potion", Quantity: 1},
			"utility2": {Code: "health_potion", Quantity: 1},
		}},
	}
	if validAllocation(current, map[string]int{"health_potion": 2}) {
		t.Fatal("same utility was assigned to both slots")
	}
}

func TestBetterPrefersDamageBeforeFirstWin(t *testing.T) {
	damage := score{winrate: 0, damage: 10, hp: 0, turns: 1}
	survival := score{winrate: 0, damage: 5, hp: 100, turns: 100}
	if !better(damage, survival) {
		t.Fatal("higher damage was not preferred before the first win")
	}
}

func TestRoleSeedsCombineTankAndSupport(t *testing.T) {
	party := []schemas.CharacterSchema{
		{Name: "one", Level: 50, WeaponSlot: "wooden_stick"},
		{Name: "two", Level: 50, WeaponSlot: "wooden_stick"},
	}
	current := allocation{
		equipment: []map[string]string{
			characterLoadout(party[0]),
			characterLoadout(party[1]),
		},
		utilities: []map[string]Utility{{}, {}},
	}
	available := map[string]int{
		"healing_aura_rune": 1,
		"mithril_ring":      1,
		"wooden_stick":      2,
	}
	seeds := roleSeeds(party, current, available)
	found := false
	for _, seed := range seeds {
		if seed.equipment[0]["ring1"] == "mithril_ring" && seed.equipment[1]["rune"] == "healing_aura_rune" {
			found = true
		}
		if seed.equipment[1]["ring1"] == "mithril_ring" && seed.equipment[0]["rune"] == "healing_aura_rune" {
			found = true
		}
	}
	if !found {
		t.Fatal("tank and support seed was not generated")
	}
}

func TestFindFightWithAvailableSkipsUtilitiesAfterEquipmentWin(t *testing.T) {
	party := []schemas.CharacterSchema{{
		Name:       "one",
		Level:      1,
		WeaponSlot: "wooden_stick",
	}}
	monster := schemas.MonsterSchema{Hp: 1, Initiative: 1}
	available := map[string]int{"health_boost_potion": 1, "wooden_stick": 1}
	result, err := FindFightWithAvailable(party, monster, available)
	if err != nil {
		t.Fatalf("find party fight: %v", err)
	}
	if len(result.Loadouts[0].Utilities) != 0 {
		t.Fatalf("utilities = %#v, want none", result.Loadouts[0].Utilities)
	}
}

func TestFindFightWithAvailableIsDeterministic(t *testing.T) {
	party := []schemas.CharacterSchema{
		{Name: "one", Level: 1, WeaponSlot: "wooden_stick"},
		{Name: "two", Level: 1, WeaponSlot: "wooden_stick"},
	}
	monster := schemas.MonsterSchema{Hp: 20, Initiative: 1}
	available := map[string]int{"wooden_stick": 2}
	first, err := FindFightWithAvailable(party, monster, available)
	if err != nil {
		t.Fatalf("find first party fight: %v", err)
	}
	second, err := FindFightWithAvailable(party, monster, available)
	if err != nil {
		t.Fatalf("find second party fight: %v", err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("results with the same inputs differ: %#v != %#v", first, second)
	}
}
