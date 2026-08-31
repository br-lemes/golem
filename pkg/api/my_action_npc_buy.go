package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionNPCBuy(name string, item schemas.SimpleItemSchema) (schemas.NpcMerchantTransactionSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/npc/buy", name)
	resp, err := PostNoCooldown(path, item)
	if err != nil {
		return schemas.NpcMerchantTransactionSchema{}, err
	}
	var data schemas.NpcMerchantTransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.NpcMerchantTransactionSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
