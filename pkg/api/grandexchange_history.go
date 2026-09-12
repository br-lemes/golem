package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const GrandexchangeHistorySize = 100

type GrandexchangeHistoryOptions struct {
	Account string `url:"account,omitempty"`
}

func GrandexchangeHistory(code string, options GrandexchangeHistoryOptions) ([]schemas.GEOrderHistorySchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.GrandexchangeHistory(code, options)
}

func (c *Client) GrandexchangeHistory(code string, options GrandexchangeHistoryOptions) ([]schemas.GEOrderHistorySchema, error) {
	params, err := query.Values(options)
	if err != nil {
		//+gocover:ignore:block typed options cannot fail query encoding
		return nil, err
	}

	result := []schemas.GEOrderHistorySchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(GrandexchangeHistorySize))
		path := fmt.Sprintf("/grandexchange/history/%s?%s", url.PathEscape(code), params.Encode())
		resp, err := c.Get(path, nil)
		if err != nil {
			return nil, err
		}
		var data schemas.DataPageGEOrderHistorySchema
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
