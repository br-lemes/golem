package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionGrandexchangeBuy(name string, buy schemas.GEBuyOrderSchema) (schemas.GETransactionListSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/grandexchange/buy", name)
	resp, err := PostNoCooldown(path, buy)
	if err != nil {
		return schemas.GETransactionListSchema{}, err
	}
	var data schemas.GETransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.GETransactionListSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
