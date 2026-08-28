package routine

import (
	"errors"
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestDepositItems(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 4},
		{Code: "copper_boots", Quantity: 1},
	}
	character := schemas.CharacterSchema{
		Name:              "hero",
		Inventory:         &inventory,
		InventoryMaxItems: 100,
		X:                 3,
		Y:                 1,
		MapId:             334,
		Layer:             "overworld",
	}
	depositedItems := []schemas.SimpleItemSchema{}
	deps := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			depositedItems = items
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	got, err := deposit(deps, character, nil)
	if err != nil {
		t.Fatal(err)
	}
	wantItems := []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 4},
		{Code: "copper_boots", Quantity: 1},
	}
	if !reflect.DeepEqual(depositedItems, wantItems) {
		t.Fatalf("deposited items = %#v, want %#v", depositedItems, wantItems)
	}
	if got.Name != character.Name {
		t.Fatalf("returned character = %#v, want %#v", got, character)
	}
}

func TestDepositGold(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		Gold:  12,
		X:     3,
		Y:     1,
		MapId: 334,
		Layer: "overworld",
	}
	depositedGold := 0
	deps := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositGold: func(_ string, quantity int) (schemas.BankGoldTransactionSchema, error) {
			depositedGold = quantity
			return schemas.BankGoldTransactionSchema{Character: character}, nil
		},
	}

	_, err := deposit(deps, character, nil)
	if err != nil {
		t.Fatal(err)
	}
	if depositedGold != character.Gold {
		t.Fatalf("deposited gold = %d, want %d", depositedGold, character.Gold)
	}
}

func TestDepositKeepsItems(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 1},
	}
	character := schemas.CharacterSchema{
		Name:      "hero",
		Inventory: &inventory,
		X:         3,
		Y:         1,
		MapId:     334,
		Layer:     "overworld",
	}
	itemsCalled := false
	deps := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, _ []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			itemsCalled = true
			return schemas.BankItemTransactionSchema{}, nil
		},
	}

	_, err := deposit(deps, character, []string{"mining"})
	if err != nil {
		t.Fatal(err)
	}
	if itemsCalled {
		t.Fatal("item was deposited despite matching keep subtype")
	}
}

func TestDepositKeepsGold(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		Gold:  12,
		X:     3,
		Y:     1,
		MapId: 334,
		Layer: "overworld",
	}
	goldCalled := false
	deps := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositGold: func(_ string, _ int) (schemas.BankGoldTransactionSchema, error) {
			goldCalled = true
			return schemas.BankGoldTransactionSchema{}, nil
		},
	}

	_, err := deposit(deps, character, []string{"gold"})
	if err != nil {
		t.Fatal(err)
	}
	if goldCalled {
		t.Fatal("gold was deposited despite being kept")
	}
}

func TestDepositReturnsError(t *testing.T) {
	wantErr := errors.New("deposit failed")
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     3,
		Y:     1,
		Gold:  1,
		MapId: 334,
		Layer: "overworld",
	}
	deps := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositGold: func(_ string, _ int) (schemas.BankGoldTransactionSchema, error) {
			return schemas.BankGoldTransactionSchema{}, wantErr
		},
	}

	_, err := deposit(deps, character, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("deposit() error = %v, want %v", err, wantErr)
	}
}

func TestDepositReturnsMoveError(t *testing.T) {
	wantErr := errors.New("move failed")
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     3,
		Y:     1,
		MapId: 334,
		Layer: "overworld",
	}
	deps := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{}, wantErr
		},
	}

	_, err := deposit(deps, character, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("deposit() error = %v, want %v", err, wantErr)
	}
}

func TestDepositReturnsItemError(t *testing.T) {
	wantErr := errors.New("item deposit failed")
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 1},
	}
	character := schemas.CharacterSchema{
		Name:      "hero",
		Inventory: &inventory,
		X:         3,
		Y:         1,
		MapId:     334,
		Layer:     "overworld",
	}
	deps := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, _ []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}

	_, err := deposit(deps, character, nil)
	if !errors.Is(err, wantErr) {
		t.Fatalf("deposit() error = %v, want %v", err, wantErr)
	}
}
