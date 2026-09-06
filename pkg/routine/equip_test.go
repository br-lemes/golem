package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestEquipValidatesBeforeFetchingCharacter(t *testing.T) {
	charactersCalled := false
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			charactersCalled = true
			return schemas.CharacterSchema{}, nil
		},
	}
	quantity := 2
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon", Quantity: &quantity},
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if charactersCalled {
		t.Fatal("characters dependency was called")
	}
}

func TestEquipSkipsAlreadyEquippedItem(t *testing.T) {
	charactersCalled := false
	equipCalled := false
	character := schemas.CharacterSchema{Level: 10, WeaponSlot: "iron_sword"}
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			charactersCalled = true
			return character, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			equipCalled = true
			return schemas.EquipmentTransactionSchema{}, nil
		},
	}
	result, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !charactersCalled {
		t.Fatal("characters dependency was not called")
	}
	if equipCalled {
		t.Fatal("equip dependency was called")
	}
	if result.WeaponSlot != character.WeaponSlot {
		t.Fatalf("weapon slot = %q, want %q", result.WeaponSlot, character.WeaponSlot)
	}
}

func TestEquipClearsUtilitiesBeforeReturning(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "small_health_potion", Quantity: 2},
	}
	character := schemas.CharacterSchema{
		Level:                10,
		X:                    4,
		Y:                    1,
		MapId:                334,
		Layer:                "overworld",
		WeaponSlot:           "iron_sword",
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: 2,
		Inventory:            &inventory,
	}
	unequipCalled := false
	depositCalls := 0
	equipCalled := false
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return character, nil
		},
		myActionUnequip: func(_ string, slots []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			unequipCalled = len(slots) == 1 && slots[0].Slot == schemas.Utility1
			character.Utility1Slot = ""
			character.Utility1SlotQuantity = 0
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			depositCalls++
			if len(items) != 1 || items[0].Code != "small_health_potion" || items[0].Quantity != 2 {
				t.Fatalf("deposited items = %#v", items)
			}
			character.Inventory = nil
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			equipCalled = true
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
	}
	result, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !unequipCalled || depositCalls != 1 || equipCalled {
		t.Fatalf("unequip called = %t, deposits = %d, equip called = %t", unequipCalled, depositCalls, equipCalled)
	}
	if result.Utility1Slot != "" {
		t.Fatalf("utility1 slot = %q, want empty", result.Utility1Slot)
	}
}

