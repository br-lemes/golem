package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionFight(name string, participants []string) (schemas.CharacterFightDataSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/fight", name),
		schemas.FightRequestSchema{Participants: &participants})
	if err != nil {
		return schemas.CharacterFightDataSchema{}, err
	}
	var data schemas.CharacterFightResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterFightDataSchema{}, err
	}
	return data.Data, nil
}
