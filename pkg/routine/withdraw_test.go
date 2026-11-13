package routine

import (
	"errors"
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func atBankCharacter() schemas.CharacterSchema {
	return schemas.CharacterSchema{
		Name:              "hero",
		InventoryMaxItems: 100,
		X:                 4,
		Y:                 1,
		Layer:             "overworld",
	}
}

func TestWithdrawMovesToBankAndWithdrawsItems(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:              "hero",
		InventoryMaxItems: 100,
		X:                 3,
		Y:                 1,
		MapId:             334,
		Layer:             "overworld",
	}
	moved := false
	var withdrawn []schemas.SimpleItemSchema
	d := deps{
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			moved = true
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "copper_ore", Quantity: 10},
			}, nil
		},
		myActionBankWithdrawItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			withdrawn = items
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	got, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 4},
	}, WithdrawOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !moved {
		t.Fatal("character did not move to the bank")
	}
	want := []schemas.SimpleItemSchema{{Code: "copper_ore", Quantity: 4}}
	if !reflect.DeepEqual(withdrawn, want) {
		t.Fatalf("withdrawn items = %#v, want %#v", withdrawn, want)
	}
	if got.Name != character.Name {
		t.Fatalf("returned character = %#v, want %#v", got, character)
	}
}

func TestWithdrawDoesNotMoveWhenAlreadyAtBank(t *testing.T) {
	character := atBankCharacter()
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			t.Fatal("movement action called while already at bank")
			return schemas.CharacterMovementDataSchema{}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "copper_ore", Quantity: 10},
			}, nil
		},
		myActionBankWithdrawItem: func(_ string, _ []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 1},
	}, WithdrawOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithdrawRejectsDuplicateItems(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:              "hero",
		InventoryMaxItems: 100,
		X:                 3,
		Y:                 1,
		MapId:             334,
		Layer:             "overworld",
	}
	bankCalled := false
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			t.Fatal("movement action called despite duplicate items")
			return schemas.CharacterMovementDataSchema{}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			bankCalled = true
			return nil, nil
		},
	}

	items := []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 3},
		{Code: "copper_ore", Quantity: 4},
		{Code: "iron_sword", Quantity: 1},
	}
	_, err := withdraw(d, character, items, WithdrawOptions{})
	if err == nil {
		t.Fatal("expected duplicate item error")
	}
	if bankCalled {
		t.Fatal("bank items loaded despite duplicate items")
	}
}

func TestWithdrawRestocksTravelPotionsBeforeWithdrawing(t *testing.T) {
	character := atBankCharacter()
	inventory := []schemas.InventorySlotSchema{}
	character.Inventory = &inventory
	var calls [][]schemas.SimpleItemSchema
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "recall_potion", Quantity: 1},
				{Code: "copper_ore", Quantity: 5},
			}, nil
		},
		myActionBankWithdrawItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			calls = append(calls, items)
			for _, item := range items {
				inventory = append(inventory, schemas.InventorySlotSchema{
					Code:     item.Code,
					Quantity: item.Quantity,
				})
			}
			character.Inventory = &inventory
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	got, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 1},
	}, WithdrawOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(calls) != 2 || calls[0][0].Code != "recall_potion" || calls[1][0].Code != "copper_ore" {
		t.Fatalf("withdraw calls = %#v, want recall potion then copper ore", calls)
	}
	if inventoryItemQuantity(got, "recall_potion") != 1 {
		t.Fatalf("travel potion was not kept after withdrawal: %#v", got.Inventory)
	}
}

func TestWithdrawRejectsInsufficientBankStock(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:              "hero",
		InventoryMaxItems: 100,
		X:                 3,
		Y:                 1,
		MapId:             334,
		Layer:             "overworld",
	}
	withdrawCalled := false
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			t.Fatal("movement action called despite insufficient bank stock")
			return schemas.CharacterMovementDataSchema{}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "copper_ore", Quantity: 5}}, nil
		},
		myActionBankWithdrawItem: func(_ string, _ []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			withdrawCalled = true
			return schemas.BankItemTransactionSchema{}, nil
		},
	}

	_, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 6},
	}, WithdrawOptions{})
	if err == nil {
		t.Fatal("expected insufficient bank stock error")
	}
	if withdrawCalled {
		t.Fatal("items were withdrawn despite insufficient bank stock")
	}
}

func TestWithdrawRejectsInsufficientInventorySpace(t *testing.T) {
	character := atBankCharacter()
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 99},
	}
	character.Inventory = &inventory
	withdrawCalled := false
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "copper_ore", Quantity: 10},
			}, nil
		},
		myActionBankWithdrawItem: func(_ string, _ []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			withdrawCalled = true
			return schemas.BankItemTransactionSchema{}, nil
		},
	}

	_, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 2},
	}, WithdrawOptions{})
	if err == nil {
		t.Fatal("expected insufficient inventory space error")
	}
	if withdrawCalled {
		t.Fatal("items were withdrawn despite insufficient inventory space")
	}
}

func TestWithdrawEmptyItemsDoesNothing(t *testing.T) {
	character := schemas.CharacterSchema{Name: "hero"}
	got, err := withdraw(deps{}, character, nil, WithdrawOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got != character {
		t.Fatalf("returned character = %#v, want %#v", got, character)
	}
}

func TestWithdrawReturnsMoveError(t *testing.T) {
	wantErr := errors.New("move failed")
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     3,
		Y:     1,
		MapId: 334,
		Layer: "overworld",
	}
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "copper_ore", Quantity: 10},
			}, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{}, wantErr
		},
	}

	_, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 1},
	}, WithdrawOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("withdraw() error = %v, want %v", err, wantErr)
	}
}

func TestWithdrawReturnsBankItemsError(t *testing.T) {
	wantErr := errors.New("bank unavailable")
	character := atBankCharacter()
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, wantErr
		},
	}

	_, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 1},
	}, WithdrawOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("withdraw() error = %v, want %v", err, wantErr)
	}
}

func TestWithdrawReturnsWithdrawError(t *testing.T) {
	wantErr := errors.New("withdraw failed")
	character := atBankCharacter()
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "copper_ore", Quantity: 10},
			}, nil
		},
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}

	_, err := withdraw(d, character, []schemas.SimpleItemSchema{
		{Code: "copper_ore", Quantity: 1},
	}, WithdrawOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("withdraw() error = %v, want %v", err, wantErr)
	}
}
