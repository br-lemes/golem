package cmd

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestXPNeeded(t *testing.T) {
	tests := []struct {
		name                 string
		level, current, goal int
		want                 int
	}{
		{"from zero to level 2", 1, 0, 2, 150},
		{"from zero to level 5", 1, 0, 5, 1200},
		{"subtract current progress", 5, 300, 6, 400},
		{"already reached", 10, 0, 10, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xpNeeded(tt.level, tt.current, tt.goal)
			if got != tt.want {
				t.Errorf("xpNeeded() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestXPValidate(t *testing.T) {
	valid := xpFlags{Group: "skill", Skill: []string{"combat"}, From: -1}
	tests := []struct {
		name    string
		target  string
		options xpFlags
		wantErr bool
	}{
		{
			name:    "valid minimum",
			target:  "1",
			options: valid,
		},
		{
			name:    "valid maximum",
			target:  "50",
			options: valid,
		},
		{
			name:    "invalid target",
			target:  "0",
			options: valid,
			wantErr: true,
		},
		{
			name:    "invalid target text",
			target:  "thirty",
			options: valid,
			wantErr: true,
		},
		{
			name:    "invalid group",
			target:  "30",
			options: xpFlags{Group: "other"},
			wantErr: true,
		},
		{
			name:    "invalid skill",
			target:  "30",
			options: xpFlags{Group: "skill", Skill: []string{"other"}},
			wantErr: true,
		},
		{
			name:    "from with name",
			target:  "30",
			options: xpFlags{From: 10, Name: []string{"Ada"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := xpValidate(tt.target, tt.options, false)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("xpValidate() error = %v, want error = %v", err, tt.wantErr)
			}
		})
	}
}

func TestCharacterSkillProgress(t *testing.T) {
	character := schemas.CharacterSchema{
		Level:                10,
		Xp:                   11,
		AlchemyLevel:         2,
		AlchemyXp:            12,
		CookingLevel:         3,
		CookingXp:            13,
		FishingLevel:         4,
		FishingXp:            14,
		GearcraftingLevel:    5,
		GearcraftingXp:       15,
		JewelrycraftingLevel: 6,
		JewelrycraftingXp:    16,
		MiningLevel:          7,
		MiningXp:             17,
		WeaponcraftingLevel:  8,
		WeaponcraftingXp:     18,
		WoodcuttingLevel:     9,
		WoodcuttingXp:        19,
	}
	want := map[string][2]int{
		"combat":          {10, 11},
		"alchemy":         {2, 12},
		"cooking":         {3, 13},
		"fishing":         {4, 14},
		"gearcrafting":    {5, 15},
		"jewelrycrafting": {6, 16},
		"mining":          {7, 17},
		"weaponcrafting":  {8, 18},
		"woodcutting":     {9, 19},
	}
	for skill, expected := range want {
		level, xp := characterSkillProgress(character, skill)
		if level != expected[0] || xp != expected[1] {
			t.Errorf("characterSkillProgress(%q) = (%d, %d), want (%d, %d)", skill, level, xp, expected[0], expected[1])
		}
	}
}

func TestXPByCharacter(t *testing.T) {
	input := map[string]map[string]int{
		"combat": {"Ada": 100, "Bob": 200},
		"mining": {"Ada": 300},
	}
	want := map[string]map[string]int{
		"Ada": {"combat": 100, "mining": 300},
		"Bob": {"combat": 200},
	}
	got := xpByCharacter(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("xpByCharacter() = %#v, want %#v", got, want)
	}
}
