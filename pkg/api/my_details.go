package api

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyDetails() (schemas.MyAccountDetails, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.MyDetails()
}

func (c *Client) MyDetails() (schemas.MyAccountDetails, error) {
	resp, err := c.Get("/my/details", nil)
	if err != nil {
		return schemas.MyAccountDetails{}, err
	}
	var data schemas.MyAccountDetailsSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.MyAccountDetails{}, err
	}
	cache.SaveAccount(data.Data.Username)
	return data.Data, nil
}
