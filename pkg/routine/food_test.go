package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFoodCurrentQtyCountsOnlyPositiveFood(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "cooked_chicken", Quantity: 3},
		{Code: "cooked_beef", Quantity: 0},
		{Code: "copper_ore", Quantity: 5},
	}

	got := foodCurrentQty(schemas.CharacterSchema{Inventory: &inventory})
	if got != 3 {
		t.Fatalf("foodCurrentQty() = %d, want 3", got)
	}
}

func TestFoodRestockReturnsEarlyForEmptyFood(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10}
	got, err := foodRestock(deps{}, character, "", nil, false)
	if err != nil || got != character {
		t.Fatalf("foodRestock() = %#v, %v, want unchanged character", got, err)
	}
}

func TestFoodRestockFoodOnlyRejectsOtherFood(t *testing.T) {
	character := schemas.CharacterSchema{Level: 10, InventoryMaxItems: 10}
	got, err := foodRestock(deps{}, character, "apple", map[string]int{
		"cooked_chicken": 5,
	}, true)
	if err != nil || got != character {
		t.Fatalf("foodRestock() = %#v, %v, want unchanged character", got, err)
	}
}

func TestFoodRestockWithdrawsFood(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:              "hero",
		InventoryMaxItems: 10,
		Level:             10,
	}
	called := false
	d := deps{
		myActionBankWithdrawItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			called = true
			if len(items) != 1 || items[0].Code != "apple" {
				t.Fatalf("withdrawn items = %#v, want apple", items)
			}
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := foodRestock(d, character, "apple", map[string]int{"apple": 10}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("food withdrawal was not called")
	}
}

func TestFoodRestockReturnsWithdrawError(t *testing.T) {
	wantErr := errors.New("withdraw failed")
	d := deps{
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}
	character := schemas.CharacterSchema{InventoryMaxItems: 10, Level: 10}

	_, err := foodRestock(d, character, "apple", map[string]int{"apple": 10}, false)
	if !errors.Is(err, wantErr) {
		t.Fatalf("foodRestock() error = %v, want %v", err, wantErr)
	}
}

func TestFoodRestockAcceptsFoodOnly(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10, Level: 10}
	called := false
	d := deps{
		myActionBankWithdrawItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			called = true
			if items[0].Code != "apple" {
				t.Fatalf("withdrawn food = %#v, want apple", items)
			}
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := foodRestock(d, character, "apple", map[string]int{"apple": 10}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("food withdrawal was not called")
	}
}

func TestFoodRestockHandlesLimitedFood(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10, Level: 10}
	_, err := foodRestock(deps{
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}, character, "apple", map[string]int{"apple": 1}, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFoodRestockHandlesNoInventorySpace(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 1},
	}
	character := schemas.CharacterSchema{
		Inventory:         &inventory,
		InventoryMaxItems: 1,
		Level:             10,
	}
	_, err := foodRestock(deps{}, character, "apple", map[string]int{"apple": 1}, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFoodRestockStopsWhenInventoryFills(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10, Level: 10}
	d := deps{
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			inventory := []schemas.InventorySlotSchema{
				{Code: "apple", Quantity: 10},
			}
			character.Inventory = &inventory
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}
	_, err := foodRestock(d, character, "apple", map[string]int{
		"apple":          2,
		"cooked_chicken": 10,
	}, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFoodRestockLimitsChunkToInventorySpace(t *testing.T) {
	character := schemas.CharacterSchema{InventoryMaxItems: 10, Level: 10}
	d := deps{
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			inventory := []schemas.InventorySlotSchema{
				{Code: "apple", Quantity: 9},
			}
			character.Inventory = &inventory
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := foodRestock(d, character, "apple", map[string]int{
		"apple":          2,
		"cooked_chicken": 10,
	}, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSelectFoodMatchesAutomaticRestockOrder(t *testing.T) {
	character := schemas.CharacterSchema{Level: 10}

	got := SelectFood(character, map[string]int{
		"apple":          10,
		"cooked_chicken": 10,
	})

	if got != "cooked_chicken" {
		t.Fatalf("SelectFood() = %q, want %q", got, "cooked_chicken")
	}
}

func TestSelectFoodIgnoresUnavailableOrTooHighLevelFood(t *testing.T) {
	character := schemas.CharacterSchema{Level: 1}

	got := SelectFood(character, map[string]int{
		"cooked_chicken": 0,
		"cooked_beef":    10,
	})

	if got != "" {
		t.Fatalf("SelectFood() = %q, want no food", got)
	}
}
