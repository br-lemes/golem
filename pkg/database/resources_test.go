package database

import "testing"

func TestResourcesSmoke(t *testing.T) {
	if len(Resources.All()) == 0 {
		t.Fatal("resources catalog is empty")
	}
}
