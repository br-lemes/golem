package cmd

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

	quantities := stockItemQuantities(bankItems, characters)
	if quantities["ash_wood"] != 7 || quantities["iron_sword"] != 1 || quantities["health_potion"] != 3 {
		t.Fatalf("stockItemQuantities() = %#v", quantities)
	}
}

func TestStockGoldQuantityIncludesCharacters(t *testing.T) {
	characters := []schemas.CharacterSchema{{Gold: 20}, {Gold: 30}}
	quantity := stockGoldQuantity(50, characters)
	if quantity != 100 {
		t.Fatalf("stockGoldQuantity() = %d, want 100", quantity)
	}
}

func TestStockRequirements(t *testing.T) {
	stocks := map[string]models.Stock{
		"test_sword": {
			Code:     "test_sword",
			Quantity: 5,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		"test_wood": {
			Code:     "test_wood",
			Quantity: 100,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
		"future_item": {
			Code:     "future_item",
			Quantity: 100,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
		"hidden_item": {
			Code:     "hidden_item",
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
	taskCodes := map[string]bool{
		"test_wood":   true,
		"future_item": true,
		"hidden_item": true,
	}
	reachable := map[string]int{"test_wood": 100}

	result := stockRequirements(stocks, quantities, nil, taskCodes, reachable, nil, false)
	if len(result) != 3 {
		t.Fatalf("stockRequirements() = %#v", result)
	}
	if result["test_sword"].Needed != 3 || result["test_sword"].Quantity == nil || *result["test_sword"].Quantity != 5 {
		t.Fatalf("exact stock = %#v", result["test_sword"])
	}
	if result["test_wood"].Needed != 43 || result["test_wood"].Safety != 100 || result["test_wood"].Target != 133 {
		t.Fatalf("safety stock = %#v", result["test_wood"])
	}
	if result["test_wood"].Required != 0 {
		t.Fatalf("direct required = %#v", result["test_wood"])
	}
	if result["future_item"].Needed != 0 || result["future_item"].Target != 133 {
		t.Fatalf("future stock = %#v", result["future_item"])
	}
	_, ok := result["hidden_item"]
	if ok {
		t.Fatalf("hidden task result = %#v", result)
	}
}

func TestStockRequirementsTarget(t *testing.T) {
	stocks := map[string]models.Stock{
		"ash_wood": {
			Code:     "ash_wood",
			Quantity: 100,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
	}
	quantities := map[string]int{"ash_wood": 110}
	taskCodes := map[string]bool{"ash_wood": true}
	reachable := map[string]int{"ash_wood": 100}

	result := stockRequirements(stocks, quantities, nil, taskCodes, reachable, nil, false)
	if len(result) != 0 {
		t.Fatalf("default stockRequirements() = %#v", result)
	}
	result = stockRequirements(stocks, quantities, nil, taskCodes, reachable, nil, true)
	if result["ash_wood"].Needed != 23 {
		t.Fatalf("target stock = %#v", result)
	}
}

func TestStockRequirementsFiltersOutputAfterExpansion(t *testing.T) {
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
	characters := []schemas.CharacterSchema{{Level: 100, MiningLevel: 100}}

	result := stockRequirements(stocks, nil, characters, nil, nil, []string{"gold_ore"}, false)
	if len(result) != 1 || result["gold_ore"].Needed != 10 || result["gold_ore"].Required != 10 {
		t.Fatalf("stockRequirements() = %#v", result)
	}
}

func TestStockItemUsable(t *testing.T) {
	conditions := []schemas.ConditionSchema{{
		Code:     "level",
		Operator: "gt",
		Value:    4,
	}}
	item := schemas.ItemSchema{Conditions: &conditions}
	characters := []schemas.CharacterSchema{{Level: 4}}
	if stockItemUsable(item, characters) {
		t.Fatal("stockItemUsable() = true for an unavailable item")
	}
	characters[0].Level = 5
	if !stockItemUsable(item, characters) {
		t.Fatal("stockItemUsable() = false for an available item")
	}
}

func TestStockRecipeRequirementsReservesDirectTargets(t *testing.T) {
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
	direct := map[string]stockGoal{
		"ash_plank": {Quantity: 10, Expand: true},
		"ash_wood":  {Quantity: 300, Expand: true},
	}
	quantities := map[string]int{"ash_wood": 200}
	characters := []schemas.CharacterSchema{{}}

	required, needed := stockRecipeRequirements(items, nil, direct, quantities, characters)
	if required["ash_wood"] != 400 || needed["ash_wood"] != 200 || needed["ash_plank"] != 10 {
		t.Fatalf("stockRecipeRequirements() = required %#v, needed %#v", required, needed)
	}

	direct["ash_wood"] = stockGoal{Quantity: 300}
	quantities["ash_wood"] = 400
	required, needed = stockRecipeRequirements(items, nil, direct, quantities, characters)
	if required["ash_wood"] != 400 || needed["ash_wood"] != 0 || needed["ash_plank"] != 10 {
		t.Fatalf("reserved stockRecipeRequirements() = required %#v, needed %#v", required, needed)
	}
}

func TestStockRecipeRequirementsExpandsNPCOffer(t *testing.T) {
	price := 220
	offer := schemas.SimpleNPCItemSchema{
		BuyPrice: &price,
		Code:     "corrupted_crown",
		Currency: "corrupted_gem",
	}
	offers := map[string][]schemas.SimpleNPCItemSchema{
		"corrupted_crown": {offer},
	}
	direct := map[string]stockGoal{
		"corrupted_crown": {Quantity: 5, Expand: true},
	}

	required, needed := stockRecipeRequirements(nil, offers, direct, nil, nil)
	if required["corrupted_gem"] != 1100 || needed["corrupted_gem"] != 1100 {
		t.Fatalf("stockRecipeRequirements() = required %#v, needed %#v", required, needed)
	}
}

func TestStockRecipeRequirementsExpandsGoldOnlyForNPCExclusiveItems(t *testing.T) {
	price := 30000
	offers := map[string][]schemas.SimpleNPCItemSchema{
		"burn_rune": {{BuyPrice: &price, Code: "burn_rune", Currency: "gold"}},
		"gold_ore":  {{BuyPrice: &price, Code: "gold_ore", Currency: "gold"}},
	}
	direct := map[string]stockGoal{
		"burn_rune": {Quantity: 1, Expand: true},
		"gold_ore":  {Quantity: 1, Expand: true},
	}
	quantities := map[string]int{"gold": 5000}

	required, needed := stockRecipeRequirements(stockCatalogItems(), offers, direct, quantities, nil)
	if required["gold"] != 30000 || needed["gold"] != 25000 {
		t.Fatalf("gold requirements = required %#v, needed %#v", required, needed)
	}
}
