package database

import "testing"

func TestItemsSmoke(t *testing.T) {
	if len(GetItems()) == 0 {
		t.Error("items list is empty")
	}
	if len(GetItemCodes()) == 0 {
		t.Error("item codes are empty")
	}
	if len(GetItemTypes()) == 0 {
		t.Error("item types are empty")
	}
}
