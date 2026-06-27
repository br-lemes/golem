package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func Effects() ([]schemas.EffectSchema, error) {
	result := []schemas.EffectSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/effects?page=%d", page), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageEffectSchema
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
