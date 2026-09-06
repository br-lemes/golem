package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const EventsSize = 10000

type EventsOptions struct {
	Type string `url:"type,omitempty"`
}

func Events(options EventsOptions) ([]schemas.EventSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.EventSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(EventsSize))
		path := fmt.Sprintf("/events?%s", params.Encode())
		resp, err := Get(path, nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageEventSchema
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
