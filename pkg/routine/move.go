package routine

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Move(character schemas.CharacterSchema, code string) (schemas.CharacterSchema, error) {
	return move(defaultDeps, character, code)
}

func move(d deps, character schemas.CharacterSchema, code string) (schemas.CharacterSchema, error) {
	result := database.FindClosest(character, code)
	if result == nil {
		return character, fmt.Errorf("no coordinates found for code %s", code)
	}
	for _, transition := range result.Transitions {
		character, err := makeMove(d, character, transition)
		if err != nil {
			return character, err
		}
		transitionData, err := d.myActionTransition(character.Name)
		if err != nil {
			return character, err
		}
		character = transitionData.Character
	}
	return makeMove(d, character, result.Target)
}

func makeMove(d deps, character schemas.CharacterSchema, target schemas.MapSchema) (schemas.CharacterSchema, error) {
	if target.X == character.X && target.Y == character.Y {
		return character, nil
	}
	moveData, err := d.myActionMove(character.Name, target.X, target.Y)
	if err != nil {
		return character, err
	}
	return moveData.Character, nil
}
