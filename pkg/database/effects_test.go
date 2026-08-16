package database

import "testing"

func TestEffectsSmoke(t *testing.T) {
	if len(Effects.All()) == 0 {
		t.Fatal("effects catalog is empty")
	}
}
