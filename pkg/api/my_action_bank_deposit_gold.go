package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionBankDepositGold(name string, quantity int) (schemas.BankGoldTransactionSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/bank/deposit/gold", name)
	resp, err := post(path, schemas.GoldSchema{Quantity: quantity})
	if err != nil {
		return schemas.BankGoldTransactionSchema{}, err
	}
	var data schemas.BankGoldTransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.BankGoldTransactionSchema{}, err
	}
	cache.UpdateBankGold(data.Data.Bank.Quantity)
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
