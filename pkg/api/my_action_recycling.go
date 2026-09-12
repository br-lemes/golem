package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionRecycling(name string, item schemas.RecyclingSchema) (schemas.RecyclingDataSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.MyActionRecycling(name, item)
}

func (c *Client) MyActionRecycling(name string, item schemas.RecyclingSchema) (schemas.RecyclingDataSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/recycling", name)
	resp, err := c.Post(path, item)
	if err != nil {
		return schemas.RecyclingDataSchema{}, err
	}
	var data schemas.RecyclingResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.RecyclingDataSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
