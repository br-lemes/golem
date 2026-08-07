package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionGrandexchangeCancel(name string, cancel schemas.GECancelOrderSchema) (schemas.GETransactionListSchema, error) {
	path := fmt.Sprintf("/my/%s/action/grandexchange/cancel", name)
	resp, err := PostNoCooldown(path, cancel)
	if err != nil {
		return schemas.GETransactionListSchema{}, err
	}
	var data schemas.GETransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.GETransactionListSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds,
		string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
