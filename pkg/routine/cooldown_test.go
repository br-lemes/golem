package routine

import (
	"testing"
	"time"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCooldownReturnsWhenAlreadyExpired(t *testing.T) {
	expiration := time.Now().Add(-time.Second)
	character := schemas.CharacterSchema{CooldownExpiration: &expiration}

	start := time.Now()
	Cooldown(character)
	elapsed := time.Since(start)
	if elapsed > time.Second {
		t.Fatalf("Cooldown() waited %s for an expired cooldown", elapsed)
	}
}
