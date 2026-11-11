package datatable

import (
	"net/http"
	"testing"
)

type dataTableTestItem struct {
	Enabled bool   `json:"enabled"`
	Level   int    `json:"level"`
	Name    string `json:"name"`
}

func TestSortBool(t *testing.T) {
	items := []dataTableTestItem{{Enabled: true}, {Enabled: false}}
	columns := []Column[dataTableTestItem]{{Field: "enabled", Sortable: true}}

	Sort(items, "enabled", "asc", columns)
	if items[0].Enabled || !items[1].Enabled {
		t.Fatalf("ascending bool sort = %#v", items)
	}

	Sort(items, "enabled", "desc", columns)
	if !items[0].Enabled || items[1].Enabled {
		t.Fatalf("descending bool sort = %#v", items)
	}
}

func TestSort(t *testing.T) {
	items := []dataTableTestItem{{Level: 20, Name: "B"}, {Level: 5, Name: "A"}}
	columns := []Column[dataTableTestItem]{
		{Field: "level", Sortable: true},
		{Field: "name", Sortable: true},
	}

	Sort(items, "level", "asc", columns)
	if items[0].Level != 5 || items[1].Level != 20 {
		t.Fatalf("ascending sort = %#v", items)
	}

	Sort(items, "name", "desc", columns)
	if items[0].Name != "B" || items[1].Name != "A" {
		t.Fatalf("descending sort = %#v", items)
	}
}

func TestSearch(t *testing.T) {
	items := []dataTableTestItem{
		{Level: 20, Name: "Bronze"},
		{Level: 5, Name: "Silver"},
	}
	columns := []Column[dataTableTestItem]{{Field: "level"}, {Field: "name"}}

	results := Search(items, "bro", columns)
	if len(results) != 1 || results[0].Name != "Bronze" {
		t.Fatalf("string search = %#v", results)
	}

	results = Search(items, "20", columns)
	if len(results) != 1 || results[0].Level != 20 {
		t.Fatalf("numeric search = %#v", results)
	}
}

func TestRefreshPeriodFromRequest(t *testing.T) {
	tests := []struct {
		cookie   string
		expected int
	}{
		{"", DefaultRefreshPeriod},
		{"off", DefaultRefreshPeriod},
		{"0", 0},
		{"5", 5},
		{"90", 90},
		{"7", DefaultRefreshPeriod},
		{" 10 ", 10},
	}
	for _, test := range tests {
		request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
		if err != nil {
			t.Fatal(err)
		}
		if test.cookie != "" {
			cookie := &http.Cookie{
				Name:  "swr-refresh-period",
				Value: test.cookie,
			}
			request.AddCookie(cookie)
		}
		actual := RefreshPeriodFromRequest(request)
		if actual != test.expected {
			t.Fatalf("RefreshPeriodFromRequest(%q) = %d, want %d", test.cookie, actual, test.expected)
		}
	}
}
