package best

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFindEquipmentUsesOwnedItems(t *testing.T) {
	character := schemas.CharacterSchema{Level: 50}
	got, err := FindEquipment(character, EquipmentOptions{
		Owned:      map[string]int{"iron_sword": 1},
		Priorities: []string{"attack_earth"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got["weapon"].Code != "iron_sword" {
		t.Fatalf("weapon = %#v, want iron_sword", got["weapon"])
	}
}

func TestFindEquipmentRejectsInvalidPriorities(t *testing.T) {
	_, err := FindEquipment(schemas.CharacterSchema{}, EquipmentOptions{
		Owned:      map[string]int{},
		Priorities: []string{"unknown"},
	})
	if err == nil {
		t.Fatal("invalid priorities were accepted")
	}
}

func TestFindEquipmentUsesGatheringSkillForTools(t *testing.T) {
	character := schemas.CharacterSchema{Level: 50, MiningLevel: 9}
	got, err := FindEquipment(character, EquipmentOptions{
		Owned:      map[string]int{"iron_pickaxe": 1},
		Priorities: []string{"mining"},
	})
	if err != nil {
		t.Fatal(err)
	}
	weapon := got["weapon"]
	if weapon.Code != "" {
		t.Fatalf("weapon = %#v, want empty result", weapon)
	}
}

func TestFindEquipmentSchemas(t *testing.T) {
	character := schemas.CharacterSchema{Level: 50}
	got, err := FindEquipmentSchemas(character, EquipmentOptions{
		Owned:      map[string]int{"iron_sword": 1},
		Priorities: []string{"attack_earth"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []schemas.EquipSchema{{Code: "iron_sword", Slot: "weapon"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("equipment = %#v, want %#v", got, want)
	}
}

func TestFindEquipmentSchemasSortsBySlot(t *testing.T) {
	character := schemas.CharacterSchema{Level: 50}
	got, err := FindEquipmentSchemas(character, EquipmentOptions{
		Owned:      map[string]int{"iron_sword": 1, "copper_boots": 1},
		Priorities: []string{"attack_earth", "wisdom"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []schemas.EquipSchema{
		{Code: "copper_boots", Slot: "boots"},
		{Code: "iron_sword", Slot: "weapon"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("equipment = %#v, want %#v", got, want)
	}
}

func TestFindEquipmentSchemasReturnsPriorityError(t *testing.T) {
	_, err := FindEquipmentSchemas(schemas.CharacterSchema{}, EquipmentOptions{
		Owned:      map[string]int{},
		Priorities: []string{"unknown"},
	})
	if err == nil {
		t.Fatal("invalid priorities were accepted")
	}
}

func TestFetchItemsCountsEquippedUtilities(t *testing.T) {
	character := schemas.CharacterSchema{
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: 7,
		Utility2Slot:         "small_mana_potion",
		Utility2SlotQuantity: 8,
	}
	ctx := bestCtx{Character: character}
	err := ctx.fetchItems(map[string]int{})
	if err != nil {
		t.Fatal(err)
	}
	if ctx.AlreadyEquipped["small_health_potion"] != 7 {
		t.Fatalf("utility1 quantity = %d, want 7", ctx.AlreadyEquipped["small_health_potion"])
	}
	if ctx.AlreadyEquipped["small_mana_potion"] != 8 {
		t.Fatalf("utility2 quantity = %d, want 8", ctx.AlreadyEquipped["small_mana_potion"])
	}
}

func TestEvaluateItem(t *testing.T) {
	ctx := bestCtx{
		Skill:   "mining",
		Weights: map[string]int{"wisdom": 100, "mining": 1},
	}
	item := schemas.ItemSchema{
		Effects: &[]schemas.SimpleEffectSchema{
			{Code: "wisdom", Value: 2},
			{Code: "mining", Value: 3},
			{Code: "unknown", Value: 100},
		},
	}
	got := ctx.evaluateItem(item)
	if got != 197 {
		t.Fatalf("evaluateItem() = %d, want 197", got)
	}
	got = ctx.evaluateItem(schemas.ItemSchema{})
	if got != 0 {
		t.Fatalf("evaluateItem() without effects = %d, want 0", got)
	}
}

func TestFormatItem(t *testing.T) {
	ctx := bestCtx{Weights: map[string]int{"wisdom": 1, "attack": 1}}
	item := schemas.ItemSchema{
		Effects: &[]schemas.SimpleEffectSchema{
			{Code: "wisdom", Value: 2},
			{Code: "attack", Value: -3},
			{Code: "unknown", Value: 4},
		},
	}
	got := ctx.formatItem(item)
	if got != "+2 wisdom, -3 attack" {
		t.Fatalf("formatItem() = %q, want %q", got, "+2 wisdom, -3 attack")
	}
	got = ctx.formatItem(schemas.ItemSchema{})
	if got != "" {
		t.Fatalf("formatItem() without effects = %q, want empty", got)
	}
}

func TestBestGroupCanMoveAnArtifactBetweenSlots(t *testing.T) {
	ctx := &bestCtx{
		ItemValues: map[string]int{"old": 10, "best": 30, "other": 20},
		OwnedItems: map[string]int{"old": 1, "best": 1, "other": 1},
		ValidItems: []schemas.ItemSchema{
			{Code: "best", Type: "artifact"},
			{Code: "other", Type: "artifact"},
			{Code: "old", Type: "artifact"},
		},
	}

	got := ctx.bestGroup([]string{"artifact1", "artifact2", "artifact3"}, "artifact")
	if got["artifact1"] != "best" || got["artifact2"] != "other" || got["artifact3"] != "old" {
		t.Fatalf("best artifact assignment = %#v", got)
	}
}

func TestBestGroupDoesNotReuseUniqueArtifacts(t *testing.T) {
	ctx := &bestCtx{
		ItemValues: map[string]int{"best": 30},
		OwnedItems: map[string]int{"best": 3},
		ValidItems: []schemas.ItemSchema{{Code: "best", Type: "artifact"}},
	}

	got := ctx.bestGroup([]string{"artifact1", "artifact2", "artifact3"}, "artifact")
	count := 0
	for _, code := range got {
		if code == "best" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("unique artifact used %d times in %#v", count, got)
	}
}

func TestBestGroupDoesNotReuseUniqueAdeptRing(t *testing.T) {
	ctx := &bestCtx{
		ItemValues:      map[string]int{"ring_of_the_adept": 30},
		OwnedItems:      map[string]int{"ring_of_the_adept": 2},
		UniqueAdeptRing: true,
		ValidItems: []schemas.ItemSchema{
			{Code: "ring_of_the_adept", Type: "ring"},
		},
	}

	got := ctx.bestGroup([]string{"ring1", "ring2"}, "ring")
	count := 0
	for _, code := range got {
		if code == "ring_of_the_adept" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("unique adept ring used %d times in %#v", count, got)
	}
}

func TestBestGroupDoesNotReuseUniqueUtilities(t *testing.T) {
	ctx := &bestCtx{
		ItemValues: map[string]int{"best": 30},
		OwnedItems: map[string]int{"best": 2},
		ValidItems: []schemas.ItemSchema{{Code: "best", Type: "utility"}},
	}

	got := ctx.bestGroup([]string{"utility1", "utility2"}, "utility")
	count := 0
	for _, code := range got {
		if code == "best" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("unique utility used %d times in %#v", count, got)
	}
}

func TestBestGroupKeepsEquippedAdeptRing(t *testing.T) {
	ctx := &bestCtx{
		ItemValues: map[string]int{
			"ring_of_the_adept": 30,
			"ring_of_chance":    30,
		},
		OwnedItems: map[string]int{"ring_of_the_adept": 1, "ring_of_chance": 1},
		Equipped: map[string]string{
			"ring1": "ring_of_chance",
			"ring2": "ring_of_the_adept",
		},
		AlreadyEquipped: map[string]int{
			"ring_of_chance":    1,
			"ring_of_the_adept": 1,
		},
		UniqueAdeptRing: true,
		ValidItems: []schemas.ItemSchema{
			{Code: "ring_of_the_adept", Type: "ring"},
			{Code: "ring_of_chance", Type: "ring"},
		},
	}

	got := ctx.bestGroup([]string{"ring1", "ring2"}, "ring")
	if got["ring1"] != "ring_of_chance" || got["ring2"] != "ring_of_the_adept" {
		t.Fatalf("equipped rings were changed: %#v", got)
	}
}

func TestHasNegativeInventorySpace(t *testing.T) {
	if hasNegativeInventorySpace(schemas.ItemSchema{}) {
		t.Fatal("item without effects was rejected")
	}
	item := schemas.ItemSchema{
		Effects: &[]schemas.SimpleEffectSchema{
			{Code: "wisdom", Value: 25},
			{Code: "inventory_space", Value: -5},
		},
	}
	if !hasNegativeInventorySpace(item) {
		t.Fatal("item with negative inventory space was not rejected")
	}
	item.Effects = &[]schemas.SimpleEffectSchema{
		{Code: "inventory_space", Value: 0},
	}
	if hasNegativeInventorySpace(item) {
		t.Fatal("item without negative inventory space was rejected")
	}
}

func TestFilterAndSortExcludesNegativeInventorySpaceItems(t *testing.T) {
	ctx := &bestCtx{
		Character: schemas.CharacterSchema{Level: 50},
		Weights: map[string]int{
			"attack_earth":    1,
			"critical_strike": 1,
			"inventory_space": 1,
		},
	}

	ctx.filterAndSort()
	for _, item := range ctx.ValidItems {
		if item.Code == "obsidian_battleaxe" {
			t.Fatal("obsidian battleaxe was included despite negative inventory space")
		}
	}
}

func TestBestGroupRecommendsAlternativeForEquippedNegativeInventoryItem(t *testing.T) {
	ctx := &bestCtx{
		Character: schemas.CharacterSchema{Level: 50},
		Weights: map[string]int{
			"attack_earth":    1,
			"critical_strike": 1,
			"inventory_space": 1,
		},
	}
	ctx.filterAndSort()
	ctx.OwnedItems = map[string]int{"obsidian_battleaxe": 1, "wooden_stick": 1}
	ctx.Equipped = map[string]string{"weapon": "obsidian_battleaxe"}
	ctx.AlreadyEquipped = map[string]int{"obsidian_battleaxe": 1}
	got := ctx.bestGroup([]string{"weapon"}, "weapon")
	if got["weapon"] != "wooden_stick" {
		t.Fatalf("equipped negative inventory item was not replaced: %#v", got)
	}
}

func TestMatchEquipmentSkipsEmptyRecommendation(t *testing.T) {
	ctx := &bestCtx{
		OwnedItems: map[string]int{},
		ValidItems: []schemas.ItemSchema{{Code: "iron_sword", Type: "weapon"}},
		ItemValues: map[string]int{"iron_sword": 1},
		Equipped:   map[string]string{},
	}

	ctx.matchEquipment()
	if len(ctx.Result) != 0 {
		t.Fatalf("result = %#v, want empty", ctx.Result)
	}
}

func TestMatchEquipmentSkipsAlreadyEquippedRecommendation(t *testing.T) {
	ctx := &bestCtx{
		OwnedItems: map[string]int{"iron_sword": 1},
		ValidItems: []schemas.ItemSchema{{
			Code: "iron_sword",
			Type: "weapon",
			Effects: &[]schemas.SimpleEffectSchema{
				{Code: "attack_earth", Value: 1},
			},
		}},
		ItemValues: map[string]int{"iron_sword": 1},
		Equipped:   map[string]string{"weapon": "iron_sword"},
		Weights:    map[string]int{"attack_earth": 1},
	}

	ctx.matchEquipment()
	if len(ctx.Result) != 0 {
		t.Fatalf("result = %#v, want empty", ctx.Result)
	}
}
