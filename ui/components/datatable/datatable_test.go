package datatable

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
)

type dataTableTestItem struct {
	Level int    `json:"level"`
	Name  string `json:"name"`
}

func renderDataTable(t *testing.T, component templ.Component) *goquery.Document {
	t.Helper()
	var output bytes.Buffer
	err := component.Render(context.Background(), &output)
	if err != nil {
		t.Fatalf("render data table: %v", err)
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(output.String()))
	if err != nil {
		t.Fatalf("parse data table: %v", err)
	}
	return document
}

func assertDataTablePanics(t *testing.T, callback func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected data table to panic")
		}
	}()
	callback()
}

func TestDataTableRendersEmptyColumn(t *testing.T) {
	items := []dataTableTestItem{{Name: "item"}}
	columns := []Column[dataTableTestItem]{
		{Header: "Actions", Class: "actions"},
	}
	component := DataTable(DataTableProps[dataTableTestItem]{
		Data:    items,
		Columns: columns,
	})

	document := renderDataTable(t, component)
	cell := document.Find("tbody tr").First().Find("td.actions")
	if cell.Length() != 1 {
		t.Fatalf("expected one empty actions cell, found %d", cell.Length())
	}
	text := strings.TrimSpace(cell.Text())
	if text != "" {
		t.Fatalf("expected empty actions cell, got %q", text)
	}
}

func TestDataTableRendersFields(t *testing.T) {
	document := renderDataTable(t, DataTable(DataTableProps[dataTableTestItem]{
		Data: []dataTableTestItem{{Name: "item", Level: 3}},
		Columns: []Column[dataTableTestItem]{
			{Header: "Name", Field: "name"},
			{Header: "Level", Field: "level"},
		},
	}))

	headers := document.Find("thead th")
	if headers.Length() != 2 {
		t.Fatalf("expected two table headers, found %d", headers.Length())
	}
	values := []string{}
	document.Find("tbody tr").First().Find("td").Each(func(_ int, cell *goquery.Selection) {
		values = append(values, strings.TrimSpace(cell.Text()))
	})
	want := []string{"item", "3"}
	if !slices.Equal(values, want) {
		t.Fatalf("expected row values %v, got %v", want, values)
	}
}

func TestDataTableRendersPointerItems(t *testing.T) {
	item := &dataTableTestItem{Name: "item", Level: 3}
	document := renderDataTable(t, DataTable(DataTableProps[*dataTableTestItem]{
		Data:    []*dataTableTestItem{item, nil},
		Columns: []Column[*dataTableTestItem]{{Header: "Name", Field: "name"}},
	}))

	rows := document.Find("tbody tr")
	text := strings.TrimSpace(rows.Eq(0).Find("td").First().Text())
	if text != "item" {
		t.Fatalf("expected pointer item to be rendered, got %q", text)
	}
	text = strings.TrimSpace(rows.Eq(1).Find("td").First().Text())
	if text != "" {
		t.Fatalf("expected nil pointer item to render empty, got %q", text)
	}
}

func TestDataTableRendersCustomCell(t *testing.T) {
	document := renderDataTable(t, DataTable(DataTableProps[dataTableTestItem]{
		Data: []dataTableTestItem{{Name: "item"}},
		Columns: []Column[dataTableTestItem]{{
			Header: "Action",
			Cell: func(item dataTableTestItem) templ.Component {
				return templ.Raw("<strong>" + item.Name + " action</strong>")
			},
		}},
	}))

	text := document.Find("tbody tr td strong").Text()
	if text != "item action" {
		t.Fatalf("expected custom cell to be rendered, got %q", text)
	}
}

func TestDataTableRendersRowHref(t *testing.T) {
	document := renderDataTable(t, DataTable(DataTableProps[dataTableTestItem]{
		Data:    []dataTableTestItem{{Name: "item"}},
		Columns: []Column[dataTableTestItem]{{Header: "Name", Field: "name"}},
		RowHref: func(item dataTableTestItem) string {
			return "/items/" + item.Name
		},
	}))

	row := document.Find("tbody tr").First()
	if row.Length() != 1 {
		t.Fatalf("expected one row, found %d", row.Length())
	}
	actual, exists := row.Attr("data-swr")
	if !exists {
		t.Error("expected row data-swr attribute")
	} else if actual != "/items/item" {
		t.Errorf("expected row data-swr attribute to be %q, got %q", "/items/item", actual)
	}
}

func TestDataTableDoesNotRenderRowAttributesWithoutHref(t *testing.T) {
	document := renderDataTable(t, DataTable(DataTableProps[dataTableTestItem]{
		Data:    []dataTableTestItem{{Name: "item"}},
		Columns: []Column[dataTableTestItem]{{Header: "Name", Field: "name"}},
	}))

	row := document.Find("tbody tr").First()
	for _, attribute := range []string{"data-swr", "data-swr-href"} {
		_, exists := row.Attr(attribute)
		if exists {
			t.Errorf("did not expect row attribute %q", attribute)
		}
	}
}

func TestDataTableRendersEmptyMessage(t *testing.T) {
	columns := []Column[dataTableTestItem]{{Header: "Name", Field: "name"}}
	document := renderDataTable(t, DataTable(DataTableProps[dataTableTestItem]{
		Columns:      columns,
		EmptyMessage: "Nothing found",
	}))

	text := document.Find("tbody tr td").Text()
	if text != "Nothing found" {
		t.Fatalf("expected custom empty message, got %q", text)
	}
}

func TestDataTableRendersDefaultEmptyMessage(t *testing.T) {
	document := renderDataTable(t, DataTable(DataTableProps[dataTableTestItem]{
		Columns: []Column[dataTableTestItem]{{Header: "Name", Field: "name"}},
	}))

	text := document.Find("tbody tr td").Text()
	if text != "No data available" {
		t.Fatalf("expected default empty message, got %q", text)
	}
}

func TestDataTablePanicsForNonStructItems(t *testing.T) {
	assertDataTablePanics(t, func() {
		DataTable(DataTableProps[int]{
			Data:    []int{1},
			Columns: []Column[int]{{Header: "Value", Field: "value"}},
		})
	})
}

func TestDataTablePanicsForUnknownField(t *testing.T) {
	assertDataTablePanics(t, func() {
		DataTable(DataTableProps[dataTableTestItem]{
			Data: []dataTableTestItem{{Name: "item"}},
			Columns: []Column[dataTableTestItem]{{
				Header: "Missing",
				Field:  "missing",
			}},
		})
	})
}
