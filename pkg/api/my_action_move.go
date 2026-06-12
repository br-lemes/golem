package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionMove(name string, x, y int) (schemas.CharacterMovementDataSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/move", name),
		map[string]int{"x": x, "y": y})
	if err != nil {
		return schemas.CharacterMovementDataSchema{}, err
	}
	var data schemas.CharacterMovementResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterMovementDataSchema{}, err
	}
	return data.Data, nil
}
