package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestBankDoesNothingWhenNoActionIsNeeded(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10}
	called := false
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			called = true
			return []schemas.SimpleItemSchema{{Code: "bread", Quantity: 5}}, nil
		},
	}

	got, err := bank(d, character, BankOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("bank items were not loaded")
	}
	if got != character {
		t.Fatalf("returned character = %#v, want %#v", got, character)
	}
}

func TestBankReturnsBankItemsError(t *testing.T) {
	wantErr := errors.New("bank unavailable")
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, wantErr
		},
	}

	_, err := bank(d, schemas.CharacterSchema{}, BankOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("bank() error = %v, want %v", err, wantErr)
	}
}

func TestBankUsesInjectedMoveAndDepositsItems(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 2},
	}
	character := schemas.CharacterSchema{
		Name:              "hero",
		Inventory:         &inventory,
		InventoryMaxItems: 5,
		Layer:             "overworld",
	}
	moveCalled := false
	depositCalled := false
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, nil
		},
		myActionMove: func(_ string, x, y int) (schemas.CharacterMovementDataSchema, error) {
			moveCalled = true
			character.X = x
			character.Y = y
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			depositCalled = true
			if len(items) != 1 || items[0].Code != "copper_ore" || items[0].Quantity != 2 {
				t.Fatalf("deposited items = %#v, want copper ore x2", items)
			}
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := bank(d, character, BankOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !moveCalled || !depositCalled {
		t.Fatalf("move called = %t, deposit called = %t, want both", moveCalled, depositCalled)
	}
}

func TestBankReturnsMoveError(t *testing.T) {
	wantErr := errors.New("move failed")
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{}, wantErr
		},
	}
	character := schemas.CharacterSchema{
		InventoryMaxItems: 5,
		Layer:             "overworld",
	}

	_, err := bank(d, character, BankOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("bank() error = %v, want %v", err, wantErr)
	}
}

func TestBankReturnsUtilityCheckError(t *testing.T) {
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, nil
		},
	}
	_, err := bank(d, schemas.CharacterSchema{}, BankOptions{
		Utility1: "small_health_potion",
	})
	if err == nil {
		t.Fatal("bank() returned nil error, want depleted utility error")
	}
}

func TestBankReturnsDepositError(t *testing.T) {
	wantErr := errors.New("deposit failed")
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 2},
	}
	character := schemas.CharacterSchema{
		Inventory:         &inventory,
		InventoryMaxItems: 5,
		Layer:             "overworld",
	}
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) { return nil, nil },
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}
	_, err := bank(d, character, BankOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("bank() error = %v, want %v", err, wantErr)
	}
}

func TestBankReturnsUtilityUnequipError(t *testing.T) {
	wantErr := errors.New("unequip failed")
	d := utilityBankDeps("old_utility", func() error { return wantErr })
	character := schemas.CharacterSchema{
		Utility1Slot:      "old_utility",
		InventoryMaxItems: 10,
		Layer:             "overworld",
	}
	_, err := bank(d, character, BankOptions{Utility1: "small_health_potion"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("bank() error = %v, want %v", err, wantErr)
	}
}

func TestBankReturnsUtilityWithdrawError(t *testing.T) {
	wantErr := errors.New("withdraw failed")
	d := utilityBankDeps("", func() error { return wantErr })
	character := schemas.CharacterSchema{
		InventoryMaxItems: 10,
		Layer:             "overworld",
	}
	_, err := bank(d, character, BankOptions{Utility1: "small_health_potion"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("bank() error = %v, want %v", err, wantErr)
	}
}

func TestBankReturnsUtilityEquipError(t *testing.T) {
	wantErr := errors.New("equip failed")
	d := utilityBankDeps("", nil)
	d.myActionEquip = func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
		return schemas.EquipmentTransactionSchema{}, wantErr
	}
	character := schemas.CharacterSchema{
		InventoryMaxItems: 10,
		Layer:             "overworld",
	}
	_, err := bank(d, character, BankOptions{Utility1: "small_health_potion"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("bank() error = %v, want %v", err, wantErr)
	}
}

func TestBankReturnsFoodWithdrawError(t *testing.T) {
	wantErr := errors.New("food withdraw failed")
	character := schemas.CharacterSchema{
		InventoryMaxItems: 10,
		Layer:             "overworld",
	}
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "cooked_chicken", Quantity: 1},
			}, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}

	_, err := bank(d, character, BankOptions{Food: "cooked_chicken"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("bank() error = %v, want %v", err, wantErr)
	}
}

func utilityBankDeps(current string, unequipErr func() error) deps {
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "small_health_potion", Quantity: 1},
			}, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			character := schemas.CharacterSchema{
				InventoryMaxItems: 10,
				Layer:             "overworld",
			}
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			if unequipErr != nil {
				return schemas.BankItemTransactionSchema{}, unequipErr()
			}
			return schemas.BankItemTransactionSchema{}, nil
		},
	}
	if unequipErr != nil && current != "" {
		d.myActionUnequip = func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			return schemas.EquipmentTransactionSchema{}, unequipErr()
		}
	}
	return d
}
