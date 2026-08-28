package parser

import "testing"

func TestItems(t *testing.T) {
	items, err := Items([]string{
		"iron_sword",
		"iron_ring@2",
		"small_health_potion@1=10",
		"iron_sword=2",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 4 || items[0].Code != "iron_sword" || items[1].Slot == nil || *items[1].Slot != "2" || items[2].Quantity == nil || *items[2].Quantity != 10 {
		t.Fatalf("unexpected parsed items: %#v", items)
	}
}

func TestItemsRejectsMalformedSyntax(t *testing.T) {
	inputs := []string{
		"",
		"@1",
		"iron_sword@",
		"iron_ring@abc",
		"iron_sword=",
		"iron_sword=abc",
		"iron_sword=0",
		"iron_sword=-1",
		"iron_sword=2@1",
		"definitely_unknown_item",
	}
	for _, input := range inputs {
		_, err := Items([]string{input})
		if err == nil {
			t.Errorf("expected %q to fail", input)
		}
	}
}

func TestSimpleItemSchemas(t *testing.T) {
	items, err := Items([]string{"iron_sword", "small_health_potion=10"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := items.SimpleItemSchemas()
	if err != nil || len(result) != 2 || result[0].Quantity != 1 || result[1].Quantity != 10 {
		t.Fatalf("unexpected simple items: %#v, %v", result, err)
	}
}

func TestEquipSchemasResolvesSlots(t *testing.T) {
	items, err := Items([]string{"iron_sword", "iron_ring", "iron_ring@2"})
	if err != nil {
		t.Fatal(err)
	}
	result, err := items.EquipSchemas()
	if err != nil {
		t.Fatal(err)
	}
	if result[0].Slot != "weapon" || result[1].Slot != "" || result[2].Slot != "ring2" {
		t.Fatalf("unexpected equipment slots: %#v", result)
	}
	if result[0].Quantity == nil || *result[0].Quantity != 1 {
		t.Fatalf("default quantity = %v, want 1", result[0].Quantity)
	}
}

func TestEquipSchemasRejectsInvalidSlotsAndItems(t *testing.T) {
	for _, input := range []string{"iron_sword@1", "iron_ring@3", "shell"} {
		items, err := Items([]string{input})
		if err != nil {
			t.Fatal(err)
		}
		_, err = items.EquipSchemas()
		if err == nil {
			t.Errorf("expected %q to fail", input)
		}
	}
}

func TestEquipSchemaRejectsUnknownManualItem(t *testing.T) {
	_, err := (Item{Code: "definitely_unknown_item"}).EquipSchema()
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestSimpleItemSchemaRejectsSlot(t *testing.T) {
	items, err := Items([]string{"iron_ring@1"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = items[0].SimpleItemSchema()
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestSimpleItemSchemasPropagatesItemError(t *testing.T) {
	items, err := Items([]string{"iron_sword", "iron_ring@1"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = items.SimpleItemSchemas()
	if err == nil {
		t.Fatal("expected an error")
	}
}
