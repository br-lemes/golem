package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionUnequip(name string, unequips []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error) {
	path := fmt.Sprintf("/my/%s/action/unequip", name)
	resp, err := PostNoCooldown(path, unequips)
	if err != nil {
		return schemas.EquipmentTransactionSchema{}, err
	}
	var data schemas.EquipmentResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.EquipmentTransactionSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds)
	return data.Data, nil
}
