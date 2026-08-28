package utils

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestGetCharacterSkillLevel(t *testing.T) {
	character := schemas.CharacterSchema{
		Level:                12,
		AlchemyLevel:         3,
		CookingLevel:         4,
		FishingLevel:         5,
		GearcraftingLevel:    6,
		JewelrycraftingLevel: 7,
		MiningLevel:          8,
		WeaponcraftingLevel:  9,
		WoodcuttingLevel:     10,
	}

	tests := []struct {
		skill string
		level int
		ok    bool
	}{
		{skill: "alchemy", level: 3, ok: true},
		{skill: "combat", level: 12, ok: true},
		{skill: "cooking", level: 4, ok: true},
		{skill: "fishing", level: 5, ok: true},
		{skill: "gearcrafting", level: 6, ok: true},
		{skill: "jewelrycrafting", level: 7, ok: true},
		{skill: "mining", level: 8, ok: true},
		{skill: "weaponcrafting", level: 9, ok: true},
		{skill: "woodcutting", level: 10, ok: true},
		{skill: "unknown", level: 0, ok: false},
	}

	for _, test := range tests {
		level, ok := GetCharacterSkillLevel(character, test.skill)
		if level != test.level || ok != test.ok {
			t.Errorf("GetCharacterSkillLevel(%q) = (%d, %t), want (%d, %t)", test.skill, level, ok, test.level, test.ok)
		}
	}
}

func TestGetCharacterConditionLevel(t *testing.T) {
	character := schemas.CharacterSchema{
		Level:                12,
		AlchemyLevel:         3,
		CookingLevel:         4,
		FishingLevel:         5,
		GearcraftingLevel:    6,
		JewelrycraftingLevel: 7,
		MiningLevel:          8,
		WeaponcraftingLevel:  9,
		WoodcuttingLevel:     10,
	}
	tests := []struct {
		code  string
		level int
		ok    bool
	}{
		{code: "level", level: 12, ok: true},
		{code: "alchemy_level", level: 3, ok: true},
		{code: "cooking_level", level: 4, ok: true},
		{code: "fishing_level", level: 5, ok: true},
		{code: "gearcrafting_level", level: 6, ok: true},
		{code: "jewelrycrafting_level", level: 7, ok: true},
		{code: "mining_level", level: 8, ok: true},
		{code: "weaponcrafting_level", level: 9, ok: true},
		{code: "woodcutting_level", level: 10, ok: true},
		{code: "unknown", level: 0, ok: false},
	}
	for _, test := range tests {
		level, ok := GetCharacterConditionLevel(character, test.code)
		if level != test.level || ok != test.ok {
			t.Errorf("GetCharacterConditionLevel(%q) = (%d, %t), want (%d, %t)", test.code, level, ok, test.level, test.ok)
		}
	}
}

func TestGetCharacterGatheringSkillLevel(t *testing.T) {
	character := schemas.CharacterSchema{FishingLevel: 5}
	level, ok := GetCharacterGatheringSkillLevel(character, "fishing")
	if level != 5 || !ok {
		t.Fatalf("GetCharacterGatheringSkillLevel() = (%d, %t), want (5, true)", level, ok)
	}
	level, ok = GetCharacterGatheringSkillLevel(character, "cooking")
	if level != 0 || ok {
		t.Fatalf("GetCharacterGatheringSkillLevel() = (%d, %t), want (0, false)", level, ok)
	}
}

func TestGetCharacterCraftingSkillLevel(t *testing.T) {
	character := schemas.CharacterSchema{CookingLevel: 4}
	level, ok := GetCharacterCraftingSkillLevel(character, "cooking")
	if level != 4 || !ok {
		t.Fatalf("GetCharacterCraftingSkillLevel() = (%d, %t), want (4, true)", level, ok)
	}
	level, ok = GetCharacterCraftingSkillLevel(character, "fishing")
	if level != 0 || ok {
		t.Fatalf("GetCharacterCraftingSkillLevel() = (%d, %t), want (0, false)", level, ok)
	}
}

func TestMeetsItemConditions(t *testing.T) {
	character := schemas.CharacterSchema{Level: 10, MiningLevel: 5}
	tests := []struct {
		name      string
		condition schemas.ConditionSchema
		want      bool
	}{
		{
			name: "gt satisfied",
			condition: schemas.ConditionSchema{
				Code:     "level",
				Operator: schemas.Gt,
				Value:    9,
			},
			want: true,
		},
		{
			name: "gt unsatisfied",
			condition: schemas.ConditionSchema{
				Code:     "level",
				Operator: schemas.Gt,
				Value:    10,
			},
		},
		{
			name: "eq satisfied",
			condition: schemas.ConditionSchema{
				Code:     "mining_level",
				Operator: schemas.Eq,
				Value:    5,
			},
			want: true,
		},
		{
			name: "eq unsatisfied",
			condition: schemas.ConditionSchema{
				Code:     "mining_level",
				Operator: schemas.Eq,
				Value:    4,
			},
		},
		{
			name: "lt satisfied",
			condition: schemas.ConditionSchema{
				Code:     "level",
				Operator: schemas.Lt,
				Value:    11,
			},
			want: true,
		},
		{
			name: "lt unsatisfied",
			condition: schemas.ConditionSchema{
				Code:     "level",
				Operator: schemas.Lt,
				Value:    10,
			},
		},
		{
			name: "ne satisfied",
			condition: schemas.ConditionSchema{
				Code:     "level",
				Operator: schemas.Ne,
				Value:    9,
			},
			want: true,
		},
		{
			name: "ne unsatisfied",
			condition: schemas.ConditionSchema{
				Code:     "level",
				Operator: schemas.Ne,
				Value:    10,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conditions := []schemas.ConditionSchema{test.condition}
			item := schemas.ItemSchema{Conditions: &conditions}
			got := MeetsItemConditions(character, item)
			if got != test.want {
				t.Fatalf("MeetsItemConditions() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestMeetsItemConditionsWithoutConditions(t *testing.T) {
	if !MeetsItemConditions(schemas.CharacterSchema{}, schemas.ItemSchema{}) {
		t.Fatal("item without conditions was rejected")
	}
}
