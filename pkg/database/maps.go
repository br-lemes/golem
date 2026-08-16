package database

import (
	_ "embed"
	"slices"
	"sync"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

type (
	Point struct {
		X     int
		Y     int
		Layer schemas.MapLayer
	}

	SearchNode struct {
		Point       Point
		Transitions []schemas.MapSchema
	}

	SearchResult struct {
		Target      schemas.MapSchema
		Transitions []schemas.MapSchema
	}

	EventPoints map[Point]bool
)

var (
	//go:embed maps.json
	maps []byte
)

var Maps = newStore(jsonLoader[schemas.MapSchema](maps), func(tile *schemas.MapSchema) Point {
	return Point{X: tile.X, Y: tile.Y, Layer: tile.Layer}
})

func FindClosest(character schemas.CharacterSchema, code string) *SearchResult {
	//+gocover:ignore:block production wrapper over tested logic
	eventPoints := getEventPoints(code, MapCodes(), EventContentCodes(), api.EventsActive)
	return findClosest(navigationContext{
		Maps:             Maps.Get,
		Character:        character,
		Code:             code,
		EventPoints:      eventPoints,
		LoadAchievements: api.AccountsAchievements,
	})
}

func getEventPoints(code string, mapCodes, eventCodes []string, load func() ([]schemas.ActiveEventSchema, error)) EventPoints {
	if slices.Contains(mapCodes, code) {
		return EventPoints{}
	}

	if !slices.Contains(eventCodes, code) {
		return EventPoints{}
	}

	events, err := load()
	if err != nil {
		return EventPoints{}
	}
	result := EventPoints{}
	for _, event := range events {
		if event.Map.Interactions.Content != nil && event.Map.Interactions.Content.Code == code {
			eventPoint := Point{
				X:     event.Map.X,
				Y:     event.Map.Y,
				Layer: event.Map.Layer,
			}
			result[eventPoint] = true
		}
	}
	return result
}

var mapCodes = sync.OnceValue(func() []string {
	seen := make(map[string]struct{})
	var codes []string

	for _, tile := range Maps.All() {
		if tile.Interactions.Content == nil {
			continue
		}
		code := tile.Interactions.Content.Code
		if code == "" {
			//+gocover:ignore:block map contents always have codes
			continue
		}
		_, exists := seen[code]
		if exists {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
	}
	return codes
})

func MapCodes() []string {
	return mapCodes()
}
