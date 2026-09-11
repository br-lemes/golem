package utils

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestCooldowns(t *testing.T) {
	got, err := GetCooldown("/does-not-exist")
	if err == nil || got {
		t.Fatalf("GetCooldown(missing) = %t, %v", got, err)
	}
	got, err = GetCooldown("/my/{name}/action/move")
	if err != nil || !got {
		t.Fatalf("GetCooldown(move) = %t, %v", got, err)
	}
	cooldowns, err := GetCooldowns()
	if err != nil || !cooldowns["/my/{name}/action/move"] {
		t.Fatalf("GetCooldowns() = %#v, %v", cooldowns, err)
	}
}

func TestHasCooldown(t *testing.T) {
	if hasCooldown(&openapi3.PathItem{}) {
		t.Fatal("empty path unexpectedly has a cooldown")
	}
}
