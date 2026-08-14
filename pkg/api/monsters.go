package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const MonstersSize = 10000

func Monsters() ([]schemas.MonsterSchema, error) {
	result := []schemas.MonsterSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/monsters?page=%d&size=%d", page, MonstersSize), nil)
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
