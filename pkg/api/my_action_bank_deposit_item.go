package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionBankDepositItem(name string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/bank/deposit/item", name),
		items)
	if err != nil {
		return schemas.BankItemTransactionSchema{}, err
	}
	var data schemas.BankItemTransactionResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.BankItemTransactionSchema{}, err
	}
	return data.Data, nil
}
