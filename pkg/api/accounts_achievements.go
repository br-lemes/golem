package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func AccountsAchievements(account string) ([]schemas.AccountAchievementSchema, error) {
	if account == "" {
		account = cache.GetAccount()
		if account == "" {
			MyDetails()
			account = cache.GetAccount()
		}
	}
	result := []schemas.AccountAchievementSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/accounts/%s/achievements?page=%d",
			account, page), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.DataPageAccountAchievementSchema
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
