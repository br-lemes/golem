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

	got, err := deposit(deps, character, DepositOptions{})
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

	_, err := deposit(deps, character, DepositOptions{})
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

	options := DepositOptions{KeepTypes: []string{"mining"}}
	_, err := deposit(deps, character, options)
	if err != nil {
		t.Fatal(err)
	}
	if itemsCalled {
		t.Fatal("item was deposited despite matching keep subtype")
	}
}

func TestDepositKeepTypeOverridesNoTeleport(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "forest_bank_potion", Quantity: 1},
		{Code: "recall_potion", Quantity: 1},
	}
	character := schemas.CharacterSchema{
		Name:      "hero",
		Inventory: &inventory,
		X:         3,
		Y:         1,
		MapId:     334,
		Layer:     "overworld",
	}
	itemsDeposited := false
	d := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, _ []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			itemsDeposited = true
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := deposit(d, character, DepositOptions{
		Movement:  MoveOptions{NoTeleport: true},
		KeepTypes: []string{"potion"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if itemsDeposited {
		t.Fatal("teleport potions were deposited despite --keep potion")
	}
}

func TestDepositNoTeleportDepositsReservedPotions(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "forest_bank_potion", Quantity: 1},
		{Code: "recall_potion", Quantity: 1},
	}
	character := schemas.CharacterSchema{
		Name:      "hero",
		Inventory: &inventory,
		X:         3,
		Y:         1,
		MapId:     334,
		Layer:     "overworld",
	}
	var deposited []schemas.SimpleItemSchema
	d := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			deposited = items
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := deposit(d, character, DepositOptions{
		Movement: MoveOptions{NoTeleport: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	valid := len(deposited) == 2
	valid = valid && deposited[0].Code == "forest_bank_potion"
	valid = valid && deposited[0].Quantity == 1
	valid = valid && deposited[1].Code == "recall_potion"
	valid = valid && deposited[1].Quantity == 1
	if !valid {
		t.Fatalf("deposited items = %#v, want one of each reserved teleport potion", deposited)
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

	options := DepositOptions{KeepTypes: []string{"gold"}}
	_, err := deposit(deps, character, options)
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

	_, err := deposit(deps, character, DepositOptions{})
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

	_, err := deposit(deps, character, DepositOptions{})
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

	_, err := deposit(deps, character, DepositOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("deposit() error = %v, want %v", err, wantErr)
	}
}
