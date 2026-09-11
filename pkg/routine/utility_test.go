package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestClearUtilitiesSkipsEmptySlots(t *testing.T) {
	character := schemas.CharacterSchema{Utility2Slot: ""}
	got, err := clearUtilities(deps{}, character, []string{
		"utility1",
		"utility2",
	})
	if err != nil || got != character {
		t.Fatalf("clearUtilities() = %#v, %v, want unchanged character", got, err)
	}
}

func TestClearUtilitiesUnequipsConfiguredSlots(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:                 "hero",
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: 4,
		Utility2Slot:         "cooked_chicken",
		Utility2SlotQuantity: 7,
	}
	called := false
	d := deps{
		myActionUnequip: func(_ string, items []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			called = true
			if len(items) != 2 || items[0].Slot != "utility1" || items[1].Slot != "utility2" {
				t.Fatalf("unequips = %#v, want both utility slots", items)
			}
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
	}

	_, err := clearUtilities(d, character, []string{"utility1", "utility2"})
	if err != nil || !called {
		t.Fatalf("clearUtilities() error = %v, called = %t", err, called)
	}
}

func TestClearUtilitiesReturnsUnequipError(t *testing.T) {
	wantErr := errors.New("unequip failed")
	d := deps{
		myActionUnequip: func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			return schemas.EquipmentTransactionSchema{}, wantErr
		},
	}
	character := schemas.CharacterSchema{Utility1Slot: "small_health_potion"}

	_, err := clearUtilities(d, character, []string{"utility1"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("clearUtilities() error = %v, want %v", err, wantErr)
	}
}

func TestUtilityCheck(t *testing.T) {
	completed := schemas.CharacterSchema{
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: utilityMinThreshold,
	}
	needs, err := utilityCheck(completed, "small_health_potion", "", nil)
	if needs || err != nil {
		t.Fatalf("utilityCheck() = %t, %v, want no attention", needs, err)
	}

	needs, err = utilityCheck(schemas.CharacterSchema{}, "small_health_potion", "", map[string]int{
		"small_health_potion": 1,
	})
	if !needs || err != nil {
		t.Fatalf("utilityCheck() = %t, %v, want attention", needs, err)
	}
}

func TestUtilityCheckReturnsDepletedErrors(t *testing.T) {
	_, err := utilityCheck(schemas.CharacterSchema{}, "small_health_potion", "", nil)
	if err == nil {
		t.Fatal("utilityCheck() returned nil error, want depleted utility error")
	}

	_, err = utilityCheck(schemas.CharacterSchema{}, "small_health_potion", "small_health_potion", nil)
	if err == nil {
		t.Fatal("utilityCheck() returned nil error, want shared depleted utility error")
	}
}

func TestUtilityCheckHandlesUnavailableBankStock(t *testing.T) {
	character := schemas.CharacterSchema{
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: utilityMinThreshold - 1,
	}

	needs, err := utilityCheck(character, "small_health_potion", "", nil)
	if needs || err != nil {
		t.Fatalf("utilityCheck() = %t, %v, want no restock", needs, err)
	}
}

func TestCheckDepletedSkipsEmptySlots(t *testing.T) {
	slots := utilitySlots(schemas.CharacterSchema{}, "", "")
	err := checkDepleted("", "", slots, func(*utilitySlot) int { return 0 })
	if err != nil {
		t.Fatalf("checkDepleted() error = %v, want nil", err)
	}
}

func TestCheckDepletedAllowsSharedUtilityWithStock(t *testing.T) {
	character := schemas.CharacterSchema{
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: 1,
		Utility2Slot:         "small_health_potion",
	}
	slots := utilitySlots(character, "small_health_potion", "small_health_potion")
	err := checkDepleted("small_health_potion", "small_health_potion", slots, func(slot *utilitySlot) int {
		return slot.curQty
	})
	if err != nil {
		t.Fatalf("checkDepleted() error = %v, want nil", err)
	}
}

