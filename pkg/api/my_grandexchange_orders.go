package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const MyGrandexchangeOrdersSize = 100

type MyGrandexchangeOrdersOptions struct {
	Code string `url:"code,omitempty"`
	Type string `url:"type,omitempty"`
}

func MyGrandexchangeOrders(options MyGrandexchangeOrdersOptions) ([]schemas.GEOrderSchema, error) {
	result := []schemas.GEOrderSchema{}
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(MyGrandexchangeOrdersSize))
		path := fmt.Sprintf("/my/grandexchange/orders?%s", params.Encode())
		resp, err := Get(path, nil)
		if err != nil {
			return nil, err
		}
		var data schemas.DataPageGEOrderSchema
		err = json.Unmarshal(resp, &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data.Data...)
		if page >= data.Pages {
			break
		}
		page++
	}
	return result, nil
}
