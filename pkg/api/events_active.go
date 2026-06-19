package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func EventsActive() ([]schemas.ActiveEventSchema, error) {
	result := []schemas.ActiveEventSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/events/active?page=%d", page), nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageActiveEventSchema
		err = json.Unmarshal(resp, &data)
		if err != nil {
			return nil, err
		}
		result = append(result, data.Data...)
		pages := 0
		if data.Pages != nil {
			pages = *data.Pages
		}
		if page >= pages {
			break
		}
		page++
	}
	return result, nil
}
