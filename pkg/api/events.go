package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

const EventsSize = 10000

func Events() ([]schemas.EventSchema, error) {
	result := []schemas.EventSchema{}
	page := 1
	for {
		resp, err := Get(fmt.Sprintf("/events?page=%d&size=%d",
			page, EventsSize), nil)
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
