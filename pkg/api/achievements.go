package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const AchievementsSize = 10000

type AchievementsOptions struct {
	Type string `url:"type,omitempty"`
}

func Achievements(options AchievementsOptions) ([]schemas.AchievementSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}

	result := []schemas.AchievementSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(AchievementsSize))
		path := fmt.Sprintf("/achievements?%s", params.Encode())
		resp, err := Get(path, nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageAchievementSchema
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
