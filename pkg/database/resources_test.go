package database

import "testing"

func TestResourcesSmoke(t *testing.T) {
	if len(GetResources()) == 0 {
		t.Error("resources list is empty")
	}
}
