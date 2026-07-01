package database

import "testing"

func TestMapsSmoke(t *testing.T) {
	if len(GetMaps()) == 0 {
		t.Error("maps list is empty")
	}
	if len(GetMapCodes()) == 0 {
		t.Error("map codes are empty")
	}
}
