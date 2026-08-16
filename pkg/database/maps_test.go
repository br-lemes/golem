package database

import "testing"

func TestMapsSmoke(t *testing.T) {
	if len(GetMaps()) == 0 {
		t.Fatal("maps catalog is empty")
	}
}
