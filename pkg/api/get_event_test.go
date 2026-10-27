package api

import (
	"net/http"
	"testing"
	"time"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestGetEventUsesCache(t *testing.T) {
	cache.CleanActiveEvents()
	expiration := time.Now().UTC().Add(time.Minute).Format(time.RFC3339)
	requests := 0
	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		requests++
		return []byte(`{"data":[{"code":"strange_apparition","expiration":"` + expiration + `","map":{"interactions":{"content":{"code":"strange_rocks"}}}}],"pages":1}`), nil
	}))

	event, err := client.GetEvent("strange_rocks")
	if err != nil {
		t.Fatal(err)
	}
	if event == nil {
		t.Fatal("expected event to be active")
	}
	event, err = client.GetEvent("strange_rocks")
	if err != nil {
		t.Fatal(err)
	}
	if event == nil {
		t.Fatal("expected cached event to be active")
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
}

func TestGetEventRefreshesExpiredCache(t *testing.T) {
	cache.CleanActiveEvents()
	expired := time.Now().UTC().Add(-time.Minute)
	cache.SaveActiveEvents([]schemas.ActiveEventSchema{{
		Expiration: expired,
		Map: schemas.MapSchema{
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "strange_rocks"},
			},
		},
	}})
	requests := 0
	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		requests++
		return []byte(`{"data":[],"pages":1}`), nil
	}))

	event, err := client.GetEvent("strange_rocks")
	if err != nil {
		t.Fatal(err)
	}
	if event != nil {
		t.Fatal("expected expired event to be inactive")
	}
	if requests != 1 {
		t.Fatalf("requests = %d, want 1", requests)
	}
}
