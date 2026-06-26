package task

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Move(character schemas.CharacterSchema, code string) (schemas.CharacterSchema, error) {
	result := database.FindClosest(character, code)
	if result == nil {
		return character, fmt.Errorf("no coordinates found for code %s", code)
	}
	if result.Transition == nil {
		return makeMove(character, result.Target)
	}
	character, err := makeMove(character, *result.Transition)
	if err != nil {
		return character, err
	}
	transitionData, err := api.MyActionTransition(character.Name)
	if err != nil {
		return character, err
	}
	return makeMove(transitionData.Character, result.Target)
}

func makeMove(character schemas.CharacterSchema, target schemas.MapSchema) (schemas.CharacterSchema, error) {
	if target.X == character.X && target.Y == character.Y {
		return character, nil
	}
	moveData, err := api.MyActionMove(character.Name, target.X, target.Y)
	if err != nil {
		return character, err
	}
	return moveData.Character, nil
}
