package database

import "testing"

func TestViewSelectsAndSharesStoreItems(t *testing.T) {
	store := newStore(func() []testCatalogItem {
		return []testCatalogItem{{Code: "one"}, {Code: "two"}, {Code: "three"}}
	}, func(item *testCatalogItem) string {
		return item.Code
	})

	view := store.View(func(item *testCatalogItem) bool {
		return item.Code != "two"
	})

	items := view.All()
	if len(items) != 2 || items[0].Code != "one" || items[1].Code != "three" {
		t.Fatalf("view.All() = %v", items)
	}

	storeItem, exists := store.Get("three")
	if !exists {
		t.Fatal("store.Get did not find the view item")
	}
	viewItem, exists := view.Get("three")
	if !exists {
		t.Fatal("view.Get did not find a selected item")
	}
	if viewItem != storeItem {
		t.Fatal("view does not share pointers with the source store")
	}

	_, exists = view.Get("two")
	if exists {
		t.Fatal("view found an item excluded by its predicate")
	}
	got := view.Keys()
	if len(got) != 2 || got[0] != "one" || got[1] != "three" {
		t.Fatalf("view.Keys() = %v", got)
	}
	filtered := view.Filter(func(item *testCatalogItem) bool { return item.Code == "three" })
	if len(filtered) != 1 || filtered[0] != items[1] {
		t.Fatalf("view.Filter() = %v", filtered)
	}
	foundItem, exists := view.Find(func(item *testCatalogItem) bool { return item.Code == "one" })
	if !exists || foundItem != items[0] {
		t.Fatal("view.Find did not find the expected item")
	}
	_, exists = view.Find(func(item *testCatalogItem) bool { return item.Code == "two" })
	if exists {
		t.Fatal("view.Find found an excluded item")
	}
}