func TestEquipReturnsUtilityClearError(t *testing.T) {
	wantErr := errors.New("clear utility failed")
	character := schemas.CharacterSchema{
		Level:        10,
		X:            4,
		Y:            1,
		MapId:        334,
		Layer:        "overworld",
		Utility1Slot: "small_health_potion",
	}
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return character, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionUnequip: func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			return schemas.EquipmentTransactionSchema{}, wantErr
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestEquipReturnsUtilityDepositError(t *testing.T) {
	wantErr := errors.New("deposit utility failed")
	inventory := []schemas.InventorySlotSchema{
		{Code: "small_health_potion", Quantity: 2},
	}
	character := schemas.CharacterSchema{
		Level:                10,
		X:                    4,
		Y:                    1,
		MapId:                334,
		Layer:                "overworld",
		WeaponSlot:           "iron_sword",
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: 2,
		Inventory:            &inventory,
	}
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return character, nil
		},
		myActionUnequip: func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			character.Utility1Slot = ""
			character.Utility1SlotQuantity = 0
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestEquipReturnsFinalDepositError(t *testing.T) {
	wantErr := errors.New("final deposit failed")
	inventory := []schemas.InventorySlotSchema{
		{Code: "small_health_potion", Quantity: 2},
	}
	character := schemas.CharacterSchema{
		Level:                10,
		X:                    4,
		Y:                    1,
		MapId:                334,
		Layer:                "overworld",
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: 2,
		Inventory:            &inventory,
	}
	depositCalls := 0
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return character, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
		},
		myActionUnequip: func(string, []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
			character.Utility1Slot = ""
			character.Utility1SlotQuantity = 0
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionBankDepositItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			depositCalls++
			if depositCalls == 2 {
				return schemas.BankItemTransactionSchema{}, wantErr
			}
			character.Inventory = nil
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			character.Inventory = &inventory
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestEquipmentItemsUsesSpecifiedQuantity(t *testing.T) {
	quantity := 7
	items := equipmentItems([]schemas.EquipSchema{{
		Code:     "small_health_potion",
		Slot:     schemas.Utility1,
		Quantity: &quantity,
	}})
	if len(items) != 1 || items[0].Code != "small_health_potion" || items[0].Quantity != quantity {
		t.Fatalf("equipmentItems() = %#v, want quantity %d", items, quantity)
	}
}

func TestValidateEquipments(t *testing.T) {
	quantity := 2
	utilityQuantity := 5
	tests := []struct {
		name       string
		equipments []schemas.EquipSchema
		wantErr    bool
	}{
		{
			name:       "missing code",
			equipments: []schemas.EquipSchema{{Slot: "weapon"}},
			wantErr:    true,
		},
		{
			name:       "missing slot",
			equipments: []schemas.EquipSchema{{Code: "iron_sword"}},
			wantErr:    true,
		},
		{
			name:       "unknown item",
			equipments: []schemas.EquipSchema{{Code: "unknown", Slot: "weapon"}},
			wantErr:    true,
		},
		{
			name:       "non equipment",
			equipments: []schemas.EquipSchema{{Code: "shell", Slot: "weapon"}},
			wantErr:    true,
		},
		{
			name: "invalid slot",
			equipments: []schemas.EquipSchema{
				{Code: "iron_sword", Slot: "ring1"},
			},
			wantErr: true,
		},
		{
			name: "slot conflict",
			equipments: []schemas.EquipSchema{
				{Code: "iron_sword", Slot: "weapon"},
				{Code: "sticky_sword", Slot: "weapon"},
			},
			wantErr: true,
		},
		{
			name: "non utility quantity",
			equipments: []schemas.EquipSchema{
				{Code: "iron_sword", Slot: "weapon", Quantity: &quantity},
			},
			wantErr: true,
		},
		{
			name: "utility quantity",
			equipments: []schemas.EquipSchema{{
				Code:     "small_health_potion",
				Slot:     "utility1",
				Quantity: &utilityQuantity,
			}},
		},
		{
			name: "equipment without conditions",
			equipments: []schemas.EquipSchema{{
				Code: "copper_boots",
				Slot: "boots",
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateEquipments(test.equipments)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateEquipments() error = %v, want error: %v", err, test.wantErr)
			}
		})
	}
}

func TestEquipReturnsCharacterError(t *testing.T) {
	wantErr := errors.New("characters failed")
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return schemas.CharacterSchema{}, wantErr
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestEquipReturnsLevelError(t *testing.T) {
	actionCalled := false
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return schemas.CharacterSchema{}, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			actionCalled = true
			return schemas.EquipmentTransactionSchema{}, nil
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if err == nil {
		t.Fatal("expected an error")
	}
	if actionCalled {
		t.Fatal("equip dependency was called")
	}
}

func TestEquipReturnsBankItemsError(t *testing.T) {
	wantErr := errors.New("bank failed")
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return schemas.CharacterSchema{Level: 10}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, wantErr
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestEquipReturnsActionError(t *testing.T) {
	wantErr := errors.New("equip action failed")
	inventory := []schemas.InventorySlotSchema{
		{Code: "iron_sword", Quantity: 1},
	}
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return schemas.CharacterSchema{Level: 10, Inventory: &inventory}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			return schemas.EquipmentTransactionSchema{}, wantErr
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestEquipUsesAvailableInventoryItem(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "iron_sword", Quantity: 1},
	}
	actionCalled := false
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return schemas.CharacterSchema{Level: 10, Inventory: &inventory}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, nil
		},
		myActionEquip: func(string, []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			actionCalled = true
			return schemas.EquipmentTransactionSchema{}, nil
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !actionCalled {
		t.Fatal("equip dependency was not called")
	}
}

func TestEquipWithdrawsMissingItem(t *testing.T) {
	withdrawCalled := false
	equipCalled := false
	character := equipTestCharacter()
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return character, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
		},
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankWithdrawItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			withdrawCalled = len(items) == 1 && items[0].Code == "iron_sword"
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
		myActionEquip: func(_ string, _ []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error) {
			equipCalled = true
			return schemas.EquipmentTransactionSchema{Character: character}, nil
		},
	}

	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !withdrawCalled || !equipCalled {
		t.Fatalf("withdraw called = %t, equip called = %t", withdrawCalled, equipCalled)
	}
}

func TestEquipReturnsMoveErrorWhenWithdrawing(t *testing.T) {
	wantErr := errors.New("move failed")
	character := equipTestCharacter()
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return character, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
		},
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{}, wantErr
		},
	}

	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestEquipReturnsWithdrawError(t *testing.T) {
	wantErr := errors.New("withdraw failed")
	character := equipTestCharacter()
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return character, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
		},
		myActionMove: func(_ string, _, _ int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
		myActionBankWithdrawItem: func(_ string, _ []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			return schemas.BankItemTransactionSchema{}, wantErr
		},
	}

	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("equip() error = %v, want %v", err, wantErr)
	}
}

