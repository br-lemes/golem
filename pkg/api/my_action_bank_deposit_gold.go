package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionBankDepositGold(name string, quantity int) (schemas.BankGoldTransactionSchema, error) {
	path := fmt.Sprintf("/my/%s/action/bank/deposit/gold", name)
	resp, err := PostNoCooldown(path, schemas.GoldSchema{Quantity: quantity})
	if err != nil {
		return schemas.BankGoldTransactionSchema{}, err
	}
	var data schemas.BankGoldTransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.BankGoldTransactionSchema{}, err
	}
	cache.UpdateBankGold(data.Data.Bank.Quantity)
	cache.SaveCharacter(name, data.Data.Character)
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
