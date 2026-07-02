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

type requestCtx struct {
	baseURL    string
	method     string
	path       string
	body       []byte
	maxRetries int
	retries    int
}

func Request(method, path string, body []byte, cooldown bool) ([]byte, error) {
	ctx := &requestCtx{baseURL: baseURL, method: method, path: path, body: body}
	return ctx.execute(cooldown)
}

func (ctx *requestCtx) execute(cooldown bool) ([]byte, error) {
	serverWait := 5 * time.Second
	maxServerWait := 60 * time.Second

	for {
		if !ctx.shouldRetry() {
			return nil, fmt.Errorf("max retries reached")
		}
		req, err := ctx.newRequest()
		if err != nil {
			cache.APILog(ctx.method, ctx.path, string(ctx.body), err.Error(), 0,
				0)
			return nil, err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			format := "Network error: %v. Retrying in 5 seconds...\n"
			ctx.handleInfraError(format, err)
			continue
		}

		respBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			format := "Read error: %v. Retrying in 5 seconds...\n"
			ctx.handleInfraError(format, err)
			continue
		}

		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, ctx.handleClientError(resp, respBytes)
		}

		if resp.StatusCode >= 500 {
			serverWait = ctx.handleBackoff(resp, respBytes, serverWait,
				maxServerWait)
			continue
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			serverWait = ctx.handleBackoff(resp, respBytes, serverWait,
				maxServerWait)
			continue
		}

		errMsg := gjson.GetBytes(respBytes, "error.message")
		if errMsg.Exists() {
			cache.APILog(ctx.method, ctx.path, string(ctx.body),
				string(respBytes), resp.StatusCode, 0)
			return nil, fmt.Errorf("%s", errMsg.String())
		}

		cdResult := gjson.GetBytes(respBytes, "data.cooldown.total_seconds")
		cd := int(cdResult.Int())

		cache.APILog(ctx.method, ctx.path, string(ctx.body), string(respBytes),
			resp.StatusCode, cd)

		if cooldown && cd > 0 {
			handleCooldown(cd)
		}

		return respBytes, nil
	}
}

func (ctx *requestCtx) newRequest() (*http.Request, error) {
	req, err := http.NewRequest(ctx.method, ctx.baseURL+ctx.path,
		bytes.NewReader(ctx.body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Token)

	console.Debugf("%s %s\n", ctx.method, ctx.path)
	if len(ctx.body) > 0 {
		console.Debugf("  Body: %s\n", string(ctx.body))
	}

	return req, nil
}

func (ctx *requestCtx) handleInfraError(format string, err error) {
	message := fmt.Sprintf(format, err)
	cache.APILog(ctx.method, ctx.path, string(ctx.body), message, 0, 0)
	console.Errorf(format, err)
	time.Sleep(5 * time.Second)
}

func (ctx *requestCtx) handleClientError(resp *http.Response, respBytes []byte) error {
	errMsg := gjson.GetBytes(respBytes, "error.message")
	message := "Client error: "
	if errMsg.Exists() {
		message += errMsg.String()
	} else {
		message += resp.Status
	}
	cache.APILog(ctx.method, ctx.path, string(ctx.body), string(respBytes),
		resp.StatusCode, 0)
	return fmt.Errorf("%s (Status: %d)", message, resp.StatusCode)
}

func (ctx *requestCtx) handleBackoff(resp *http.Response, respBytes []byte, currentWait, maxWait time.Duration) time.Duration {
	errMsg := gjson.GetBytes(respBytes, "error.message")
	message := "Server error: "
	if errMsg.Exists() {
		message += errMsg.String()
	} else {
		message += resp.Status
	}

	format := "%s (Status: %d). Retrying in %v...\n"
	logMessage := fmt.Sprintf(format, message, resp.StatusCode, currentWait)
	cache.APILog(ctx.method, ctx.path, string(ctx.body), logMessage,
		resp.StatusCode, 0)
	console.Errorf(format, message, resp.StatusCode, currentWait)
	time.Sleep(currentWait)

	nextWait := currentWait * 2
	if nextWait > maxWait {
		return maxWait
	}
	return nextWait
}

func (ctx *requestCtx) shouldRetry() bool {
	if ctx.maxRetries <= 0 {
		return true
	}
	if ctx.retries >= ctx.maxRetries {
		return false
	}
	ctx.retries++
	return true
}
