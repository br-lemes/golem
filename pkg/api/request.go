package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/tidwall/gjson"
)

const baseURL = "https://api.artifactsmmo.com"

var Token string

var httpClient = &http.Client{Timeout: 30 * time.Second}

func Request(method, path string, body []byte) ([]byte, error) {
	for {
		req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
		if err != nil {
			cache.APILog(method, path, string(body), err.Error(), 0, 0)
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+Token)

		console.Debugf("%s %s\n", method, path)
		if len(body) > 0 {
			console.Debugf("  Body: %s\n", string(body))
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			format := "Network error: %v. Retrying in 5 seconds...\n"
			message := fmt.Sprintf(format, err)
			cache.APILog(method, path, string(body), message, 0, 0)
			console.Errorf(format, err)
			time.Sleep(5 * time.Second)
			continue
		}

		respBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			console.Errorf("Read error: %v. Retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		errMsg := gjson.GetBytes(respBytes, "error.message")
		if errMsg.Exists() {
			cache.APILog(method, path, string(body), string(respBytes),
				resp.StatusCode, 0)
			return nil, fmt.Errorf("%s", errMsg.String())
		}

		cdResult := gjson.GetBytes(respBytes, "data.cooldown.total_seconds")
		cd := int(cdResult.Int())

		cache.APILog(method, path, string(body), string(respBytes),
			resp.StatusCode, cd)
		if cd > 0 {
			console.Errorf("⏳ Cooldown started: %d seconds\n", cd)
			time.Sleep(time.Duration(cd) * time.Second)
		}

		return respBytes, nil
	}
}
