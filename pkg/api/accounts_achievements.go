package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const AccountsAchievementsSize = 100

type AccountsAchievementsOptions struct {
	Completed bool   `url:"completed,omitempty"`
	Type      string `url:"type,omitempty"`
}

func AccountsAchievements(account string, options AccountsAchievementsOptions) ([]schemas.AccountAchievementSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.AccountsAchievements(account, options)
}

func (c *Client) AccountsAchievements(account string, options AccountsAchievementsOptions) ([]schemas.AccountAchievementSchema, error) {
	if account == "" {
		var err error
		account, err = c.accountName()
		if err != nil {
			return nil, err
		}
	}
	params, err := query.Values(options)
	if err != nil {
		//+gocover:ignore:block typed options cannot fail query encoding
		return nil, err
	}
	result := []schemas.AccountAchievementSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(AccountsAchievementsSize))
		path := fmt.Sprintf("/accounts/%s/achievements?%s", url.PathEscape(account), params.Encode())
		resp, err := c.Get(path, nil)
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

func (c *Client) accountName() (string, error) {
	account := cache.GetAccount()
	if account != "" {
		return account, nil
	}
	_, err := c.MyDetails()
	if err != nil {
		return "", err
	}
	return cache.GetAccount(), nil
}
