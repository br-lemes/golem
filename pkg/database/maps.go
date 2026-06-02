package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/br-lemes/golem/pkg/schemas"
)

//go:embed maps.json
var maps []byte

type Point struct {
	X     int
	Y     int
	Layer MapLayer
}

var tilesCache map[Point]MapSchema

func FindClosest(character CharacterSchema, code string) *MapSchema {
	initTilesCache()

	startPoint := Point{X: character.X, Y: character.Y, Layer: character.Layer}
	queue := []Point{startPoint}
	visited := make(map[Point]bool)
	visited[startPoint] = true

	var foundTile *MapSchema
	dx := []int{0, 0, 1, -1}
	dy := []int{1, -1, 0, 0}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		currentTile, exists := tilesCache[current]
		if exists {
			if currentTile.Interactions.Content != nil {
				if currentTile.Interactions.Content.Code == code {
					foundTile = &currentTile
					break
				}
			}
		}

		for i := 0; i < 4; i++ {
			nextPoint := Point{X: current.X + dx[i], Y: current.Y + dy[i], Layer: current.Layer}
			if !visited[nextPoint] {
				nextTile, tileExists := tilesCache[nextPoint]
				if tileExists {
					if nextTile.Access.Type != "blocked" {
						visited[nextPoint] = true
						queue = append(queue, nextPoint)
					}
				}
			}
		}
	}

	return foundTile
}

var initTilesCache = sync.OnceFunc(func() {
	tilesCache = make(map[Point]MapSchema)
	var tiles []MapSchema
	err := json.Unmarshal(maps, &tiles)
	if err != nil {
		return
	}
	for _, tile := range tiles {
		p := Point{X: tile.X, Y: tile.Y, Layer: tile.Layer}
		tilesCache[p] = tile
	}
})
