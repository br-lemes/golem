package api

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyDetails() (schemas.MyAccountDetails, error) {
	resp, err := Get("/my/details", nil)
	if err != nil {
		return schemas.MyAccountDetails{}, err
	}
	var data schemas.MyAccountDetailsSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.MyAccountDetails{}, err
	}
	return data.Data, nil
}
