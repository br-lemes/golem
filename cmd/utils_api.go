package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

const baseURL = "https://api.artifactsmmo.com"

var httpClient = &http.Client{Timeout: 30 * time.Second}

func apiGet(path string, data map[string]string) ([]byte, error) {
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
	return apiRequest(http.MethodGet, path, nil)
}

func apiPost(path string, data any) ([]byte, error) {
	var body []byte
	var err error
	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	return apiRequest(http.MethodPost, path, body)
}

func apiRequest(method, path string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)

	if debugFlag {
		fmt.Fprintf(os.Stderr, "→ %s %s\n", method, path)
		if len(body) > 0 {
			fmt.Fprintf(os.Stderr, "  Body: %s\n", string(body))
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errMsg := gjson.GetBytes(respBytes, "error.message")
	if errMsg.Exists() {
		return nil, fmt.Errorf("%s", errMsg.String())
	}

	cdResult := gjson.GetBytes(respBytes, "data.cooldown.total_seconds")
	cd := int(cdResult.Int())

	if cd > 0 {
		fmt.Fprintf(os.Stderr, "⏳ Cooldown started: %d seconds\n", cd)
		time.Sleep(time.Duration(cd) * time.Second)
	}

	return respBytes, nil
}
