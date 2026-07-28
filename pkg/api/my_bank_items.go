package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyBankItems() ([]schemas.SimpleItemSchema, error) {
	bankItems := cache.GetBankItems()
	if bankItems != nil {
		return bankItems, nil
	}
	result := []schemas.SimpleItemSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/my/bank/items?page=%d", page), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.DataPageSimpleItemSchema
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
	cache.SaveBankItems(result)
	return result, nil
}
