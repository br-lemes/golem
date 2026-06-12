package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionRest(name string) (schemas.CharacterRestDataSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/rest", name), nil)
	if err != nil {
		return schemas.CharacterRestDataSchema{}, err
	}
	var data schemas.CharacterRestResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterRestDataSchema{}, err
	}
	return data.Data, nil
}
