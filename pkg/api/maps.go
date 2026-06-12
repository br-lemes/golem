package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func Maps() ([]schemas.MapSchema, error) {
	result := []schemas.MapSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/maps?page=%d", page), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageMapSchema
		err = json.Unmarshal(resp, &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data.Data...)
		pages := 0
		if data.Pages != nil {
			pages = *data.Pages
		}
		if page >= pages {
			break
		}
		page++
	}
	return result, nil
}
