package cmd

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCheapestValidate(t *testing.T) {
	tests := []struct {
		name    string
		skill   string
		wantErr bool
	}{
		{name: "gearcrafting", skill: "gearcrafting"},
		{name: "jewelrycrafting", skill: "jewelrycrafting"},
		{name: "weaponcrafting", skill: "weaponcrafting"},
		{name: "invalid", skill: "mining", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cheapestValidate(tt.skill)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("cheapestValidate() error = %v, want error = %v", err, tt.wantErr)
			}
		})
	}
}

func TestCraftCostLess(t *testing.T) {
	base := craftCost{
		basic:        1,
		drops:        1,
		rarity:       1,
		monsterLevel: 1,
		steps:        1,
		units:        1,
	}
	basicA := craftCost{basic: 1}
	basicB := craftCost{basic: 2}
	dropsA := craftCost{basic: 1, drops: 1}
	dropsB := craftCost{basic: 1, drops: 2}
	rarityA := craftCost{basic: 1, drops: 1, rarity: 1}
	rarityB := craftCost{basic: 1, drops: 1, rarity: 2}
	monsterA := craftCost{basic: 1, drops: 1, rarity: 1, monsterLevel: 1}
	monsterB := craftCost{basic: 1, drops: 1, rarity: 1, monsterLevel: 2}
	stepsA := craftCost{
		basic:        1,
		drops:        1,
		rarity:       1,
		monsterLevel: 1,
		steps:        1,
	}
	stepsB := craftCost{
		basic:        1,
		drops:        1,
		rarity:       1,
		monsterLevel: 1,
		steps:        2,
	}
	unitsB := craftCost{
		basic:        1,
		drops:        1,
		rarity:       1,
		monsterLevel: 1,
		steps:        1,
		units:        2,
	}
	checkCraftCostLess(t, "basic", basicA, basicB, true)
	checkCraftCostLess(t, "drops", dropsA, dropsB, true)
	checkCraftCostLess(t, "rarity", rarityA, rarityB, true)
	checkCraftCostLess(t, "monster level", monsterA, monsterB, true)
	checkCraftCostLess(t, "steps", stepsA, stepsB, true)
	checkCraftCostLess(t, "units", base, unitsB, true)
	checkCraftCostLess(t, "equal", base, base, false)
}

func checkCraftCostLess(t *testing.T, name string, a, b craftCost, want bool) {
	t.Helper()
	got := craftCostLess(a, b)
	if got != want {
		t.Errorf("%s: craftCostLess() = %v, want %v", name, got, want)
	}
}

func TestDependsOnTasksCoin(t *testing.T) {
	if !dependsOnTasksCoin("tasks_coin", map[string]bool{}) {
		t.Error("tasks_coin should be detected as a dependency")
	}
	if dependsOnTasksCoin("unknown_test_item", map[string]bool{}) {
		t.Error("unknown item should not be detected as a dependency")
	}
}

func TestItemCraftCost(t *testing.T) {
	got := itemCraftCost("unknown_test_item", map[string]bool{})
	if got.steps != 0 {
		t.Errorf("itemCraftCost() for database-missing item steps = %d, want 0", got.steps)
	}
}

func TestCheapestItemsFiltersSkillAndRange(t *testing.T) {
	gear := schemas.CraftSkill("gearcrafting")
	weapon := schemas.CraftSkill("weaponcrafting")
	levelFive := 5
	levelFifteen := 15
	items := []*schemas.ItemSchema{
		{
			Code:  "outside_low",
			Craft: &schemas.CraftSchema{Level: &levelFive, Skill: &gear},
		},
		{
			Code:  "inside",
			Craft: &schemas.CraftSchema{Level: &levelFifteen, Skill: &gear},
		},
		{
			Code:  "wrong_skill",
			Craft: &schemas.CraftSchema{Level: &levelFifteen, Skill: &weapon},
		},
	}
	got := cheapestItems("gearcrafting", 15, items)
	want := []string{"inside", "outside_low"}
	codes := make([]string, 0, len(got))
	for _, item := range got {
		codes = append(codes, item.Code)
	}
	if !reflect.DeepEqual(codes, want) {
		t.Errorf("cheapestItems() codes = %v, want %v", codes, want)
	}
}
