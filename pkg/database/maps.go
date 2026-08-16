package database

import (
	_ "embed"
	"encoding/json"
	"slices"
	"sync"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/tidwall/gjson"
)

type (
	Point struct {
		X     int
		Y     int
		Layer schemas.MapLayer
	}

	SearchNode struct {
		Point           Point
		FirstTransition *schemas.MapSchema
	}

	SearchResult struct {
		Target     schemas.MapSchema
		Transition *schemas.MapSchema
	}

	EventPoints map[Point]bool
)

var (
	//go:embed maps.json
	maps []byte

	mapsList      []schemas.MapSchema
	mapsCache     map[Point]schemas.MapSchema
	mapsCodesList []string
)

func GetMaps() []schemas.MapSchema {
	initMapsCache()
	return mapsList
}

func GetMap(x, y int, layer schemas.MapLayer) (schemas.MapSchema, bool) {
	initMapsCache()
	tile, exists := mapsCache[Point{X: x, Y: y, Layer: layer}]
	return tile, exists
}

func GetMapCodes() []string {
	initMapsCodesCache()
	return mapsCodesList
}

func FindClosest(character schemas.CharacterSchema, code string) *SearchResult {
	initMapsCache()

	eventPoints := getEventPoints(code)

	startPoint := Point{X: character.X, Y: character.Y, Layer: character.Layer}
	startNode := SearchNode{Point: startPoint, FirstTransition: nil}

	queue := []SearchNode{startNode}
	visited := make(map[Point]bool)
	visited[startPoint] = true

	dx := []int{0, 0, 1, -1}
	dy := []int{1, -1, 0, 0}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		currentTile, exists := mapsCache[current.Point]
		if exists {
			var isTarget bool

			if currentTile.Interactions.Content != nil {
				if currentTile.Interactions.Content.Code == code {
					isTarget = true
				}
			}

			if !isTarget {
				if eventPoints[current.Point] {
					isTarget = true
				}
			}

			if isTarget {
				result := &SearchResult{
					Target:     currentTile,
					Transition: current.FirstTransition,
				}
				return result
			}

			if currentTile.Interactions.Transition != nil {
				transition := currentTile.Interactions.Transition
				nextPoint := Point{
					X:     transition.X,
					Y:     transition.Y,
					Layer: transition.Layer,
				}

				if !visited[nextPoint] {
					visited[nextPoint] = true

					var nextFirstTransition *schemas.MapSchema
					if current.FirstTransition != nil {
						nextFirstTransition = current.FirstTransition
					} else {
						nextFirstTransition = &currentTile
					}

					nextNode := SearchNode{
						Point:           nextPoint,
						FirstTransition: nextFirstTransition,
					}
					queue = append(queue, nextNode)
				}
			}
		}

		for i := 0; i < 4; i++ {
			nextPoint := Point{
				X:     current.Point.X + dx[i],
				Y:     current.Point.Y + dy[i],
				Layer: current.Point.Layer,
			}

			if !visited[nextPoint] {
				nextTile, tileExists := mapsCache[nextPoint]
				if tileExists {
					if nextTile.Access.Type != "blocked" {
						visited[nextPoint] = true
						nextNode := SearchNode{
							Point:           nextPoint,
							FirstTransition: current.FirstTransition,
						}
						queue = append(queue, nextNode)
					}
				}
			}
		}
	}

	return nil
}

func getEventPoints(code string) EventPoints {
	initMapsCodesCache()
	if slices.Contains(mapsCodesList, code) {
		return EventPoints{}
	}

	if !slices.Contains(EventContentCodes(), code) {
		return EventPoints{}
	}

	events, err := api.EventsActive()
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

var initMapsCache = sync.OnceFunc(func() {
	mapsCache = make(map[Point]schemas.MapSchema)
	err := json.Unmarshal(maps, &mapsList)
	if err != nil {
		panic("failed to unmarshal maps: " + err.Error())
	}
	for _, tile := range mapsList {
		p := Point{X: tile.X, Y: tile.Y, Layer: tile.Layer}
		mapsCache[p] = tile
	}
})

var initMapsCodesCache = sync.OnceFunc(func() {
	uniqueCodes := make(map[string]struct{})
	gjson.GetBytes(maps, "#.interactions.content.code").ForEach(func(_, value gjson.Result) bool {
		s := value.String()
		if s != "" {
			uniqueCodes[s] = struct{}{}
		}
		return true
	})
	for code := range uniqueCodes {
		mapsCodesList = append(mapsCodesList, code)
	}
})
