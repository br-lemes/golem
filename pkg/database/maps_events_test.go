package database

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestGetEventPointsWithMapCodeDoesNotLoadEvents(t *testing.T) {
	loaded := false
	points := getEventPoints("slime", []string{"slime"}, nil, func() ([]schemas.ActiveEventSchema, error) {
		loaded = true
		return nil, nil
	})

	if loaded {
		t.Fatal("getEventPointsWith loaded events for a map code")
	}
	if len(points) != 0 {
		t.Fatalf("getEventPointsWith() = %#v, want empty", points)
	}
}

func TestGetEventPointsWithUnknownCodeDoesNotLoadEvents(t *testing.T) {
	loaded := false
	points := getEventPoints("unknown", nil, []string{"event"}, func() ([]schemas.ActiveEventSchema, error) {
		loaded = true
		return nil, nil
	})

	if loaded {
		t.Fatal("getEventPointsWith loaded events for an unknown code")
	}
	if len(points) != 0 {
		t.Fatalf("getEventPointsWith() = %#v, want empty", points)
	}
}

func TestGetEventPointsWithLoaderError(t *testing.T) {
	points := getEventPoints("event", nil, []string{"event"}, func() ([]schemas.ActiveEventSchema, error) {
		return nil, errors.New("events unavailable")
	})

	if len(points) != 0 {
		t.Fatalf("getEventPointsWith() = %#v, want empty", points)
	}
}

func TestGetEventPointsWithMatchingEvents(t *testing.T) {
	points := getEventPoints("event", nil, []string{"event"}, func() ([]schemas.ActiveEventSchema, error) {
		return []schemas.ActiveEventSchema{
			{
				Map: schemas.MapSchema{
					X:     10,
					Y:     20,
					Layer: "overworld",
					Interactions: schemas.InteractionSchema{
						Content: &schemas.MapContentSchema{Code: "event"},
					},
				},
			},
			{
				Map: schemas.MapSchema{X: 30, Y: 40, Layer: "underground"},
			},
		}, nil
	})

	want := EventPoints{{X: 10, Y: 20, Layer: "overworld"}: true}
	if len(points) != len(want) || !points[Point{
		X:     10,
		Y:     20,
		Layer: "overworld",
	}] {
		t.Fatalf("getEventPointsWith() = %#v, want %#v", points, want)
	}
}
