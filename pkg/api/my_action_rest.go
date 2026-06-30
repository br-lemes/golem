package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionRest(name string) (schemas.CharacterRestDataSchema, error) {
	path := fmt.Sprintf("/my/%s/action/rest", name)
	resp, err := PostNoCooldown(path, nil)
	if err != nil {
		return schemas.CharacterRestDataSchema{}, err
	}
	var data schemas.CharacterRestResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterRestDataSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds)
	return data.Data, nil
}
