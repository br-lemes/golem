package cache

import (
	"testing"
	"time"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestEventCache(t *testing.T) {
	initializeTestCache(t)

	expiration := time.Now().UTC().Add(time.Minute)
	SaveActiveEvents([]schemas.ActiveEventSchema{{
		Code:       "strange_apparition",
		Expiration: expiration,
		Map: schemas.MapSchema{
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "strange_rocks"},
			},
		},
	}})

	found, refresh := HasEvent("strange_rocks")
	if !found || refresh {
		t.Fatal("active event was not saved")
	}
	found, refresh = HasEvent("missing")
	if found || refresh {
		t.Fatal("missing event did not use the valid cache")
	}
}

func TestEventCacheMissRequiresRefresh(t *testing.T) {
	initializeTestCache(t)

	SaveActiveEvents(nil)
	found, refresh := HasEvent("missing")
	if found || refresh {
		t.Fatal("empty cache did not require refresh")
	}
}
