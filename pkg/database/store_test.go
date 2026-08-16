package database

import "testing"

type testCatalogItem struct {
	Code string
}

func TestStoreIndexesItemsWithoutCopying(t *testing.T) {
	store := newStore(func() []testCatalogItem {
		return []testCatalogItem{{Code: "one"}, {Code: "two"}}
	}, func(item *testCatalogItem) string {
		return item.Code
	})

	items := store.All()
	item, exists := store.Get("two")
	if !exists {
		t.Fatal("store.Get did not find an existing item")
	}
	if item != items[1] {
		t.Fatal("store.Get returned a different pointer than store.All")
	}
	got := store.Keys()
	if len(got) != 2 || got[0] != "one" || got[1] != "two" {
		t.Fatalf("store.Keys() = %v", got)
	}
	filtered := store.Filter(func(item *testCatalogItem) bool { return item.Code == "two" })
	if len(filtered) != 1 || filtered[0] != items[1] {
		t.Fatalf("store.Filter() = %v", filtered)
	}
	foundItem, exists := store.Find(func(item *testCatalogItem) bool { return item.Code == "one" })
	if !exists || foundItem != items[0] {
		t.Fatal("store.Find did not find the expected item")
	}
	_, exists = store.Find(func(item *testCatalogItem) bool { return item.Code == "missing" })
	if exists {
		t.Fatal("store.Find found a missing item")
	}
}

func TestStoreRejectsDuplicateKeys(t *testing.T) {
	store := newStore(func() []testCatalogItem {
		return []testCatalogItem{{Code: "duplicate"}, {Code: "duplicate"}}
	}, func(item *testCatalogItem) string {
		return item.Code
	})

	defer func() {
		if recover() == nil {
			t.Fatal("store did not panic on a duplicate key")
		}
	}()

	store.All()
}

func TestJSONLoaderRejectsInvalidJSON(t *testing.T) {
	loader := jsonLoader[testCatalogItem]([]byte("{"))

	defer func() {
		if recover() == nil {
			t.Fatal("jsonLoader did not panic on invalid JSON")
		}
	}()

	loader()
}
