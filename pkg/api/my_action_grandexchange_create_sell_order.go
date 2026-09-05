package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionGrandexchangeCreateSellOrder(name string, order schemas.GEOrderCreationSchema) (schemas.GEOrderTransactionSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/grandexchange/create_sell_order", name)
	resp, err := post(path, order)
	if err != nil {
		return schemas.GEOrderTransactionSchema{}, err
	}
	var data schemas.GECreateOrderTransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.GEOrderTransactionSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
