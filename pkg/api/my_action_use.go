package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionUse(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
	path := fmt.Sprintf("/my/%s/action/use", name)
	resp, err := PostNoCooldown(path, item)
	if err != nil {
		return schemas.UseItemSchema{}, err
	}
	var data schemas.UseItemResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.UseItemSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds,
		string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