func TestFilterNeededEquipmentsUsesQuantities(t *testing.T) {
	quantity := 2
	character := schemas.CharacterSchema{
		Utility1Slot:         "small_health_potion",
		Utility1SlotQuantity: 1,
	}
	equipments := []schemas.EquipSchema{{
		Code:     "small_health_potion",
		Slot:     "utility1",
		Quantity: &quantity,
	}}
	got := filterNeededEquipments(character, equipments)
	if len(got) != 1 {
		t.Fatalf("needed equipments = %#v, want one equipment", got)
	}
}

func TestValidateTotalStockUsesInventory(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "iron_sword", Quantity: 1},
	}
	deps := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, nil
		},
	}
	character := schemas.CharacterSchema{Inventory: &inventory}
	quantity := 1
	err := validateTotalStock(deps, character, []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon", Quantity: &quantity},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func equipTestCharacter() schemas.CharacterSchema {
	return schemas.CharacterSchema{
		Level: 10,
		X:     3,
		Y:     1,
		MapId: 334,
		Layer: "overworld",
	}
}

func TestEquipReturnsInsufficientStockError(t *testing.T) {
	deps := deps{
		characters: func(string) (schemas.CharacterSchema, error) {
			return schemas.CharacterSchema{Level: 10}, nil
		},
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{}, nil
		},
	}
	_, err := equip(deps, "hero", []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestValidateTotalStockUsesBank(t *testing.T) {
	deps := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 1}}, nil
		},
	}
	err := validateTotalStock(deps, schemas.CharacterSchema{}, []schemas.EquipSchema{
		{Code: "iron_sword", Slot: "weapon"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCalculateMissingItemsAccountsForInventory(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{
		{Code: "iron_sword", Quantity: 1},
	}
	quantity := 2
	character := schemas.CharacterSchema{Inventory: &inventory}
	equipments := []schemas.EquipSchema{{
		Code:     "iron_sword",
		Slot:     "weapon",
		Quantity: &quantity,
	}}
	missing := calculateMissingItems(character, equipments)
	if len(missing) != 1 {
		t.Fatalf("missing items = %#v, want one item", missing)
	}
	if missing[0].Code != "iron_sword" || missing[0].Quantity != 1 {
		t.Fatalf("missing item = %#v, want iron_sword x1", missing[0])
	}
}

func TestCheckLevelRequirementsUsesToolSkill(t *testing.T) {
	tests := []struct {
		name      string
		character schemas.CharacterSchema
		wantErr   bool
	}{
		{
			name:      "insufficient fishing level",
			character: schemas.CharacterSchema{Level: 50},
			wantErr:   true,
		},
		{
			name:      "sufficient fishing level",
			character: schemas.CharacterSchema{Level: 50, FishingLevel: 10},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkLevelRequirements(test.character, []schemas.EquipSchema{
				{Code: "spruce_fishing_rod", Slot: "weapon"},
			})
			if (err != nil) != test.wantErr {
				t.Fatalf("checkLevelRequirements() error = %v, want error: %v", err, test.wantErr)
			}
		})
	}
}
