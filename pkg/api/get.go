package api

import (
	"net/http"
	"net/url"
)

func Get(path string, data map[string]string) ([]byte, error) {
	if len(data) > 0 {
		values := url.Values{}
		for key, val := range data {
			if val != "" {
				values.Add(key, val)
			}
		}
		queryString := values.Encode()
		if queryString != "" {
			path = path + "?" + queryString
		}
	}
	return Request(http.MethodGet, path, nil)
}
