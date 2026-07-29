package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionTransition(name string) (schemas.CharacterTransitionDataSchema, error) {
	path := fmt.Sprintf("/my/%s/action/transition", name)
	resp, err := PostNoCooldown(path, nil)
	if err != nil {
		return schemas.CharacterTransitionDataSchema{}, err
	}
	var data schemas.CharacterTransitionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterTransitionDataSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds,
		string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
