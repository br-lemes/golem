package cmd

import (
	"slices"
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

func TestCountCodesWithoutArguments(t *testing.T) {
	bankItems := []schemas.SimpleItemSchema{{Code: "ash_wood", Quantity: 2}}
	inventoryItem := schemas.InventorySlotSchema{
		Code:     "iron_sword",
		Quantity: 1,
	}
	inventory := []schemas.InventorySlotSchema{inventoryItem}
	characters := []schemas.CharacterSchema{{Inventory: &inventory}}
	slots := []map[string]int{{"health_potion": 3}}

	codes := countCodes(nil, bankItems, characters, slots)
	want := []string{"ash_wood", "health_potion", "iron_sword"}
	if !slices.Equal(codes, want) {
		t.Fatalf("countCodes() = %#v, want %#v", codes, want)
	}
}

func TestCountValidateLevelRange(t *testing.T) {
	flags := countFlags{MinLevel: minGameLevel, MaxLevel: maxGameLevel}
	err := countValidate(nil, flags)
	if err != nil {
		t.Fatal(err)
	}
	flags.MinLevel = 0
	err = countValidate(nil, flags)
	if err == nil {
		t.Fatal("countValidate() accepted a level below the game range")
	}
	flags.MinLevel = minGameLevel
	flags.MaxLevel = maxGameLevel + 1
	err = countValidate(nil, flags)
	if err == nil {
		t.Fatal("countValidate() accepted a level above the game range")
	}
}

func TestCountItemMatches(t *testing.T) {
	recyclable := true
	skill := "gearcrafting"
	item := schemas.ItemSchema{
		Type:       "leg_armor",
		Subtype:    "",
		Level:      30,
		Tradeable:  true,
		Recyclable: &recyclable,
		Craft:      &schemas.CraftSchema{Skill: &skill},
	}
	flags := countFlags{
		Type:       []string{"leg_armor"},
		MinLevel:   20,
		MaxLevel:   40,
		Skill:      []string{"gearcrafting"},
		Tradeable:  true,
		Recyclable: true,
	}
	if !countItemMatches(item, flags) {
		t.Fatal("countItemMatches() = false")
	}
	flags.MaxLevel = 29
	if countItemMatches(item, flags) {
		t.Fatal("countItemMatches() = true above maximum level")
	}
}
