package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const NpcsDetailsSize = 10000

type NpcsDetailsOptions struct {
	Currency string `url:"currency,omitempty"`
	Item     string `url:"item,omitempty"`
	Name     string `url:"name,omitempty"`
	Type     string `url:"type,omitempty"`
}

func NpcsDetails(options NpcsDetailsOptions) ([]schemas.NPCSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.NPCSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(NpcsDetailsSize))
		path := fmt.Sprintf("/npcs/details?%s", params.Encode())
		resp, err := Get(path, nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageNPCSchema
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
