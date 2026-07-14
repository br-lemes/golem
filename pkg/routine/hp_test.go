package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestHp(t *testing.T) {
	tests := []struct {
		name          string
		character     schemas.CharacterSchema
		minHp         int
		mockUse       func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error)
		mockRest      func(name string) (schemas.CharacterRestDataSchema, error)
		expectHp      int
		expectRest    bool
		expectUseQty  int
		expectUseCode string
		expectError   bool
	}{
		{
			name: "above minimum",
			character: schemas.CharacterSchema{
				Name:  "Hero",
				Hp:    310,
				MaxHp: 350,
			},
			minHp:       300,
			expectHp:    310,
			expectRest:  false,
			expectError: false,
		},
		{
			name: "minor damage rest",
			character: schemas.CharacterSchema{
				Name:  "Hero",
				Hp:    348,
				MaxHp: 350,
				Inventory: &[]schemas.InventorySlotSchema{
					{Code: "cooked_chicken", Quantity: 10, Slot: 0},
				},
			},
			minHp: 349,
			mockRest: func(name string) (schemas.CharacterRestDataSchema, error) {
				return schemas.CharacterRestDataSchema{
					Character: schemas.CharacterSchema{
						Name: name, Hp: 350, MaxHp: 350},
				}, nil
			},
			expectHp:    350,
			expectRest:  true,
			expectError: false,
		},
		{
			name: "nil inventory",
			character: schemas.CharacterSchema{
				Name:      "Hero",
				Hp:        100,
				MaxHp:     350,
				Inventory: nil,
			},
			minHp: 200,
			mockRest: func(name string) (schemas.CharacterRestDataSchema, error) {
				return schemas.CharacterRestDataSchema{
					Character: schemas.CharacterSchema{
						Name: name, Hp: 350, MaxHp: 350},
				}, nil
			},
			expectHp:    350,
			expectRest:  true,
			expectError: false,
		},
		{
			name: "exact healing loop break",
			character: schemas.CharacterSchema{
				Name:  "Hero",
				Hp:    190,
				MaxHp: 350,
				Level: 10,
				Inventory: &[]schemas.InventorySlotSchema{
					{Code: "cooked_chicken", Quantity: 5, Slot: 0},
					{Code: "apple", Quantity: 10, Slot: 1},
				},
			},
			minHp: 250,
			mockUse: func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
				if item.Code == "cooked_chicken" && item.Quantity == 2 {
					return schemas.UseItemSchema{
						Character: schemas.CharacterSchema{
							Name: name, Hp: 350, MaxHp: 350,
							Inventory: &[]schemas.InventorySlotSchema{
								{Code: "cooked_chicken", Quantity: 3, Slot: 0},
								{Code: "apple", Quantity: 10, Slot: 1},
							},
						},
					}, nil
				}
				return schemas.UseItemSchema{},
					errors.New("processed inventory after full recovery")
			},
			expectHp:      350,
			expectUseCode: "cooked_chicken",
			expectUseQty:  2,
			expectError:   false,
		},
		{
			name: "ceiled division healing",
			character: schemas.CharacterSchema{
				Name:  "Hero",
				Hp:    240,
				MaxHp: 350,
				Level: 10,
				Inventory: &[]schemas.InventorySlotSchema{
					{Code: "cooked_chicken", Quantity: 5, Slot: 0},
				},
			},
			minHp: 260,
			mockUse: func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
				if item.Code == "cooked_chicken" && item.Quantity == 2 {
					return schemas.UseItemSchema{
						Character: schemas.CharacterSchema{
							Name: name, Hp: 350, MaxHp: 350},
					}, nil
				}
				return schemas.UseItemSchema{},
					errors.New("should have requested exactly 2 items")
			},
			expectHp:      350,
			expectUseCode: "cooked_chicken",
			expectUseQty:  2,
			expectError:   false,
		},
		{
			name: "insufficient food level",
			character: schemas.CharacterSchema{
				Name:  "Newbie",
				Hp:    50,
				MaxHp: 350,
				Level: 1,
				Inventory: &[]schemas.InventorySlotSchema{
					{Code: "cooked_beef", Quantity: 5, Slot: 0},
				},
			},
			minHp: 200,
			mockRest: func(name string) (schemas.CharacterRestDataSchema, error) {
				return schemas.CharacterRestDataSchema{
					Character: schemas.CharacterSchema{
						Name: name, Hp: 350, MaxHp: 350},
				}, nil
			},
			expectHp:    350,
			expectRest:  true,
			expectError: false,
		},
		{
			name: "fallback waterfall combo",
			character: schemas.CharacterSchema{
				Name:  "Hero",
				Hp:    140,
				MaxHp: 350,
				Level: 10,
				Inventory: &[]schemas.InventorySlotSchema{
					{Code: "apple", Quantity: 5, Slot: 0},
					{Code: "cooked_chicken", Quantity: 2, Slot: 1},
				},
			},
			minHp: 200,
			mockUse: func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
				if item.Code == "cooked_chicken" && item.Quantity == 2 {
					return schemas.UseItemSchema{
						Character: schemas.CharacterSchema{
							Name: name, Hp: 300, MaxHp: 350,
							Inventory: &[]schemas.InventorySlotSchema{
								{Code: "apple", Quantity: 5, Slot: 0},
								{Code: "cooked_chicken", Quantity: 0, Slot: 1},
							},
						},
					}, nil
				}
				if item.Code == "apple" && item.Quantity == 1 {
					return schemas.UseItemSchema{
						Character: schemas.CharacterSchema{
							Name: name, Hp: 350, MaxHp: 350},
					}, nil
				}
				return schemas.UseItemSchema{},
					errors.New("unexpected food item combo or quantity")
			},
			expectHp:    350,
			expectError: false,
		},
		{
			name: "out of food fallback rest",
			character: schemas.CharacterSchema{
				Name:  "Hero",
				Hp:    150,
				MaxHp: 350,
				Level: 10,
				Inventory: &[]schemas.InventorySlotSchema{
					{Code: "cooked_chicken", Quantity: 1, Slot: 0},
				},
			},
			minHp: 200,
			mockUse: func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
				return schemas.UseItemSchema{
					Character: schemas.CharacterSchema{
						Name: name, Hp: 230, MaxHp: 350},
				}, nil
			},
			mockRest: func(name string) (schemas.CharacterRestDataSchema, error) {
				return schemas.CharacterRestDataSchema{
					Character: schemas.CharacterSchema{
						Name: name, Hp: 350, MaxHp: 350},
				}, nil
			},
			expectHp:    350,
			expectRest:  true,
			expectError: false,
		},
		{
			name: "use item api error",
			character: schemas.CharacterSchema{
				Name:  "Hero",
				Hp:    100,
				MaxHp: 350,
				Level: 10,
				Inventory: &[]schemas.InventorySlotSchema{
					{Code: "cooked_chicken", Quantity: 5, Slot: 0},
				},
			},
			minHp: 200,
			mockUse: func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
				return schemas.UseItemSchema{}, errors.New("timeout")
			},
			expectError: true,
		},
		{
			name: "rest api error",
			character: schemas.CharacterSchema{
				Name:      "InjuredHero",
				Hp:        100,
				MaxHp:     350,
				Inventory: nil,
			},
			minHp: 200,
			mockRest: func(name string) (schemas.CharacterRestDataSchema, error) {
				return schemas.CharacterRestDataSchema{}, errors.New("timeout")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restCalled := false
			useCalledWithQty := 0
			useCalledWithCode := ""
			dMock := deps{
				myActionUse: func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
					useCalledWithQty = item.Quantity
					useCalledWithCode = item.Code
					if tt.mockUse != nil {
						return tt.mockUse(name, item)
					}
					return schemas.UseItemSchema{}, nil
				},
				myActionRest: func(name string) (schemas.CharacterRestDataSchema, error) {
					restCalled = true
					if tt.mockRest != nil {
						return tt.mockRest(name)
					}
					return schemas.CharacterRestDataSchema{}, nil
				},
			}
			resultChar, err := hp(dMock, tt.character, tt.minHp)
			if (err != nil) != tt.expectError {
				t.Fatalf(
					"hp() unexpected error: expected error = %v, got = %v",
					tt.expectError, err)
			}
			if tt.expectError {
				return
			}
			if resultChar.Hp != tt.expectHp {
				t.Errorf("hp() final HP discrepancy: expected = %d, got = %d",
					tt.expectHp, resultChar.Hp)
			}
			if restCalled != tt.expectRest {
				t.Errorf(
					"hp() rest action call mismatch: expected = %v, got = %v",
					tt.expectRest, restCalled)
			}
			if tt.expectUseQty > 0 && useCalledWithQty != tt.expectUseQty {
				t.Errorf(
					"hp() consumed quantity mismatch: expected = %d, got = %d",
					tt.expectUseQty, useCalledWithQty)
			}
			if tt.expectUseCode != "" && useCalledWithCode != tt.expectUseCode {
				t.Errorf(
					"hp() consumed item code mismatch: expected = %s, got = %s",
					tt.expectUseCode, useCalledWithCode)
			}
		})
	}
}
