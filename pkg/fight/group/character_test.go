package group

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFighterFromCharacterOmitsPlayerEffects(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:                 "test",
		MaxHp:                100,
		Threat:               25,
		Utility1Slot:         "health_potion",
		Utility1SlotQuantity: 10,
	}
	fighter := FighterFromCharacter(character)
	if fighter.Stats.Threat != 25 {
		t.Fatalf("threat = %d, want 25", fighter.Stats.Threat)
	}
	if len(fighter.Utilities) != 0 {
		t.Fatalf("utilities = %v, want none", fighter.Utilities)
	}
}
