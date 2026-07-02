package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/tidwall/gjson"
)

const baseURL = "https://api.artifactsmmo.com"

var Token string

var httpClient = &http.Client{Timeout: 30 * time.Second}

func Request(method, path string, body []byte, cooldown bool) ([]byte, error) {
	serverWait := 5 * time.Second
	maxServerWait := 60 * time.Second

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

		if resp.StatusCode >= http.StatusInternalServerError {
			format := "Server error: %d. Retrying in %v...\n"
			message := fmt.Sprintf(format, resp.StatusCode, serverWait)
			cache.APILog(method, path, string(body), message, resp.StatusCode, 0)
			console.Errorf(format, resp.StatusCode, serverWait)
			time.Sleep(serverWait)
			serverWait = serverWait * 2
			if serverWait > maxServerWait {
				serverWait = maxServerWait
			}
			continue
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			format := "Unexpected non-JSON response: %d. Retrying in %v...\n"
			message := fmt.Sprintf(format, resp.StatusCode, serverWait)
			cache.APILog(method, path, string(body), message, resp.StatusCode, 0)
			console.Errorf(format, resp.StatusCode, serverWait)
			time.Sleep(serverWait)
			serverWait = serverWait * 2
			if serverWait > maxServerWait {
				serverWait = maxServerWait
			}
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

		if cooldown && cd > 0 {
			handleCooldown(cd)
		}

		return respBytes, nil
	}
}
