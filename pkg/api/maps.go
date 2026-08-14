package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const MapsSize = 10000

func Maps() ([]schemas.MapSchema, error) {
	result := []schemas.MapSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/maps?page=%d&size=%d", page, MapsSize), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageMapSchema
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
