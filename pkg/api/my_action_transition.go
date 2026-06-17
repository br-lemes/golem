package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionTransition(name string) (schemas.CharacterTransitionDataSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/transition", name), nil)
	if err != nil {
		return schemas.CharacterTransitionDataSchema{}, err
	}
	var data schemas.CharacterTransitionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterTransitionDataSchema{}, err
	}
	return data.Data, nil
}
