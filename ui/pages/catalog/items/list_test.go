package items

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
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
			pageItems, page := paginateItems(items, test.page, size)
			totalPages := totalPagesForItems(len(items), size)
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
	pageItems, page := paginateItems(nil, 1, 10)
	totalPages := totalPagesForItems(0, 10)
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
