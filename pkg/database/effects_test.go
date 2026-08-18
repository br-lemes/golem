package database

import "testing"

func TestEffectsSmoke(t *testing.T) {
	if len(Effects().All()) == 0 {
		t.Fatal("effects catalog is empty")
	}
}

func TestEffectsEquipmentsView(t *testing.T) {
	for _, effect := range Effects().Equipments().All() {
		if effect.Type != "equipment" {
			t.Fatalf("Effects().Equipments() returned effect with type %q", effect.Type)
		}
	}
}
