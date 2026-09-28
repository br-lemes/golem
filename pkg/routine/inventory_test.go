package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestInventoryDoesNothingWithEnoughSpace(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10}
	got, err := inventory(deps{}, character, nil, MoveOptions{})
	if err != nil || got != character {
		t.Fatalf("inventory() = %#v, %v, want unchanged character", got, err)
	}
}

func TestInventoryUsesInjectedDeposit(t *testing.T) {
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 2},
	}
	character := schemas.CharacterSchema{
		Name:              "hero",
		Inventory:         &inventoryItems,
		InventoryMaxItems: 5,
		Layer:             "overworld",
	}
	called := false
	d := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			called = true
			if len(items) != 1 || items[0].Code != "copper_ore" {
				t.Fatalf("deposited items = %#v, want copper ore", items)
			}
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := inventory(d, character, nil, MoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("deposit was not called")
	}
}

func TestInventoryReturnsMoveError(t *testing.T) {
	wantErr := errors.New("move failed")
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{}, wantErr
		},
	}
	character := schemas.CharacterSchema{
		InventoryMaxItems: 5,
		Layer:             "overworld",
	}

	_, err := inventory(d, character, nil, MoveOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("inventory() error = %v, want %v", err, wantErr)
	}
}

func TestInventoryReturnsDepositError(t *testing.T) {
	wantErr := errors.New("deposit failed")
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 2},
	}
	character := schemas.CharacterSchema{
		Inventory:         &inventoryItems,
		InventoryMaxItems: 5,
		Layer:             "overworld",
	}
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}

	_, err := inventory(d, character, nil, MoveOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("inventory() error = %v, want %v", err, wantErr)
	}
}

func TestGetInventoryItemsIncludesUnknownItem(t *testing.T) {
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "unknown_item", Quantity: 1},
	}
	character := schemas.CharacterSchema{Inventory: &inventoryItems}

	options := InventoryItemsOptions{KeepTypes: []string{"resource"}}
	items := GetInventoryItems(character, options)
	if len(items) != 1 || items[0].Code != "unknown_item" {
		t.Fatalf("GetInventoryItems() = %#v, want unknown item", items)
	}
}

func TestGetInventoryItemsSkipsEmptySlots(t *testing.T) {
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "", Quantity: 0},
		{Code: "copper_ore", Quantity: 0},
		{Code: "", Quantity: 1},
	}
	character := schemas.CharacterSchema{Inventory: &inventoryItems}

	options := InventoryItemsOptions{KeepTravelPotions: true}
	items := GetInventoryItems(character, options)
	if len(items) != 0 {
		t.Fatalf("GetInventoryItems() = %#v, want no items", items)
	}
}

func TestGetInventoryItemsPreservesOneTravelPotion(t *testing.T) {
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "forest_bank_potion", Quantity: 3},
		{Code: "recall_potion", Quantity: 1},
		{Code: "copper_ore", Quantity: 2},
	}
	character := schemas.CharacterSchema{Inventory: &inventoryItems}

	options := InventoryItemsOptions{KeepTravelPotions: true}
	got := GetInventoryItems(character, options)
	valid := len(got) == 2
	valid = valid && got[0].Code == "forest_bank_potion"
	valid = valid && got[0].Quantity == 2
	valid = valid && got[1].Code == "copper_ore"
	valid = valid && got[1].Quantity == 2
	if !valid {
		t.Fatalf("GetInventoryItems() = %#v, want excess potion and ore", got)
	}

	options.KeepTravelPotions = false
	got = GetInventoryItems(character, options)
	if len(got) != 3 || got[0].Quantity != 3 || got[1].Quantity != 1 {
		t.Fatalf("GetInventoryItems() without travel reserve = %#v, want all items", got)
	}
}

func TestKeepTravelReserveFiltersCustomDepositList(t *testing.T) {
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "forest_bank_potion", Quantity: 2},
		{Code: "recall_potion", Quantity: 1},
	}
	character := schemas.CharacterSchema{Inventory: &inventoryItems}
	items := []schemas.SimpleItemSchema{
		{Code: "forest_bank_potion", Quantity: 2},
		{Code: "recall_potion", Quantity: 1},
		{Code: "copper_ore", Quantity: 3},
	}

	got := KeepTravelReserve(character, items)
	valid := len(got) == 2
	valid = valid && got[0].Code == "forest_bank_potion"
	valid = valid && got[0].Quantity == 1
	valid = valid && got[1].Code == "copper_ore"
	valid = valid && got[1].Quantity == 3
	if !valid {
		t.Fatalf("KeepTravelReserve() = %#v, want excess potion and ore", got)
	}
}

func TestSpaceAfterDepositKeepsTravelPotionReserve(t *testing.T) {
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "forest_bank_potion", Quantity: 2},
		{Code: "recall_potion", Quantity: 1},
		{Code: "copper_ore", Quantity: 4},
	}
	character := schemas.CharacterSchema{
		Inventory:         &inventoryItems,
		InventoryMaxItems: 10,
	}

	got := SpaceAfterDeposit(character)
	if got != 8 {
		t.Fatalf("SpaceAfterDeposit() = %d, want 8", got)
	}
}

func TestShouldKeepItemRejectsNonMatchingType(t *testing.T) {
	if shouldKeepItem("copper_ore", []string{"food"}) {
		t.Fatal("shouldKeepItem() = true, want false")
	}
}
