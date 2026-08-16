package database

import "testing"

func TestEquipmentsSmoke(t *testing.T) {
	if len(Equipments()) == 0 {
		t.Fatal("equipment list is empty")
	}
}
