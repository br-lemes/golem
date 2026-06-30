package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionRecycling(name string, item schemas.SimpleItemSchema) (schemas.RecyclingDataSchema, error) {
	path := fmt.Sprintf("/my/%s/action/recycling", name)
	resp, err := PostNoCooldown(path, item)
	if err != nil {
		return schemas.RecyclingDataSchema{}, err
	}
	var data schemas.RecyclingResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.RecyclingDataSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds)
	return data.Data, nil
}
