package cmd

import (
	"fmt"
	"os"
	"time"

	. "github.com/br-lemes/golem/pkg/schemas"
)

func handleCooldown(character CharacterSchema) {
	if character.CooldownExpiration == nil {
		return
	}
	now := time.Now()
	if !character.CooldownExpiration.After(now) {
		return
	}
	duration := character.CooldownExpiration.Sub(now)
	fmt.Fprintf(os.Stderr, "⏳ Cooldown active: %.0f seconds\n",
		duration.Round(time.Second).Seconds())
	time.Sleep(duration)
}
