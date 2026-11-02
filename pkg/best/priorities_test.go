package best

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestGatheringPriorities(t *testing.T) {
	tests := []struct {
		name       string
		level      int
		resource   schemas.ResourceSchema
		priorities []string
	}{
		{
			name:       "near resource level",
			level:      5,
			resource:   schemas.ResourceSchema{Level: 10, Skill: "mining"},
			priorities: []string{"mining_gathering", "wisdom", "prospecting"},
		},
		{
			name:       "far above resource level",
			level:      25,
			resource:   schemas.ResourceSchema{Level: 10, Skill: "mining"},
			priorities: []string{"mining_gathering", "prospecting"},
		},
		{
			name:       "maximum level near resource level",
			level:      50,
			resource:   schemas.ResourceSchema{Level: 45, Skill: "mining"},
			priorities: []string{"mining_gathering", "prospecting"},
		},
		{
			name: "single drop resource with XP",
			resource: schemas.ResourceSchema{
				Level: 10,
				Skill: "mining",
				Drops: []schemas.DropRateSchema{{Code: "iron_ore"}},
			},
			priorities: []string{"mining_gathering", "wisdom"},
		},
		{
			name:  "single drop resource without XP",
			level: 25,
			resource: schemas.ResourceSchema{
				Level: 10,
				Skill: "mining",
				Drops: []schemas.DropRateSchema{{Code: "iron_ore"}},
			},
			priorities: []string{"mining_gathering"},
		},
		{
			name:  "event near resource level",
			level: 5,
			resource: schemas.ResourceSchema{
				Code:  "magic_tree",
				Level: 10,
				Skill: "mining",
			},
			priorities: []string{"mining_gathering", "prospecting", "wisdom"},
		},
		{
			name:  "event far above resource level",
			level: 25,
			resource: schemas.ResourceSchema{
				Code:  "magic_tree",
				Level: 10,
				Skill: "mining",
			},
			priorities: []string{"mining_gathering", "prospecting"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			character := schemas.CharacterSchema{MiningLevel: test.level}
			got := GatheringPriorities(character, &test.resource)
			if !reflect.DeepEqual(got, test.priorities) {
				t.Fatalf("priorities = %#v, want %#v", got, test.priorities)
			}
		})
	}
}

func TestCraftingPriorities(t *testing.T) {
	skill := "mining"
	tests := []struct {
		name       string
		level      int
		itemLevel  int
		priorities []string
	}{
		{
			name:       "near item level",
			level:      5,
			itemLevel:  10,
			priorities: []string{"mining_crafting", "wisdom"},
		},
		{
			name:       "far above item level",
			level:      25,
			itemLevel:  10,
			priorities: []string{"mining_crafting"},
		},
		{
			name:       "maximum level near item level",
			level:      50,
			itemLevel:  45,
			priorities: []string{"mining_crafting"},
		},
		{
			name:       "missing skill level",
			level:      0,
			itemLevel:  1,
			priorities: []string{"mining_crafting"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			item := schemas.ItemSchema{
				Level: test.itemLevel,
				Craft: &schemas.CraftSchema{Skill: &skill},
			}
			character := schemas.CharacterSchema{MiningLevel: test.level}
			got := CraftingPriorities(character, &item)
			if !reflect.DeepEqual(got, test.priorities) {
				t.Fatalf("priorities = %#v, want %#v", got, test.priorities)
			}
		})
	}
}

func TestNormalizePriorities(t *testing.T) {
	got, err := NormalizePriorities([]string{
		"wisdom",
		"mining_gathering",
		"prospecting",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"mining_gathering",
		"wisdom",
		"prospecting",
		"inventory_space",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %#v, want %#v", got, want)
	}
}

func TestNormalizePrioritiesDefaultsToInventorySpace(t *testing.T) {
	got, err := NormalizePriorities(nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"inventory_space"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %#v, want %#v", got, want)
	}
}

func TestNormalizePrioritiesPreservesExplicitInventoryPosition(t *testing.T) {
	priorities := []string{"wisdom", "inventory_space", "prospecting"}
	got, err := NormalizePriorities(priorities)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"wisdom", "inventory_space", "prospecting"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %#v, want %#v", got, want)
	}
}

func TestNormalizePrioritiesRejectsInvalidEffect(t *testing.T) {
	_, err := NormalizePriorities([]string{"wisdom", "unknown"})
	if err == nil {
		t.Fatal("expected invalid effect error")
	}
}

func TestNormalizePrioritiesRejectsDuplicateEffect(t *testing.T) {
	_, err := NormalizePriorities([]string{"wisdom", "wisdom"})
	if err == nil {
		t.Fatal("expected duplicate effect error")
	}
}

func TestNormalizePrioritiesRejectsMultipleGatheringSkills(t *testing.T) {
	_, err := NormalizePriorities([]string{
		"mining_gathering",
		"woodcutting_gathering",
	})
	if err == nil {
		t.Fatal("expected multiple gathering skills error")
	}
}

func TestNormalizePrioritiesRejectsMultipleCraftingSkills(t *testing.T) {
	_, err := NormalizePriorities([]string{
		"mining_crafting",
		"woodcutting_crafting",
	})
	if err == nil {
		t.Fatal("expected multiple crafting skills error")
	}
}

func TestNormalizePrioritiesRejectsGatheringAndCrafting(t *testing.T) {
	_, err := NormalizePriorities([]string{
		"mining_gathering",
		"mining_crafting",
	})
	if err == nil {
		t.Fatal("expected gathering and crafting conflict error")
	}
}
