package task

import (
	"time"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Cooldown(character schemas.CharacterSchema) {
	if character.CooldownExpiration == nil {
		return
	}
	now := time.Now()
	if !character.CooldownExpiration.After(now) {
		return
	}
	duration := character.CooldownExpiration.Sub(now)
	console.Errorf("⏳ Cooldown active: %.0f seconds\n",
		duration.Round(time.Second).Seconds())
	time.Sleep(duration)
}
