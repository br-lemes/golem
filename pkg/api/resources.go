package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const ResourcesSize = 10000

func Resources() ([]schemas.ResourceSchema, error) {
	result := []schemas.ResourceSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/resources?page=%d&size=%d", page, ResourcesSize), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageResourceSchema
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
