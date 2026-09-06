package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const ResourcesSize = 10000

type ResourcesOptions struct {
	Drop     string `url:"drop,omitempty"`
	MaxLevel int    `url:"max_level,omitempty"`
	MinLevel int    `url:"min_level,omitempty"`
	Skill    string `url:"skill,omitempty"`
}

func Resources(options ResourcesOptions) ([]schemas.ResourceSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.ResourceSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(ResourcesSize))
		path := fmt.Sprintf("/resources?%s", params.Encode())
		resp, err := Get(path, nil)
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
