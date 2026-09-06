package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const MonstersSize = 10000

type MonstersOptions struct {
	Drop     string `url:"drop,omitempty"`
	MaxLevel int    `url:"max_level,omitempty"`
	MinLevel int    `url:"min_level,omitempty"`
	Name     string `url:"name,omitempty"`
}

func Monsters(options MonstersOptions) ([]schemas.MonsterSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.MonsterSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(MonstersSize))
		path := fmt.Sprintf("/monsters?%s", params.Encode())
		resp, err := Get(path, nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageMonsterSchema
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
