package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionUse(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/use", name)
	resp, err := post(path, item)
	if err != nil {
		return schemas.UseItemSchema{}, err
	}
	var data schemas.UseItemResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.UseItemSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
