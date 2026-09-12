package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionUnequip(name string, unequips []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.MyActionUnequip(name, unequips)
}

func (c *Client) MyActionUnequip(name string, unequips []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/unequip", name)
	resp, err := c.Post(path, unequips)
	if err != nil {
		return schemas.EquipmentTransactionSchema{}, err
	}
	var data schemas.EquipmentResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.EquipmentTransactionSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
