package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const GrandexchangeOrdersSize = 100

type GrandexchangeOrdersOptions struct {
	Account  string `url:"account,omitempty"`
	Code     string `url:"code,omitempty"`
	ItemType string `url:"item_type,omitempty"`
	Type     string `url:"type,omitempty"`
}

func GrandexchangeOrder(id string) (schemas.GEOrderSchema, error) {
	path := fmt.Sprintf("/grandexchange/orders/%s", url.PathEscape(id))
	resp, err := Get(path, nil)
	if err != nil {
		return schemas.GEOrderSchema{}, err
	}
	var data schemas.GEOrderResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.GEOrderSchema{}, err
	}
	return data.Data, nil
}

func GrandexchangeOrders(options GrandexchangeOrdersOptions) ([]schemas.GEOrderSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}

	var result []schemas.GEOrderSchema
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(GrandexchangeOrdersSize))
		path := fmt.Sprintf("/grandexchange/orders?%s", params.Encode())
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
