package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func Characters(name string) (schemas.CharacterSchema, error) {
	resp, err := Get(fmt.Sprintf("/characters/%s", name), nil)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	var data schemas.CharacterResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	return data.Data, nil
}
