package utils

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

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
