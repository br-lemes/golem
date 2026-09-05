package database

import "testing"

func TestItemsSmoke(t *testing.T) {
	if len(Items().All()) == 0 {
		t.Fatal("items catalog is empty")
	}
}

func TestItemCatalogViews(t *testing.T) {
	t.Run("equipments", func(t *testing.T) {
		equipments := Items().Equipments().All()
		if len(equipments) == 0 {
			t.Fatal("equipment view is empty")
		}
		for _, item := range equipments {
			_, exists := EquipmentTypeToSlots[item.Type]
			if !exists {
				t.Fatalf("equipment view contains %q with type %q", item.Code, item.Type)
			}
		}
	})

	t.Run("potions", func(t *testing.T) {
		potions := Items().Potions().All()
		if len(potions) == 0 {
			t.Fatal("potion view is empty")
		}
		for _, item := range potions {
			if item.Subtype != "potion" {
				t.Fatalf("potion view contains %q with subtype %q", item.Code, item.Subtype)
			}
		}
	})

	t.Run("foods", func(t *testing.T) {
		foods := Items().Foods().All()
		if len(foods) == 0 {
			t.Fatal("food view is empty")
		}
		for _, item := range foods {
			if item.Type != "consumable" || item.Subtype != "food" {
				t.Fatalf("food view contains %q with type %q and subtype %q", item.Code, item.Type, item.Subtype)
			}
		}
	})

	t.Run("tradeables", func(t *testing.T) {
		tradeables := Items().Tradeables().All()
		if len(tradeables) == 0 {
			t.Fatal("tradeable view is empty")
		}
		for _, item := range tradeables {
			if !item.Tradeable {
				t.Fatalf("tradeable view contains %q", item.Code)
			}
		}
	})
}
