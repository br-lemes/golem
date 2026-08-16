package database

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMapsSmoke(t *testing.T) {
	if len(Maps.All()) == 0 {
		t.Fatal("maps catalog is empty")
	}
}

func TestFindClosest(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {X: 0, Y: 0},
		{X: 1, Y: 0}: {
			X: 1,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "slime"},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result == nil {
		t.Fatal("findClosest returned nil")
	}
	if result.Target.X != 1 || result.Target.Y != 0 {
		t.Fatalf("findClosest returned (%d, %d), want (1, 0)", result.Target.X, result.Target.Y)
	}
}

func TestFindClosestReturnsStartWhenItIsTarget(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {
			X: 0,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "slime"},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result == nil || result.Target.X != 0 || result.Target.Y != 0 {
		t.Fatalf("findClosest() = %#v, want start tile", result)
	}
}

func TestFindClosestReturnsNilWhenTargetDoesNotExist(t *testing.T) {
	maps := map[Point]schemas.MapSchema{{X: 0, Y: 0}: {X: 0, Y: 0}}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result != nil {
		t.Fatalf("findClosest() = %#v, want nil", result)
	}
}

func TestFindClosestIgnoresBlockedTiles(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {X: 0, Y: 0},
		{X: 1, Y: 0}: {
			X:      1,
			Y:      0,
			Access: schemas.AccessSchema{Type: "blocked"},
		},
		{X: 2, Y: 0}: {
			X: 2,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "slime"},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result != nil {
		t.Fatalf("findClosest found target through blocked tile: %+v", result.Target)
	}
}

func TestFindClosestFindsPathAroundBlockedTile(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {X: 0, Y: 0},
		{X: 1, Y: 0}: {
			X:      1,
			Y:      0,
			Access: schemas.AccessSchema{Type: "blocked"},
		},
		{X: 0, Y: 1}: {X: 0, Y: 1},
		{X: 1, Y: 1}: {X: 1, Y: 1},
		{X: 2, Y: 1}: {X: 2, Y: 1},
		{X: 2, Y: 0}: {
			X: 2,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "slime"},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result == nil || result.Target.X != 2 || result.Target.Y != 0 {
		t.Fatalf("findClosest() = %#v, want target at (2, 0)", result)
	}
}

func TestFindClosestUsesTransition(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {
			X: 0,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Transition: &schemas.TransitionSchema{X: 2, Y: 0},
			},
		},
		{X: 2, Y: 0}: {
			X: 2,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "slime"},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result == nil {
		t.Fatal("findClosest() returned nil")
	}
	if len(result.Transitions) != 1 {
		t.Fatalf("findClosest() transitions = %d, want 1", len(result.Transitions))
	}
}

func TestFindClosestPreservesFirstTransition(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {
			X: 0,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Transition: &schemas.TransitionSchema{X: 1, Y: 0},
			},
		},
		{X: 1, Y: 0}: {
			X: 1,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Transition: &schemas.TransitionSchema{X: 2, Y: 0},
			},
		},
		{X: 2, Y: 0}: {
			X: 2,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "slime"},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result == nil {
		t.Fatal("findClosest() returned nil")
	}
	if len(result.Transitions) != 2 || result.Transitions[0].X != 0 || result.Transitions[1].X != 1 {
		t.Fatalf("findClosest() transitions = %#v, want both transitions", result.Transitions)
	}
}

func TestFindClosestDoesNotCrossLayers(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0, Layer: "overworld"}: {X: 0, Y: 0, Layer: "overworld"},
		{X: 0, Y: 0, Layer: "underground"}: {
			X:     0,
			Y:     0,
			Layer: "underground",
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: "slime"},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X:     0,
		Y:     0,
		Layer: "overworld",
	}, "slime", nil)
	if result != nil {
		t.Fatalf("findClosest() crossed layers: %#v", result.Target)
	}
}

func TestFindClosestHandlesTransitionCycle(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {
			X: 0,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Transition: &schemas.TransitionSchema{X: 1, Y: 0},
			},
		},
		{X: 1, Y: 0}: {
			X: 1,
			Y: 0,
			Interactions: schemas.InteractionSchema{
				Transition: &schemas.TransitionSchema{X: 0, Y: 0},
			},
		},
	}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "slime", nil)
	if result != nil {
		t.Fatalf("findClosest() found target in cycle: %#v", result)
	}
}

func TestFindClosestFindsEventPoint(t *testing.T) {
	maps := map[Point]schemas.MapSchema{
		{X: 0, Y: 0}: {X: 0, Y: 0},
		{X: 1, Y: 0}: {X: 1, Y: 0},
	}
	eventPoints := EventPoints{{X: 1, Y: 0}: true}

	result := findClosest(testMapLookup(maps), schemas.CharacterSchema{
		X: 0,
		Y: 0,
	}, "event", eventPoints)
	if result == nil {
		t.Fatal("findClosest did not find event point")
	}
	if result.Target.X != 1 || result.Target.Y != 0 {
		t.Fatalf("findClosest returned (%d, %d), want (1, 0)", result.Target.X, result.Target.Y)
	}
}

func testMapLookup(maps map[Point]schemas.MapSchema) mapLookup {
	return func(point Point) (*schemas.MapSchema, bool) {
		tile, exists := maps[point]
		return &tile, exists
	}
}
