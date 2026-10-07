package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCountInSlots(t *testing.T) {
	character := schemas.CharacterSchema{
		WeaponSlot:           "iron_sword",
		Utility1Slot:         "health_potion",
		Utility1SlotQuantity: 3,
	}

	quantities := countInSlots(character)
	if quantities["iron_sword"] != 1 || quantities["health_potion"] != 3 {
		t.Fatalf("countInSlots() = %#v", quantities)
	}
}
