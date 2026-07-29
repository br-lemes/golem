package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func GrandexchangeOrders(id string) ([]schemas.GEOrderSchema, error) {
	if id == "" {
		return grandexchangeOrdersAll()
	}
	return grandexchangeOrdersById(id)
}

func grandexchangeOrdersAll() ([]schemas.GEOrderSchema, error) {
	result := []schemas.GEOrderSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/grandexchange/orders?page=%d", page), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.DataPageGEOrderSchema
		err = json.Unmarshal(resp, &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data.Data...)
		if data.Page >= data.Pages {
			break
		}
		page++
	}
	return result, nil
}

func grandexchangeOrdersById(id string) ([]schemas.GEOrderSchema, error) {
	path := fmt.Sprintf("/grandexchange/orders/%s", id)
	resp, err := Get(path, nil)
	if err != nil {
		return nil, err
	}
	var data schemas.GEOrderResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return []schemas.GEOrderSchema{data.Data}, nil
}
