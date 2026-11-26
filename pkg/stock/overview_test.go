package stock

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestStockItemQuantitiesIncludesEquippedItems(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{{Code: "ash_wood", Quantity: 2}}
	characters := []schemas.CharacterSchema{{
		Inventory:            &inventory,
		WeaponSlot:           "iron_sword",
		Utility1Slot:         "health_potion",
		Utility1SlotQuantity: 3,
	}}
	bankItems := []schemas.SimpleItemSchema{{Code: "ash_wood", Quantity: 5}}

	quantities := itemQuantities(bankItems, characters)
	if quantities["ash_wood"] != 7 || quantities["iron_sword"] != 1 || quantities["health_potion"] != 3 {
		t.Fatalf("stockItemQuantities() = %#v", quantities)
	}
}

func TestStockGoldQuantityIncludesCharacters(t *testing.T) {
	characters := []schemas.CharacterSchema{{Gold: 20}, {Gold: 30}}
	quantity := goldQuantity(50, characters)
	if quantity != 100 {
		t.Fatalf("stockGoldQuantity() = %d, want 100", quantity)
	}
}

func TestRequirementsSafetyZeroHasNoDirectRequirement(t *testing.T) {
	result := requirements(map[string]models.Stock{
		"test_wood": {
			Code:    "test_wood",
			Mode:    models.StockModeSafety,
			Enabled: true,
		},
	}, map[string]int{"test_wood": 10}, nil, false)
	if len(result) != 0 {
		t.Fatalf("requirements() = %#v, want no direct requirement", result)
	}
}

func TestRequirements(t *testing.T) {
	stocks := map[string]models.Stock{
		"test_sword": {
			Code:     "test_sword",
			Quantity: 5,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		"test_wood": {
			Code:         "test_wood",
			Quantity:     100,
			Mode:         models.StockModeSafety,
			Enabled:      true,
			Replenishing: true,
		},
		"future_item": {
			Code:     "future_item",
			Quantity: 100,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
		"empty_item": {
			Code:     "empty_item",
			Quantity: 100,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
	}
	quantities := map[string]int{
		"test_sword":  2,
		"test_wood":   90,
		"future_item": 200,
	}

	result := requirements(stocks, quantities, nil, false)
	if len(result) != 3 {
		t.Fatalf("requirements() = %#v", result)
	}
	if result["test_sword"].Needed != 3 || result["test_sword"].Target != 5 {
		t.Fatalf("exact stock = %#v", result["test_sword"])
	}
	if result["test_wood"].Needed != 43 || result["test_wood"].Safety != 100 || result["test_wood"].Target != 133 {
		t.Fatalf("safety stock = %#v", result["test_wood"])
	}
	if result["test_wood"].Demand != 0 {
		t.Fatalf("direct required = %#v", result["test_wood"])
	}
	_, ok := result["future_item"]
	if ok {
		t.Fatalf("replenished stock = %#v", result)
	}
	if result["empty_item"].Needed != 100 || result["empty_item"].Target != 133 {
		t.Fatalf("empty stock = %#v", result["empty_item"])
	}
}

func TestRequirementsTarget(t *testing.T) {
	stocks := map[string]models.Stock{
		"ash_wood": {
			Code:     "ash_wood",
			Quantity: 100,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
	}
	quantities := map[string]int{"ash_wood": 110}

	result := requirements(stocks, quantities, nil, false)
	if len(result) != 0 {
		t.Fatalf("default requirements() = %#v", result)
	}
	stocks["ash_wood"] = models.Stock{
		Code:         "ash_wood",
		Quantity:     100,
		Mode:         models.StockModeSafety,
		Enabled:      true,
		Replenishing: true,
	}
	result = requirements(stocks, quantities, nil, false)
	if result["ash_wood"].Needed != 23 {
		t.Fatalf("target stock = %#v", result)
	}
}

func TestRequirementsFiltersOutputAfterExpansion(t *testing.T) {
	stocks := map[string]models.Stock{
		"gold_bar": {
			Code:     "gold_bar",
			Quantity: 1,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		"gold_ore": {
			Code:     "gold_ore",
			Quantity: 0,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
	}
	result := requirements(stocks, nil, []string{"gold_ore"}, false)
	if len(result) != 1 || result["gold_ore"].Needed != 10 || result["gold_ore"].Demand != 10 {
		t.Fatalf("requirements() = %#v", result)
	}
}

func TestRecipeRequirementsReservesDirectTargets(t *testing.T) {
	craftItems := []schemas.SimpleItemSchema{{Code: "ash_wood", Quantity: 10}}
	craftQuantity := 1
	items := map[string]*schemas.ItemSchema{
		"ash_plank": {
			Code: "ash_plank",
			Craft: &schemas.CraftSchema{
				Items:    &craftItems,
				Quantity: &craftQuantity,
			},
		},
		"ash_wood": {Code: "ash_wood"},
	}
	direct := map[string]goal{
		"ash_plank": {Quantity: 10, Expand: true},
		"ash_wood":  {Quantity: 300, Expand: true},
	}
	quantities := map[string]int{"ash_wood": 200}

	required, needed := recipeRequirements(items, nil, direct, quantities)
	if required["ash_wood"] != 400 || needed["ash_wood"] != 200 || needed["ash_plank"] != 10 {
		t.Fatalf("recipeRequirements() = required %#v, needed %#v", required, needed)
	}

	direct["ash_wood"] = goal{Quantity: 300}
	quantities["ash_wood"] = 400
	required, needed = recipeRequirements(items, nil, direct, quantities)
	if required["ash_wood"] != 400 || needed["ash_wood"] != 0 || needed["ash_plank"] != 10 {
		t.Fatalf("reserved recipeRequirements() = required %#v, needed %#v", required, needed)
	}
}

func TestRecipeRequirementsExpandsNPCOffer(t *testing.T) {
	price := 220
	offer := schemas.SimpleNPCItemSchema{
		BuyPrice: &price,
		Code:     "corrupted_crown",
		Currency: "corrupted_gem",
	}
	offers := map[string][]schemas.SimpleNPCItemSchema{
		"corrupted_crown": {offer},
	}
	direct := map[string]goal{"corrupted_crown": {Quantity: 5, Expand: true}}

	required, needed := recipeRequirements(nil, offers, direct, nil)
	if required["corrupted_gem"] != 1100 || needed["corrupted_gem"] != 1100 {
		t.Fatalf("recipeRequirements() = required %#v, needed %#v", required, needed)
	}
}

func TestRecipeRequirementsExpandsGoldOnlyForNPCExclusiveItems(t *testing.T) {
	price := 30000
	offers := map[string][]schemas.SimpleNPCItemSchema{
		"burn_rune": {{BuyPrice: &price, Code: "burn_rune", Currency: "gold"}},
		"gold_ore":  {{BuyPrice: &price, Code: "gold_ore", Currency: "gold"}},
	}
	direct := map[string]goal{
		"burn_rune": {Quantity: 1, Expand: true},
		"gold_ore":  {Quantity: 1, Expand: true},
	}
	quantities := map[string]int{"gold": 5000}

	required, needed := recipeRequirements(catalogItems(), offers, direct, quantities)
	if required["gold"] != 30000 || needed["gold"] != 25000 {
		t.Fatalf("gold requirements = required %#v, needed %#v", required, needed)
	}
}
