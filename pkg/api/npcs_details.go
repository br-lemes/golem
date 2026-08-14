package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const NpcsDetailsSize = 10000

func NpcsDetails() ([]schemas.NPCSchema, error) {
	result := []schemas.NPCSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/npcs/details?page=%d&size=%d", page, NpcsDetailsSize), nil)
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
