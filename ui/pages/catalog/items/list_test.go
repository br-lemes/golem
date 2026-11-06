package items

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/ui/components/datatable"
)

func TestPaginateItems(t *testing.T) {
	size := 10
	items := make([]*schemas.ItemSchema, size*2+1)

	tests := []struct {
		name      string
		page      int
		wantPage  int
		wantPages int
		wantItems int
	}{
		{
			name:      "first page",
			page:      1,
			wantPage:  1,
			wantPages: 3,
			wantItems: size,
		},
		{
			name:      "last page",
			page:      3,
			wantPage:  3,
			wantPages: 3,
			wantItems: 1,
		},
		{
			name:      "page below one",
			page:      0,
			wantPage:  1,
			wantPages: 3,
			wantItems: size,
		},
		{
			name:      "page above total",
			page:      4,
			wantPage:  3,
			wantPages: 3,
			wantItems: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pageItems, page := datatable.Paginate(items, test.page, size)
			totalPages := max(1, (len(items)+size-1)/size)
			if page != test.wantPage {
				t.Errorf("page = %d, want %d", page, test.wantPage)
			}
			if totalPages != test.wantPages {
				t.Errorf("total pages = %d, want %d", totalPages, test.wantPages)
			}
			if len(pageItems) != test.wantItems {
				t.Errorf("page items = %d, want %d", len(pageItems), test.wantItems)
			}
		})
	}
}

func TestPaginateItemsEmpty(t *testing.T) {
	pageItems, page := datatable.Paginate([]*schemas.ItemSchema(nil), 1, 10)
	totalPages := 1
	if page != 1 {
		t.Errorf("page = %d, want 1", page)
	}
	if totalPages != 1 {
		t.Errorf("total pages = %d, want 1", totalPages)
	}
	if len(pageItems) != 0 {
		t.Errorf("page items = %d, want 0", len(pageItems))
	}
}
