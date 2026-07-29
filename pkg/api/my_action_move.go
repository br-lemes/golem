package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionMove(name string, x, y int) (schemas.CharacterMovementDataSchema, error) {
	path := fmt.Sprintf("/my/%s/action/move", name)
	resp, err := PostNoCooldown(path, map[string]int{"x": x, "y": y})
	if err != nil {
		return schemas.CharacterMovementDataSchema{}, err
	}
	var data schemas.CharacterMovementResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterMovementDataSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds,
		string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
