package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestDependencyTreeIncludesOnlyRequestedPaths(t *testing.T) {
	tree := dependencyTree([]string{"copper_ore"}, catalog.Items().All(), nil, nil, false)
	ore := dependencyChild(t, tree, "copper_ore")
	dependencyNoQuantity(t, ore)
	if len(ore) != 0 {
		t.Fatalf("ore descendants = %#v, want none", ore)
	}
}

func TestDependencyTreeIncludesBankQuantities(t *testing.T) {
	tree := dependencyTree([]string{"copper_ore"}, catalog.Items().All(), nil, map[string]int{
		"copper_ore": 42,
		"copper_bar": 8,
	}, false)
	ore := dependencyChild(t, tree, "copper_ore")
	dependencyQuantity(t, ore, 42)
	if len(ore) != 1 {
		t.Fatalf("ore descendants = %#v, want only quantity", ore)
	}
}

func TestDependencyTreeResolvesToCurrencySource(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, catalog.Items().All(), catalog.NpcsItems.All(), nil, false)
	coins := dependencyChild(t, tree, "sonnengott_coin")
	dependencyChild(t, coins, "powerful_rune")
	_, exists := coins["sonnengott_cloak"]
	if exists {
		t.Fatalf("unexpected sibling = %#v", coins["sonnengott_cloak"])
	}
}

func TestDependencyTreeIncludesZeroBankQuantities(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, catalog.Items().All(), catalog.NpcsItems.All(), map[string]int{}, false)
	coins := dependencyChild(t, tree, "sonnengott_coin")
	dependencyQuantity(t, coins, 0)
	rune := dependencyChild(t, coins, "powerful_rune")
	dependencyQuantity(t, rune, 0)
}

func TestDependenciesExpand(t *testing.T) {
	tests := []struct {
		mode     string
		hasCodes bool
		expand   bool
	}{
		{mode: "auto", expand: true},
		{mode: "auto", hasCodes: true},
		{mode: "expand", expand: true},
		{mode: "paths"},
	}
	for _, test := range tests {
		expand := dependenciesExpand(test.mode, test.hasCodes)
		if expand != test.expand {
			t.Fatalf("dependenciesExpand(%q, %t) = %t, want %t", test.mode, test.hasCodes, expand, test.expand)
		}
	}
	err := dependenciesValidate(dependenciesFlags{Mode: "invalid"})
	if err == nil {
		t.Fatal("dependenciesValidate(invalid) error = nil, want error")
	}
}

func TestDependencyTreeExpandsSourceUses(t *testing.T) {
	tree := dependencyTree([]string{"powerful_rune"}, catalog.Items().All(), catalog.NpcsItems.All(), map[string]int{}, true)
	coins := dependencyChild(t, tree, "sonnengott_coin")
	dependencyQuantity(t, coins, 0)
	rune := dependencyChild(t, coins, "powerful_rune")
	dependencyQuantity(t, rune, 0)
	cloak := dependencyChild(t, coins, "sonnengott_cloak")
	dependencyQuantity(t, cloak, 0)
}

func TestDependencyTreeResolvesToCraftingSource(t *testing.T) {
	tree := dependencyTree([]string{"copper_bar"}, catalog.Items().All(), nil, nil, false)
	ore := dependencyChild(t, tree, "copper_ore")
	dependencyChild(t, ore, "copper_bar")
}

func TestDependencyTreeStopsAtCycles(t *testing.T) {
	recipe := func(code, ingredient string) *schemas.ItemSchema {
		items := []schemas.SimpleItemSchema{{Code: ingredient, Quantity: 1}}
		return &schemas.ItemSchema{
			Code:  code,
			Craft: &schemas.CraftSchema{Items: &items},
		}
	}
	tree := dependencyTree([]string{"a"}, []*schemas.ItemSchema{
		recipe("a", "b"),
		recipe("b", "a"),
	}, nil, nil, false)
	a := dependencyChild(t, tree, "a")
	b := dependencyChild(t, a, "b")
	if len(b) != 0 {
		t.Fatalf("cycle descendants = %#v, want none", b)
	}
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
