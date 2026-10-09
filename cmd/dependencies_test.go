package cmd

import (
	"testing"
)

func TestDependencyTreeIncludesOnlyRequestedPaths(t *testing.T) {
	tree := dependencyTree([]string{"copper_ore"}, nil, false, false, nil)
	ore := dependencyChild(t, tree, "copper_ore")
	dependencyNoQuantity(t, ore)
	if len(ore) != 0 {
		t.Fatalf("ore descendants = %#v, want none", ore)
	}
}

func TestDependencyTreeIncludesBankQuantities(t *testing.T) {
	tree := dependencyTree([]string{"copper_ore"}, map[string]int{
		"copper_ore": 42,
		"copper_bar": 8,
	}, false, false, nil)
	ore := dependencyChild(t, tree, "copper_ore")
	dependencyQuantity(t, ore, 42)
	if len(ore) != 1 {
		t.Fatalf("ore descendants = %#v, want only quantity", ore)
	}
}

func TestDependencyTreeResolvesToCurrencySource(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, nil, false, false, nil)
	coins := dependencyChild(t, tree, "sonnengott_coin")
	dependencyChild(t, coins, "powerful_rune")
	_, exists := coins["sonnengott_cloak"]
	if exists {
		t.Fatalf("unexpected sibling = %#v", coins["sonnengott_cloak"])
	}
}

func TestDependencyTreeExcludesNonItemCurrencies(t *testing.T) {
	tree := dependencyTree([]string{"cultist_boots"}, nil, false, false, nil)
	if tree["gold"] != nil {
		t.Fatalf("gold = %#v, want absent", tree["gold"])
	}
}

func TestDependencyTreeIncludesZeroBankQuantities(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, map[string]int{}, false, false, nil)
	coins := dependencyChild(t, tree, "sonnengott_coin")
	dependencyQuantity(t, coins, 0)
	rune := dependencyChild(t, coins, "powerful_rune")
	dependencyQuantity(t, rune, 0)
}

func TestDependencyTreeExpandsSourceUses(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, map[string]int{}, true, false, nil)
	coins := dependencyChild(t, tree, "sonnengott_coin")
	dependencyQuantity(t, coins, 0)
	rune := dependencyChild(t, coins, "powerful_rune")
	dependencyQuantity(t, rune, 0)
	cloak := dependencyChild(t, coins, "sonnengott_cloak")
	dependencyQuantity(t, cloak, 0)
}

func TestDependencyTreeExcludesExactZeroProducts(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, nil, true, false, map[string]bool{
		"powerful_rune": true,
	})
	coins := dependencyChild(t, tree, "sonnengott_coin")
	_, exists := coins["powerful_rune"]
	if exists {
		t.Fatalf("exact zero product = %#v", coins["powerful_rune"])
	}
	dependencyChild(t, coins, "sonnengott_cloak")
}

func TestDependencyPathsExcludeExactZeroProducts(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, nil, false, false, map[string]bool{
		"powerful_rune": true,
	})
	coins := dependencyChild(t, tree, "sonnengott_coin")
	if len(coins) != 0 {
		t.Fatalf("exact zero path = %#v", coins)
	}
}

func TestDependencyTreeResolvesToCraftingSource(t *testing.T) {
	tree := dependencyTree([]string{"copper_bar"}, nil, false, false, nil)
	ore := dependencyChild(t, tree, "copper_ore")
	dependencyChild(t, ore, "copper_bar")
}

func dependencyChild(t *testing.T, tree map[string]any, code string) map[string]any {
	t.Helper()
	next, ok := tree[code].(map[string]any)
	if !ok {
		t.Fatalf("%s = %#v, want map", code, tree[code])
	}
	return next
}

func dependencyQuantity(t *testing.T, tree map[string]any, want int) {
	t.Helper()
	got, ok := tree["_quantity"].(int)
	if !ok || got != want {
		t.Fatalf("_quantity = %#v, want %d", tree["_quantity"], want)
	}
}

func dependencyNoQuantity(t *testing.T, tree map[string]any) {
	t.Helper()
	_, exists := tree["_quantity"]
	if exists {
		t.Fatalf("_quantity = %#v, want absent", tree["_quantity"])
	}
}

func TestDependencyTreeIncludesItemLevels(t *testing.T) {
	tree := dependencyTree([]string{"copper_ore"}, nil, false, true, nil)
	dependencyLevel(t, dependencyChild(t, tree, "copper_ore"), 1)

	tree = dependencyTree([]string{"copper_ore"}, nil, true, true, nil)
	dependencyLevel(t, dependencyChild(t, tree, "copper_ore"), 1)
}

func dependencyLevel(t *testing.T, tree map[string]any, want int) {
	t.Helper()
	got, ok := tree["_level"].(int)
	if !ok || got != want {
		t.Fatalf("_level = %#v, want %d", tree["_level"], want)
	}
}
