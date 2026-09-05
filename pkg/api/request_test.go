package api

import (
	"errors"
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
	database := config.Database{Driver: "sqlite", Path: directory + "/cache.db"}
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

	ctx := &requestCtx{
		baseURL:     server.URL,
		body:        []byte(`{}`),
		client:      server.Client(),
		initialWait: 1 * time.Millisecond,
		maxWait:     8 * time.Millisecond,
		method:      http.MethodPost,
		path:        "/test",
	}

	resp, err := ctx.execute(false)
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

	ctx := &requestCtx{
		baseURL:     server.URL,
		body:        []byte(`{}`),
		client:      server.Client(),
		initialWait: 1 * time.Millisecond,
		maxWait:     8 * time.Millisecond,
		method:      http.MethodPost,
		path:        "/test",
	}

	_, err := ctx.execute(false)
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

	ctx := &requestCtx{
		baseURL:     server.URL,
		body:        []byte(`{}`),
		client:      server.Client(),
		initialWait: 1 * time.Millisecond,
		maxRetries:  2,
		maxWait:     8 * time.Millisecond,
		method:      http.MethodGet,
		path:        "/test",
	}

	_, err := ctx.execute(false)
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

	ctx := &requestCtx{
		baseURL:     server.URL,
		body:        []byte(`{}`),
		client:      server.Client(),
		initialWait: 1 * time.Millisecond,
		maxRetries:  1,
		maxWait:     8 * time.Millisecond,
		method:      http.MethodGet,
		path:        "/test",
	}

	_, err := ctx.execute(false)
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

	ctx := &requestCtx{
		baseURL:     server.URL,
		client:      server.Client(),
		method:      http.MethodPost,
		initialWait: 1 * time.Millisecond,
	}

	_, err := ctx.execute(true)
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

	ctx := &requestCtx{
		baseURL: server.URL,
		client:  server.Client(),
		method:  http.MethodGet,
	}

	_, err := ctx.execute(false)
	expectedMsg := "Client error: 418 I'm a teapot (Status: 418)"
	if err == nil || err.Error() != expectedMsg {
		t.Errorf("expected %q, got %v", expectedMsg, err)
	}
}

func TestExecuteNetworkErrorRetry(t *testing.T) {
	customClient := &http.Client{}
	customClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network connection reset")
	})

	ctx := &requestCtx{
		baseURL:     "https://invalid-url-target.local",
		body:        []byte(`{}`),
		client:      customClient,
		initialWait: 1 * time.Millisecond,
		maxRetries:  1,
		maxWait:     8 * time.Millisecond,
		method:      http.MethodGet,
		path:        "/test",
	}

	_, err := ctx.execute(false)
	if err == nil {
		t.Fatal("expected max retries error, got nil")
	}

	expectedMsg := "max retries reached"
	if err.Error() != expectedMsg {
		t.Errorf("expected %q, got %q", expectedMsg, err.Error())
	}
}

func TestExecuteResponseWithEmbeddedError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"error": {"message": "character not found"}}`))
	}))
	defer server.Close()

	ctx := &requestCtx{
		baseURL:     server.URL,
		body:        []byte(`{}`),
		client:      server.Client(),
		initialWait: 1 * time.Millisecond,
		maxWait:     8 * time.Millisecond,
		method:      http.MethodPost,
		path:        "/test",
	}

	_, err := ctx.execute(false)
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

	ctx := &requestCtx{
		baseURL:     server.URL,
		body:        []byte(`{}`),
		client:      server.Client(),
		initialWait: 2 * time.Millisecond,
		maxRetries:  3,
		maxWait:     3 * time.Millisecond,
		method:      http.MethodGet,
		path:        "/test",
	}

	_, err := ctx.execute(false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	wait := ctx.initialWait
	nextWait := wait * 2
	if nextWait <= ctx.maxWait {
		t.Errorf("test parameters setup incorrectly")
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
