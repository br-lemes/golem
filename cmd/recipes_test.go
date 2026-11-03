package cmd

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestRecipesItemsSortsByXPAndFiltersSkillLevel(t *testing.T) {
	gear := "gearcrafting"
	level := 20
	quantity := 1
	items := []*schemas.ItemSchema{
		recipeItem("low_xp", gear, level, quantity, "copper_ore", 1),
		recipeItem("high_xp", gear, level, quantity, "copper_ore", 20),
		recipeItem("too_low", gear, 9, quantity, "copper_ore", 20),
	}

	got := recipesItems(gear, level, items, nil)
	codes := recipeCodes(got)
	if !reflect.DeepEqual(codes, []string{"high_xp", "low_xp"}) {
		t.Fatalf("recipe codes = %v, want [high_xp low_xp]", codes)
	}
}

func TestRecipesItemsExcludesTransitiveIngredients(t *testing.T) {
	gear := "gearcrafting"
	level := 20
	quantity := 1
	items := []*schemas.ItemSchema{
		recipeItem("copper_bar", gear, level, quantity, "copper_ore", 1),
		recipeItem("sword", gear, level, quantity, "copper_bar", 2),
	}

	got := recipesItems(gear, level, items, []string{"copper_ore"})
	if len(got) != 0 {
		t.Fatalf("recipes with excluded ingredient = %v, want none", recipeCodes(got))
	}
}

func TestRecipesItemsUsesTheTenLevelXPWindow(t *testing.T) {
	skill := "gearcrafting"
	items := []*schemas.ItemSchema{
		recipeItem("level_one", skill, 1, 1, "copper_ore", 1),
		recipeItem("level_two", skill, 2, 1, "copper_ore", 1),
		recipeItem("level_eleven", skill, 11, 1, "copper_ore", 1),
		recipeItem("level_twelve", skill, 12, 1, "copper_ore", 1),
	}

	got := recipeCodes(recipesItems(skill, 11, items, nil))
	want := []string{"level_eleven", "level_two", "level_one"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("recipes at level 11 = %v, want level 1 through 11", got)
	}
}

func TestRecipeDetailsUsesDirectIngredients(t *testing.T) {
	level := 20
	quantity := 1
	skill := "gearcrafting"
	sword := recipeItem("sword", skill, level, quantity, "", 0)
	sword.Craft.Items = &[]schemas.SimpleItemSchema{
		{Code: "copper_bar", Quantity: 2},
		{Code: "ash_wood", Quantity: 1},
	}

	items := []*schemas.ItemSchema{sword}

	want := []map[string]map[string]int{{
		"sword": {"copper_bar": 2, "ash_wood": 1},
	}}
	got := recipeDetails(items)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("recipe details = %#v, want %#v", got, want)
	}
}

func recipeItem(code, skill string, level, outputQuantity int, ingredientCode string, ingredientQuantity int) *schemas.ItemSchema {
	return &schemas.ItemSchema{
		Code: code,
		Craft: &schemas.CraftSchema{
			Items: &[]schemas.SimpleItemSchema{{
				Code:     ingredientCode,
				Quantity: ingredientQuantity,
			}},
			Level:    &level,
			Quantity: &outputQuantity,
			Skill:    &skill,
		},
	}
}

func recipeCodes(items []*schemas.ItemSchema) []string {
	codes := make([]string, 0, len(items))
	for _, item := range items {
		codes = append(codes, item.Code)
	}
	return codes
}
