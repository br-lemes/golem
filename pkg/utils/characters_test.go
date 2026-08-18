package utils

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestGetCharacterSkillLevel(t *testing.T) {
	character := schemas.CharacterSchema{
		Level:            12,
		AlchemyLevel:     3,
		MiningLevel:      7,
		WoodcuttingLevel: 9,
	}

	tests := []struct {
		skill string
		level int
		ok    bool
	}{
		{skill: "alchemy", level: 3, ok: true},
		{skill: "combat", level: 12, ok: true},
		{skill: "mining", level: 7, ok: true},
		{skill: "woodcutting", level: 9, ok: true},
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
	character := schemas.CharacterSchema{Level: 12, MiningLevel: 8}
	tests := []struct {
		code  string
		level int
		ok    bool
	}{
		{code: "level", level: 12, ok: true},
		{code: "mining_level", level: 8, ok: true},
		{code: "unknown", level: 0, ok: false},
	}
	for _, test := range tests {
		level, ok := GetCharacterConditionLevel(character, test.code)
		if level != test.level || ok != test.ok {
			t.Errorf("GetCharacterConditionLevel(%q) = (%d, %t), want (%d, %t)", test.code, level, ok, test.level, test.ok)
		}
	}
}
