package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/google/go-querystring/query"
)

const MapsSize = 10000

type MapsOptions struct {
	ContentCode     string `url:"content_code,omitempty"`
	ContentType     string `url:"content_type,omitempty"`
	HideBlockedMaps bool   `url:"hide_blocked_maps,omitempty"`
	HideEvent       bool   `url:"hide_event,omitempty"`
	Layer           string `url:"layer,omitempty"`
	Transition      bool   `url:"transition,omitempty"`
}

func Maps(options MapsOptions) ([]schemas.MapSchema, error) {
	params, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	result := []schemas.MapSchema{}
	page := 1
	for {
		params.Set("page", strconv.Itoa(page))
		params.Set("size", strconv.Itoa(MapsSize))
		path := fmt.Sprintf("/maps?%s", params.Encode())
		resp, err := Get(path, nil)
		if err != nil {
			return nil, err
		}
		var data schemas.StaticDataPageMapSchema
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
