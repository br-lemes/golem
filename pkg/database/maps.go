package database

import (
	_ "embed"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

type (
	Point struct {
		X     int
		Y     int
		Layer string
	}
)

var (
	//go:embed maps.json
	maps []byte
)

var Maps = newStore(jsonLoader[schemas.MapSchema](maps), func(tile *schemas.MapSchema) Point {
	return Point{X: tile.X, Y: tile.Y, Layer: tile.Layer}
})

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
