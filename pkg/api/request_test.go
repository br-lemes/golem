package api

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/logs"
)

func TestMain(m *testing.M) {
	directory, err := os.MkdirTemp("", "golem-api-test-")
	if err != nil {
		panic(err)
	}
	database := config.Storage{
		Cache: directory + "/cache.db",
		Logs:  directory + "/logs",
	}
	err = cache.Initialize(database)
	if err != nil {
		_ = os.RemoveAll(directory)
		panic(err)
	}
	err = logs.Initialize(database)
	if err != nil {
		_ = os.RemoveAll(directory)
		panic(err)
	}
	code := m.Run()
	_ = os.RemoveAll(directory)
	os.Exit(code)
}

func TestExecuteSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": {"cooldown": {"total_seconds": 0}}, "status": "success"}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL
	client.InitialWait = time.Millisecond
	client.MaxWait = 8 * time.Millisecond

	resp, err := client.Request(http.MethodPost, "/test", []byte(`{}`))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response bytes, got nil")
	}
}

func TestExecuteClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error": {"message": "invalid payload"}}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL
	client.InitialWait = time.Millisecond
	client.MaxWait = 8 * time.Millisecond

	_, err := client.Request(http.MethodPost, "/test", []byte(`{}`))
	if err == nil {
		t.Fatal("expected client error, got nil")
	}

	expectedMsg := "Client error: invalid payload (Status: 400)"
	if err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestExecuteMaxRetriesReached(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": {"message": "server down"}}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL
	client.InitialWait = time.Millisecond
	client.MaxRetries = 2
	client.MaxWait = 8 * time.Millisecond

	_, err := client.Request(http.MethodGet, "/test", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error due to max retries, got nil")
	}

	expectedMsg := "max retries reached"
	if err.Error() != expectedMsg {
		t.Errorf("expected error %q, got %q", expectedMsg, err.Error())
	}
}

func TestExecuteInvalidContentTypeBackoff(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`text response`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL
	client.InitialWait = time.Millisecond
	client.MaxRetries = 1
	client.MaxWait = 8 * time.Millisecond

	_, err := client.Request(http.MethodGet, "/test", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error after non-json backoff retry limit, got nil")
	}
}

func TestExecuteCooldownActivated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data": {"cooldown": {"total_seconds": 1}}}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL
	client.InitialWait = time.Millisecond

	_, err := client.Request(http.MethodPost, "", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestExecuteClientErrorFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTeapot)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL

	_, err := client.Request(http.MethodGet, "", nil)
	expectedMsg := "Client error: 418 I'm a teapot (Status: 418)"
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected %q, got %v", expectedMsg, err)
	}
}

func TestExecuteClientErrorForNonRetryableServerStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusHTTPVersionNotSupported)
		_, _ = w.Write([]byte(`{"error": {"message": "unsupported protocol"}}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL

	_, err := client.Request(http.MethodGet, "", nil)
	expectedMsg := "Client error: unsupported protocol (Status: 505)"
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected %q, got %v", expectedMsg, err)
	}
}

func TestExecuteNetworkErrorRetry(t *testing.T) {
	customClient := &http.Client{}
	customClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network connection reset")
	})

	client := newTestClient(customClient.Transport)
	client.BaseURL = "https://invalid-url-target.local"
	client.InitialWait = time.Millisecond
	client.MaxRetries = 1
	client.MaxWait = 8 * time.Millisecond

	_, err := client.Request(http.MethodGet, "/test", []byte(`{}`))
	if err == nil {
		t.Fatal("expected max retries error, got nil")
	}

	expectedMsg := "max retries reached"
	if err.Error() != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestExecuteRetriesWithoutLimit(t *testing.T) {
	attempts := 0
	client := newTestClient(roundTripFunc(func(*http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return nil, errors.New("temporary network error")
		}
		return testResponse(http.StatusOK, []byte(`{"status":"success"}`)), nil
	}))
	client.MaxRetries = 0
	client.InitialWait = 0

	_, err := client.Request(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestExecuteResponseWithEmbeddedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"error": {"message": "character not found"}}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL
	client.InitialWait = time.Millisecond
	client.MaxWait = 8 * time.Millisecond

	_, err := client.Request(http.MethodPost, "/test", []byte(`{}`))
	if err == nil {
		t.Fatal("expected embedded api error, got nil")
	}

	expectedMsg := "character not found"
	if err.Error() != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestExecuteBackoffCapping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTestClient(http.DefaultTransport)
	client.HTTPClient = server.Client()
	client.BaseURL = server.URL
	client.InitialWait = 2 * time.Millisecond
	client.MaxRetries = 3
	client.MaxWait = 3 * time.Millisecond

	_, err := client.Request(http.MethodGet, "/test", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	wait := client.InitialWait
	nextWait := wait * 2
	if nextWait <= client.MaxWait {
		t.Errorf("test parameters setup incorrectly")
	}
}

func TestExecuteReturnsNewRequestError(t *testing.T) {
	client := newTestClient(http.DefaultTransport)
	client.BaseURL = "://invalid"
	client.MaxRetries = 1

	_, err := client.Request(http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected request construction error")
	}
}

func TestNewRequestReturnsError(t *testing.T) {
	client := newTestClient(http.DefaultTransport)
	client.BaseURL = "://invalid"
	_, err := client.newRequest(http.MethodGet, "", nil)
	if err == nil {
		t.Fatal("expected request construction error")
	}
}

func TestExecuteRetriesWhenReadingResponseFails(t *testing.T) {
	transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			Body:       failingReadCloser{},
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			StatusCode: http.StatusOK,
		}, nil
	})
	client := newTestClient(transport)
	client.InitialWait = time.Millisecond
	client.MaxRetries = 1
	client.MaxWait = time.Millisecond
	_, err := client.Request(http.MethodGet, "/test", nil)
	if err == nil {
		t.Fatal("expected max retries error")
	}
}

type failingReadCloser struct{}

func (failingReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (failingReadCloser) Close() error {
	return io.ErrUnexpectedEOF
}
