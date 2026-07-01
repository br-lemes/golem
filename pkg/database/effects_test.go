package database

import "testing"

func TestEffectsSmoke(t *testing.T) {
	if len(GetEffects()) == 0 {
		t.Error("effects list is empty")
	}
	if len(GetEffectCodes()) == 0 {
		t.Error("effect codes are empty")
	}
}
