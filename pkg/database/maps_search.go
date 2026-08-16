package database

import "github.com/br-lemes/golem/pkg/schemas"

type mapLookup func(Point) (*schemas.MapSchema, bool)

func findClosest(maps mapLookup, character schemas.CharacterSchema, code string, eventPoints EventPoints) *SearchResult {
	startPoint := Point{X: character.X, Y: character.Y, Layer: character.Layer}
	startNode := SearchNode{Point: startPoint}

	queue := []SearchNode{startNode}
	visited := make(map[Point]bool)
	visited[startPoint] = true

	dx := []int{0, 0, 1, -1}
	dy := []int{1, -1, 0, 0}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		currentTile, exists := maps(current.Point)
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
				return &SearchResult{
					Target:      *currentTile,
					Transitions: current.Transitions,
				}
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

					nextTransitions := make([]schemas.MapSchema, len(current.Transitions), len(current.Transitions)+1)
					copy(nextTransitions, current.Transitions)
					nextTransitions = append(nextTransitions, *currentTile)

					queue = append(queue, SearchNode{
						Point:       nextPoint,
						Transitions: nextTransitions,
					})
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
				nextTile, tileExists := maps(nextPoint)
				if tileExists {
					if nextTile.Access.Type != "blocked" {
						visited[nextPoint] = true
						queue = append(queue, SearchNode{
							Point:       nextPoint,
							Transitions: current.Transitions,
						})
					}
				}
			}
		}
	}

	return nil
}
