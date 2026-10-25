package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const MaxInventoryReserve = 90

type InventoryPlan struct {
	Capacity                int `json:"capacity"`
	FixedItems              int `json:"fixed_items"`
	Available               int `json:"available"`
	FoodReserve             int `json:"food_reserve"`
	FoodQuantity            int `json:"food_quantity"`
	Utility1Reserve         int `json:"utility1_reserve"`
	Utility1ReserveQuantity int `json:"utility1_reserve_quantity"`
	Utility2Reserve         int `json:"utility2_reserve"`
	Utility2ReserveQuantity int `json:"utility2_reserve_quantity"`
	PlannedItems            int `json:"planned_items"`
	RemainingItems          int `json:"remaining_items"`
}

func ValidateInventoryReserve(foodReserve, utility1Reserve, utility2Reserve int) error {
	reserves := []struct {
		name  string
		value int
	}{
		{"food", foodReserve},
		{"utility1", utility1Reserve},
		{"utility2", utility2Reserve},
	}
	for _, reserve := range reserves {
		if reserve.value < 0 || reserve.value > MaxInventoryReserve {
			return fmt.Errorf("%s reserve must be between 0 and %d", reserve.name, MaxInventoryReserve)
		}
	}
	if foodReserve+utility1Reserve+utility2Reserve > MaxInventoryReserve {
		return fmt.Errorf("food and utility reserves must total at most %d", MaxInventoryReserve)
	}
	return nil
}

func PlanInventory(character schemas.CharacterSchema, noTeleport bool, foodReserve, utility1Reserve, utility2Reserve int) (InventoryPlan, error) {
	err := ValidateInventoryReserve(foodReserve, utility1Reserve, utility2Reserve)
	if err != nil {
		return InventoryPlan{}, err
	}

	fixedItems := 0
	if !noTeleport {
		fixedItems = len(reservedTravelPotionCodes)
	}
	capacity := character.InventoryMaxItems
	available := max(0, capacity-fixedItems)
	foodQuantity := available * foodReserve / 100
	utility1ReserveQuantity := available * utility1Reserve / 100
	utility2ReserveQuantity := available * utility2Reserve / 100
	plannedItems := fixedItems + foodQuantity + utility1ReserveQuantity + utility2ReserveQuantity
	return InventoryPlan{
		Capacity:                capacity,
		FixedItems:              fixedItems,
		Available:               available,
		FoodReserve:             foodReserve,
		FoodQuantity:            foodQuantity,
		Utility1Reserve:         utility1Reserve,
		Utility1ReserveQuantity: utility1ReserveQuantity,
		Utility2Reserve:         utility2Reserve,
		Utility2ReserveQuantity: utility2ReserveQuantity,
		PlannedItems:            plannedItems,
		RemainingItems:          capacity - plannedItems,
	}, nil
}
