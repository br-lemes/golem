package api

import (
	"time"

	"github.com/br-lemes/golem/pkg/console"
)

func handleCooldown(seconds int, reason string) {
	console.Errorf("⏳ Cooldown started: %d seconds (%s)\n", seconds, reason)
	time.Sleep(time.Duration(seconds) * time.Second)
}
