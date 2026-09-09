package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMoveUsesFirstFreePath(t *testing.T) {
	character := schemas.CharacterSchema{Name: "hero", Layer: "overworld"}
	var movedTo []struct{ x, y int }
	d := deps{
		accountsAchievements: func(string) ([]schemas.AccountAchievementSchema, error) {
			return nil, nil
		},
		myActionMove: func(_ string, x, y int) (schemas.CharacterMovementDataSchema, error) {
			movedTo = append(movedTo, struct{ x, y int }{x, y})
			character.X = x
			character.Y = y
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
	}

	got, err := move(d, character, "bank")
	if err != nil {
		t.Fatal(err)
	}
	if len(movedTo) != 1 {
		t.Fatalf("move calls = %d, want 1", len(movedTo))
	}
	if got.X != movedTo[0].x || got.Y != movedTo[0].y {
		t.Fatalf("returned character = %#v, want position %#v", got, movedTo[0])
	}
}

func TestMoveReturnsErrorWithoutFreePath(t *testing.T) {
	character := schemas.CharacterSchema{Layer: "overworld"}

	_, err := move(deps{}, character, "invalid_code")
	if err == nil {
		t.Fatal("move() returned nil error, want error")
	}
}

func TestMoveExecutesTransitions(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     5,
		Y:     -4,
		Layer: "underground",
	}
	moves := 0
	transitions := 0
	d := deps{
		myActionMove: func(_ string, x, y int) (schemas.CharacterMovementDataSchema, error) {
			moves++
			character.X = x
			character.Y = y
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			transitions++
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
	}

	_, err := move(d, character, "mithril_rocks")
	if err != nil {
		t.Fatal(err)
	}
	if moves != 3 {
		t.Fatalf("move calls = %d, want 3", moves)
	}
	if transitions != 2 {
		t.Fatalf("transition calls = %d, want 2", transitions)
	}
}

func TestMoveReturnsTransitionError(t *testing.T) {
	character := schemas.CharacterSchema{X: 5, Y: -4, Layer: "underground"}
	wantErr := errors.New("transition failed")
	d := deps{
		myActionMove: func(_ string, x, y int) (schemas.CharacterMovementDataSchema, error) {
			character.X = x
			character.Y = y
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{}, wantErr
		},
	}

	_, err := move(d, character, "mithril_rocks")
	if !errors.Is(err, wantErr) {
		t.Fatalf("move() error = %v, want %v", err, wantErr)
	}
}

func TestMoveReturnsMoveErrorAtTransition(t *testing.T) {
	character := schemas.CharacterSchema{X: 5, Y: -4, Layer: "underground"}
	wantErr := errors.New("move failed")
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{}, wantErr
		},
	}

	_, err := move(d, character, "mithril_rocks")
	if !errors.Is(err, wantErr) {
		t.Fatalf("move() error = %v, want %v", err, wantErr)
	}
}
