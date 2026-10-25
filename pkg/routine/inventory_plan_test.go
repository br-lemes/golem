package routine

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestPlanInventoryUsesOneSharedSpaceBase(t *testing.T) {
	plan, err := PlanInventory(schemas.CharacterSchema{InventoryMaxItems: 223}, false, 70, 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Available != 221 || plan.FoodQuantity != 154 || plan.Utility2ReserveQuantity != 44 {
		t.Fatalf("plan = %#v", plan)
	}
	if plan.PlannedItems != 200 || plan.RemainingItems != 23 {
		t.Fatalf("planned = %d, remaining = %d", plan.PlannedItems, plan.RemainingItems)
	}
}

func TestPlanInventoryRejectsPercentagesWithoutMargin(t *testing.T) {
	_, err := PlanInventory(schemas.CharacterSchema{InventoryMaxItems: 223}, false, 70, 21, 0)
	if err == nil {
		t.Fatal("PlanInventory() returned nil error")
	}
}