func TestUtilityRestockUnequipsWithdrawsAndEquips(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:              "hero",
		Utility1Slot:      "old_utility",
		InventoryMaxItems: 10,
	}
	unequipped, withdrawn, equipped := false, false, false
	d := deps{
		myActionUnequip: func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			unequipped = true
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			withdrawn = true
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			equipped = true
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
	}

	_, err := utilityRestock(d, character, "small_health_potion", "", map[string]int{
		"small_health_potion": 2,
	})
	if err != nil || !unequipped || !withdrawn || !equipped {
		t.Fatalf("utilityRestock() error = %v, calls = %t/%t/%t", err, unequipped, withdrawn, equipped)
	}
}

func TestUtilityRestockStopsWithoutInventorySpace(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 1},
	}
	character := schemas.CharacterSchema{
		Inventory:         &inventory,
		InventoryMaxItems: 1,
	}

	got, err := utilityRestock(deps{}, character, "small_health_potion", "", map[string]int{
		"small_health_potion": 1,
	})
	if err != nil || got != character {
		t.Fatalf("utilityRestock() = %#v, %v, want unchanged character", got, err)
	}
}

func TestUtilityRestockSkipsFullUtility(t *testing.T) {
	character := schemas.CharacterSchema{
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: utilityMaxStack,
	}
	d := deps{
		myActionUnequip: func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			t.Fatal("unequip should not be called")
			return schemas.EquipmentTransactionSchema{}, nil
		},
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			t.Fatal("withdraw should not be called")
			return schemas.BankItemTransactionSchema{}, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			t.Fatal("equip should not be called")
			return schemas.EquipmentTransactionSchema{}, nil
		},
	}

	got, err := utilityRestock(d, character, "small_health_potion", "", map[string]int{
		"small_health_potion": 1,
	})
	if err != nil || got != character {
		t.Fatalf("utilityRestock() = %#v, %v, want unchanged character", got, err)
	}
}

func TestUtilityRestockLimitsChunkToInventorySpace(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "copper_ore", Quantity: 98},
	}
	character := schemas.CharacterSchema{
		Inventory:         &inventory,
		InventoryMaxItems: 100,
	}
	called := false
	d := deps{
		myActionBankWithdrawItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			called = true
			if items[0].Quantity != 2 {
				t.Fatalf("withdrawn quantity = %d, want 2", items[0].Quantity)
			}
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
	}

	_, err := utilityRestock(d, character, "small_health_potion", "", map[string]int{
		"small_health_potion": 100,
	})
	if err != nil || !called {
		t.Fatalf("utilityRestock() error = %v, called = %t", err, called)
	}
}

func TestUtilityRestockReturnsErrors(t *testing.T) {
	checks := []struct {
		name string
		deps deps
		want error
	}{
		{
			name: "unequip",
			deps: deps{
				myActionUnequip: func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
					return schemas.EquipmentTransactionSchema{}, errors.New("unequip failed")
				},
			},
			want: errors.New("unequip failed"),
		},
		{
			name: "withdraw",
			deps: deps{
				myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
					return schemas.BankItemTransactionSchema{}, errors.New("withdraw failed")
				},
			},
			want: errors.New("withdraw failed"),
		},
		{
			name: "equip",
			deps: deps{
				myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
					return schemas.BankItemTransactionSchema{}, nil
				},
				myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
					return schemas.EquipmentTransactionSchema{}, errors.New("equip failed")
				},
			},
			want: errors.New("equip failed"),
		},
	}

	for _, test := range checks {
		t.Run(test.name, func(t *testing.T) {
			character := schemas.CharacterSchema{InventoryMaxItems: 10}
			if test.name == "unequip" {
				character.Utility1Slot = "old_utility"
			}
			_, err := utilityRestock(test.deps, character, "small_health_potion", "", map[string]int{
				"small_health_potion": 1,
			})
			if err == nil || err.Error() != test.want.Error() {
				t.Fatalf("utilityRestock() error = %v, want %v", err, test.want)
			}
		})
	}
}
