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

func TestEvaluateItemAccountsForInventoryPenaltyAgainstWisdom(t *testing.T) {
	ctx := &bestCtx{
		Weights: map[string]int{"wisdom": 10000, "inventory_space": 1},
	}

	penalized := schemas.ItemSchema{
		Effects: &[]schemas.SimpleEffectSchema{
			{Code: "wisdom", Value: 25},
			{Code: "inventory_space", Value: -5},
		},
	}
	plain := schemas.ItemSchema{
		Effects: &[]schemas.SimpleEffectSchema{{Code: "wisdom", Value: 24}},
	}

	got := ctx.evaluateItem(penalized)
	want := 25*10000 - 5*10000 - 5
	if got != want {
		t.Fatalf("penalized item score = %d, want %d", got, want)
	}
	if ctx.evaluateItem(penalized) >= ctx.evaluateItem(plain) {
		t.Fatal("small wisdom gain should not outweigh inventory loss")
	}
}
