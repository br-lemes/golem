package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionBankDepositItem(name string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/bank/deposit/item", name)
	resp, err := post(path, items)
	if err != nil {
		return schemas.BankItemTransactionSchema{}, err
	}
	var data schemas.BankItemTransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.BankItemTransactionSchema{}, err
	}
	cache.SaveBankItems(data.Data.Bank)
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
