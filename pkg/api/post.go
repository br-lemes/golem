package api

import (
	"encoding/json"
	"net/http"

	"github.com/tidwall/gjson"
)

func Post(path string, data any) ([]byte, error) {
	resp, err := post(path, data)
	if err != nil {
		return nil, err
	}
	cd := int(gjson.GetBytes(resp, "data.cooldown.total_seconds").Int())
	if cd > 0 {
		reason := gjson.GetBytes(resp, "data.cooldown.reason")
		handleCooldown(cd, reason.String())
	}
	return resp, nil
}

func post(path string, data any) ([]byte, error) {
	var body []byte
	var err error
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	return Request(http.MethodPost, path, body)
}
