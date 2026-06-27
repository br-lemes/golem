package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func Monsters() ([]schemas.MonsterSchema, error) {
	result := []schemas.MonsterSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/monsters?page=%d", page), nil)
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
