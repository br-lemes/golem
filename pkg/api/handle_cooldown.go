package api

import (
	"time"

	"github.com/br-lemes/golem/pkg/console"
)

func handleCooldown(seconds int) {
	console.Errorf("⏳ Cooldown started: %d seconds\n", seconds)
	time.Sleep(time.Duration(seconds) * time.Second)
}
