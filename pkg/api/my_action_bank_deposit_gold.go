package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionBankDepositGold(name string, quantity int) (schemas.BankGoldTransactionSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/bank/deposit/gold", name),
		schemas.GoldSchema{Quantity: quantity})
	if err != nil {
		return schemas.BankGoldTransactionSchema{}, err
	}
	var data schemas.BankGoldTransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.BankGoldTransactionSchema{}, err
	}
	return data.Data, nil
}
