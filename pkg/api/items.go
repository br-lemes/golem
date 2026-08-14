package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const ItemsSize = 10000

func Items() ([]schemas.ItemSchema, error) {
	result := []schemas.ItemSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/items?page=%d&size=%d", page, ItemsSize), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageItemSchema
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
