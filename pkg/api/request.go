package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/logs"
	"github.com/tidwall/gjson"
)

var defaultClient *Client

type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	InitialWait time.Duration
	MaxRetries  int
	MaxWait     time.Duration
	Token       string
}

type request struct {
	body    []byte
	client  *Client
	method  string
	path    string
	retries int
	wait    time.Duration
}

func NewClient(apiConfig config.API) *Client {
	baseURL := "https://api.artifactsmmo.com"
	switch apiConfig.Environment {
	case "sandbox":
		baseURL = "https://api.sandbox.artifactsmmo.com"
	case "beta":
		baseURL = "https://api.beta.artifactsmmo.com"
	}
	return &Client{
		BaseURL:     baseURL,
		HTTPClient:  &http.Client{Timeout: 5 * time.Minute},
		InitialWait: 5 * time.Second,
		MaxWait:     60 * time.Second,
		Token:       apiConfig.Token,
	}
}

func Request(method, path string, body []byte) ([]byte, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.Request(method, path, body)
}

func (c *Client) Request(method, path string, body []byte) ([]byte, error) {
	r := &request{
		body:   body,
		client: c,
		method: method,
		path:   path,
		wait:   c.InitialWait,
	}
	return r.do()
}

func (r *request) do() ([]byte, error) {
	for {
		if !r.shouldRetry() {
			return nil, fmt.Errorf("max retries reached")
		}
		req, err := r.client.newRequest(r.method, r.path, r.body)
		if err != nil {
			logs.Record(logs.Event{
				Method:  r.method,
				Path:    r.path,
				Body:    string(r.body),
				Message: err.Error(),
				Status:  0,
			})
			return nil, err
		}

		resp, err := r.client.HTTPClient.Do(req)
		if err != nil {
			r.handleInfraError("Network error", err)
			continue
		}

		respBytes, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			r.handleInfraError("Read error", err)
			continue
		}

		if resp.StatusCode >= 400 && (resp.StatusCode < 500 || resp.StatusCode > 504) {
			return nil, r.handleClientError(resp, respBytes)
		}

		if resp.StatusCode >= 500 && resp.StatusCode <= 504 {
			r.wait = r.handleBackoff(resp, respBytes)
			continue
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			r.wait = r.handleBackoff(resp, respBytes)
			continue
		}

		errMsg := gjson.GetBytes(respBytes, "error.message")
		if errMsg.Exists() {
			logs.Record(logs.Event{
				Method:   r.method,
				Path:     r.path,
				Body:     string(r.body),
				Response: string(respBytes),
				Status:   resp.StatusCode,
			})
			return nil, fmt.Errorf("%s", errMsg.String())
		}

		logs.Record(logs.Event{
			Method:   r.method,
			Path:     r.path,
			Body:     string(r.body),
			Response: string(respBytes),
			Status:   resp.StatusCode,
		})

		return respBytes, nil
	}
}

func (c *Client) newRequest(method, path string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	console.Debugf("%s %s\n", method, path)
	if len(body) > 0 {
		console.Debugf("  Body: %s\n", string(body))
	}

	return req, nil
}

func (r *request) handleInfraError(reason string, err error) {
	format := "%s: %v. Retrying in %v...\n"
	message := fmt.Sprintf(format, reason, err, r.wait)
	logs.Record(logs.Event{
		Method:  r.method,
		Path:    r.path,
		Body:    string(r.body),
		Level:   logs.Error,
		Message: message,
	})
	time.Sleep(r.wait)
}

func (r *request) handleClientError(resp *http.Response, respBytes []byte) error {
	errMsg := gjson.GetBytes(respBytes, "error.message")
	message := "Client error: "
	if errMsg.Exists() {
		message += errMsg.String()
	} else {
		message += resp.Status
	}
	logs.Record(logs.Event{
		Method:   r.method,
		Path:     r.path,
		Body:     string(r.body),
		Response: string(respBytes),
		Status:   resp.StatusCode,
	})
	return fmt.Errorf("%s (Status: %d)", message, resp.StatusCode)
}

func (r *request) handleBackoff(resp *http.Response, respBytes []byte) time.Duration {
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
	logMessage := fmt.Sprintf(format, message, resp.StatusCode, r.wait)
	logs.Record(logs.Event{
		Method:   r.method,
		Path:     r.path,
		Body:     string(r.body),
		Response: string(respBytes),
		Status:   resp.StatusCode,
		Level:    logs.Error,
		Message:  logMessage,
	})
	time.Sleep(r.wait)

	nextWait := r.wait * 2
	if nextWait > r.client.MaxWait {
		return r.client.MaxWait
	}
	return nextWait
}

func (r *request) shouldRetry() bool {
	if r.client.MaxRetries <= 0 {
		return true
	}
	if r.retries >= r.client.MaxRetries {
		return false
	}
	r.retries++
	return true
}
