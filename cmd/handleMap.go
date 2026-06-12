package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func handleMap(character schemas.CharacterSchema, code string) (schemas.CharacterSchema, error) {
	target := database.FindClosest(character, code)
	if target == nil {
		return character, fmt.Errorf("no coordinates found for code %s", code)
	}
	if target.X == character.X &&
		target.Y == character.Y &&
		target.Layer == character.Layer {
		return character, nil
	}
	moveData, err := api.MyActionMove(character.Name, target.X, target.Y)
	if err != nil {
		return character, err
	}
	return moveData.Character, nil
}
