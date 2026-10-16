package group

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestAvailableItemsIncludesOnlyPartyHoldings(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "iron_sword", Quantity: 3},
		{Code: "health_potion", Quantity: 4},
	}
	party := []schemas.CharacterSchema{{
		Inventory:            &inventory,
		WeaponSlot:           "wooden_stick",
		Ring1Slot:            "copper_ring",
		Ring2Slot:            "copper_ring",
		Utility1Slot:         "health_potion",
		Utility1SlotQuantity: 5,
		Utility2Slot:         "health_potion",
		Utility2SlotQuantity: 2,
	}}
	bank := []schemas.SimpleItemSchema{
		{Code: "wooden_stick", Quantity: 1},
		{Code: "health_potion", Quantity: 2},
	}

	got := AvailableItems(bank, party)
	want := map[string]int{
		"copper_ring":   2,
		"health_potion": 13,
		"iron_sword":    3,
		"wooden_stick":  2,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("available items = %#v, want %#v", got, want)
	}
}

func TestAvailableItemsOmitsEmptyAndNonpositiveQuantities(t *testing.T) {
	bank := []schemas.SimpleItemSchema{
		{Code: "", Quantity: 2},
		{Code: "ignored", Quantity: 0},
	}
	party := []schemas.CharacterSchema{{
		Utility1Slot:         "ignored",
		Utility1SlotQuantity: 0,
	}}

	got := AvailableItems(bank, party)
	if len(got) != 0 {
		t.Fatalf("available items = %#v, want empty", got)
	}
}
