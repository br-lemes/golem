package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

const AccountsAchievementsSize = 100

func AccountsAchievements(account string) ([]schemas.AccountAchievementSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.AccountsAchievements(account)
}

func (c *Client) AccountsAchievements(account string) ([]schemas.AccountAchievementSchema, error) {
	if account == "" {
		account = cache.GetAccount()
		if account == "" {
			_, err := c.MyDetails()
			if err != nil {
				return nil, err
			}
			account = cache.GetAccount()
		}
	}
	result := []schemas.AccountAchievementSchema{}
	page := 1
	for {
		resp, err := c.Get(fmt.Sprintf("/accounts/%s/achievements?page=%d&size=%d", account, page, AccountsAchievementsSize), nil)
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
