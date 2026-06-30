package api

import (
	"encoding/json"
	"net/http"
)

func Post(path string, data any) ([]byte, error) {
	return post(path, data, true)
}

func PostNoCooldown(path string, data any) ([]byte, error) {
	return post(path, data, false)
}

func post(path string, data any, cooldown bool) ([]byte, error) {
	var body []byte
	var err error
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	return Request(http.MethodPost, path, body, cooldown)
}
