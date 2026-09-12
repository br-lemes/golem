package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionRest(name string) (schemas.CharacterRestDataSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.MyActionRest(name)
}

func (c *Client) MyActionRest(name string) (schemas.CharacterRestDataSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/rest", name)
	resp, err := c.Post(path, nil)
	if err != nil {
		return schemas.CharacterRestDataSchema{}, err
	}
	var data schemas.CharacterRestResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterRestDataSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
