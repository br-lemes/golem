package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const ItemsSize = 10000

type ItemsOptions struct {
	CraftMaterial string `url:"craft_material,omitempty"`
	CraftSkill    string `url:"craft_skill,omitempty"`
	MaxLevel      int    `url:"max_level,omitempty"`
	MinLevel      int    `url:"min_level,omitempty"`
	Name          string `url:"name,omitempty"`
	Type          string `url:"type,omitempty"`
}

func Items(options ItemsOptions) ([]schemas.ItemSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.ItemSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(ItemsSize))
		path := fmt.Sprintf("/items?%s", params.Encode())
		resp, err := Get(path, nil)
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
