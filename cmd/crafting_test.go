package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCraftingCapacityUsesRecipeOutputQuantity(t *testing.T) {
	output := 2
	item := schemas.ItemSchema{
		Craft: &schemas.CraftSchema{
			Items: &[]schemas.SimpleItemSchema{
				{Code: "water", Quantity: 3},
				{Code: "herb", Quantity: 2},
			},
			Quantity: &output,
		},
	}

	quantity := craftingCapacity(item, map[string]int{"water": 10, "herb": 7})
	if quantity != 6 {
		t.Fatalf("quantity = %d, want 6", quantity)
	}
}

func TestCraftingTarget(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		maximum  int
		want     int
		actions  int
		wantErr  bool
	}{
		{name: "maximum", quantity: 0, maximum: 6, want: 6, actions: 3},
		{name: "no materials", quantity: 0, maximum: 0, want: 0, actions: 0},
		{name: "exact quantity", quantity: 4, maximum: 6, want: 4, actions: 2},
		{name: "partial batch", quantity: 5, maximum: 6, wantErr: true},
		{name: "exceeds maximum", quantity: 8, maximum: 6, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			quantity, actions, err := craftingTarget(test.quantity, test.maximum, 2)
			if test.wantErr {
				if err == nil {
					t.Fatal("craftingTarget() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("craftingTarget() error = %v", err)
			}
			if quantity != test.want || actions != test.actions {
				t.Fatalf("craftingTarget() = (%d, %d), want (%d, %d)", quantity, actions, test.want, test.actions)
			}
		})
	}
}
