package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const EventsActiveSize = 10000

func EventsActive() ([]schemas.ActiveEventSchema, error) {
	result := []schemas.ActiveEventSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/events/active?page=%d&size=%d", page, EventsActiveSize), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageActiveEventSchema
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
