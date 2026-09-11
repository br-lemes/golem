package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestInventoryDoesNothingWithEnoughSpace(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10}
	got, err := inventory(deps{}, character, nil)
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

	_, err := inventory(d, character, nil)
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

	_, err := inventory(d, character, nil)
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

	_, err := inventory(d, character, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("inventory() error = %v, want %v", err, wantErr)
	}
}

func TestGetInventoryItemsIncludesUnknownItem(t *testing.T) {
	inventoryItems := []schemas.InventorySlotSchema{
		{Code: "unknown_item", Quantity: 1},
	}
	character := schemas.CharacterSchema{Inventory: &inventoryItems}

	items := GetInventoryItems(character, []string{"resource"})
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

	items := GetInventoryItems(character, nil)
	if len(items) != 0 {
		t.Fatalf("GetInventoryItems() = %#v, want no items", items)
	}
}

func TestShouldKeepItemRejectsNonMatchingType(t *testing.T) {
	if shouldKeepItem("copper_ore", []string{"food"}) {
		t.Fatal("shouldKeepItem() = true, want false")
	}
}
