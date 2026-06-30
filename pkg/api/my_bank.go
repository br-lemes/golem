package api

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyBank() (schemas.BankSchema, error) {
	bank := cache.GetBank()
	if bank != nil {
		return *bank, nil
	}
	resp, err := Get("/my/bank", nil)
	if err != nil {
		return schemas.BankSchema{}, err
	}
	var data schemas.BankResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.BankSchema{}, err
	}
	cache.SaveBank(data.Data)
	return data.Data, nil
}
