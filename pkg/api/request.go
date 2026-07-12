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

var (
	token         string
	defaultClient = &http.Client{Timeout: 30 * time.Second}
)

type requestCtx struct {
	baseURL     string
	body        []byte
	client      *http.Client
	initialWait time.Duration
	maxRetries  int
	maxWait     time.Duration
	method      string
	path        string
	retries     int
}

func Request(method, path string, body []byte, cooldown bool) ([]byte, error) {
	ctx := &requestCtx{
		baseURL:     baseURL,
		body:        body,
		client:      defaultClient,
		initialWait: 5 * time.Second,
		maxWait:     60 * time.Second,
		method:      method,
		path:        path,
	}
	return ctx.execute(cooldown)
}

func (ctx *requestCtx) execute(cooldown bool) ([]byte, error) {
	wait := ctx.initialWait
	maxWait := ctx.maxWait

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

		resp, err := ctx.client.Do(req)
		if err != nil {
			ctx.handleInfraError("Network error", err)
			continue
		}

		respBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			ctx.handleInfraError("Read error", err)
			continue
		}

		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, ctx.handleClientError(resp, respBytes)
		}

		if resp.StatusCode >= 500 {
			wait = ctx.handleBackoff(resp, respBytes, wait, maxWait)
			continue
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			wait = ctx.handleBackoff(resp, respBytes, wait, maxWait)
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
	req.Header.Set("Authorization", "Bearer "+token)

	console.Debugf("%s %s\n", ctx.method, ctx.path)
	if len(ctx.body) > 0 {
		console.Debugf("  Body: %s\n", string(ctx.body))
	}

	return req, nil
}

func (ctx *requestCtx) handleInfraError(reason string, err error) {
	format := "%s: %v. Retrying in %v...\n"
	message := fmt.Sprintf(format, reason, err, ctx.initialWait)
	cache.APILog(ctx.method, ctx.path, string(ctx.body), message, 0, 0)
	console.Errorf(format, reason, err, ctx.initialWait)
	time.Sleep(ctx.initialWait)
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
	if resp.StatusCode < 400 {
		message = "Invalid response format: "
	}
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
