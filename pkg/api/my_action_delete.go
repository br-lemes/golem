package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionDelete(name string, item schemas.SimpleItemSchema) (schemas.DeleteItemSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.MyActionDelete(name, item)
}

func (c *Client) MyActionDelete(name string, item schemas.SimpleItemSchema) (schemas.DeleteItemSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/delete", name)
	resp, err := c.Post(path, item)
	if err != nil {
		return schemas.DeleteItemSchema{}, err
	}
	var data schemas.DeleteItemResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.DeleteItemSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
