package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const EffectsSize = 10000

type EffectsOptions struct{}

func Effects(options EffectsOptions) ([]schemas.EffectSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.EffectSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(EffectsSize))
		path := fmt.Sprintf("/effects?%s", params.Encode())
		resp, err := Get(path, nil)
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
