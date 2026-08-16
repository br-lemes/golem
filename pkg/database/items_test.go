package database

import "testing"

func TestItemsSmoke(t *testing.T) {
	if len(Items.All()) == 0 {
		t.Fatal("items catalog is empty")
	}
}
