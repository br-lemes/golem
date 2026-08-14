package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const MyGrandexchangeOrdersSize = 100

func MyGrandexchangeOrders(code string, orderType string) ([]schemas.GEOrderSchema, error) {
	result := []schemas.GEOrderSchema{}
	page := 1
	for {
		path := fmt.Sprintf("/my/grandexchange/orders?page=%d&size=%d", page, MyGrandexchangeOrdersSize)
		if code != "" {
			path += fmt.Sprintf("&code=%s", code)
		}
		if orderType != "" {
			path += fmt.Sprintf("&type=%s", orderType)
		}
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
