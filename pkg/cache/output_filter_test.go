package cache

import "testing"

func TestOutputFiltersCache(t *testing.T) {
	initializeTestCache(t)
	err := AddOutputFilter("fight", "field", "hp")
	if err != nil {
		t.Fatal(err)
	}
	err = AddOutputFilter("fight", "field", "xp")
	if err != nil {
		t.Fatal(err)
	}
	filters := GetOutputFilters("fight")
	if len(filters["field"]) != 2 {
		t.Fatalf("GetOutputFilters() = %#v", filters)
	}
	err = EditOutputFilter("fight", "field", "xp", "gold")
	if err != nil {
		t.Fatalf("EditOutputFilter() = %v", err)
	}
	filters = GetOutputFilters("fight")
	if len(filters["field"]) != 2 || filters["field"][0] != "gold" {
		t.Fatalf("edited filters = %#v", filters)
	}
}

func TestOutputFilterCommandsAndRemoval(t *testing.T) {
	initializeTestCache(t)
	err := AddOutputFilter("fight", "field", "hp")
	if err != nil {
		t.Fatal(err)
	}
	err = AddOutputFilter("gather", "field", "xp")
	if err != nil {
		t.Fatal(err)
	}
	all := ListOutputFilters("")
	if len(all) != 2 {
		t.Fatalf("ListOutputFilters() = %#v", all)
	}
	removed, err := RemoveOutputFilter("fight", "field", "hp")
	if err != nil || !removed {
		t.Fatalf("RemoveOutputFilter(existing) = %t, %v", removed, err)
	}
	removed, err = RemoveOutputFilter("fight", "field", "missing")
	if err != nil || removed {
		t.Fatalf("RemoveOutputFilter(missing) = %t, %v", removed, err)
	}
}

func TestEditOutputFilterReportsMissingFilter(t *testing.T) {
	initializeTestCache(t)
	err := EditOutputFilter("fight", "field", "missing", "new")
	if err == nil || err.Error() != "output filter not found" {
		t.Fatalf("EditOutputFilter(missing) = %v", err)
	}
}
