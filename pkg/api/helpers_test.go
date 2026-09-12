package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func newTestClient(transport http.RoundTripper) *Client {
	return &Client{
		BaseURL:     "https://api.test",
		HTTPClient:  &http.Client{Transport: transport},
		InitialWait: time.Millisecond,
		MaxRetries:  1,
		MaxWait:     8 * time.Millisecond,
	}
}

func testPaginatedEndpoint(t *testing.T, path string, size int, call func(*Client) (int, error)) {
	t.Helper()
	for _, test := range []struct {
		name      string
		transport http.RoundTripper
		wantErr   bool
	}{
		{
			name: "valid pagination",
			transport: func() http.RoundTripper {
				requests := 0
				return testTransport(func(r *http.Request) ([]byte, error) {
					requests++
					if r.URL.Path != path {
						t.Errorf("path = %s, want %s", r.URL.Path, path)
					}
					queryPage(t, r.URL.Query(), requests, size)
					return []byte(`{"data":[{},{}],"pages":2}`), nil
				})
			}(),
		},
		{
			name:      "HTTP error",
			transport: responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"bad request"}}`)),
			wantErr:   true,
		},
		{
			name:      "JSON error",
			transport: responseTransport(http.StatusOK, []byte("invalid json")),
			wantErr:   true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := newTestClient(test.transport)
			got, err := call(client)
			if test.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != 4 {
				t.Fatalf("got %d items, want 4", got)
			}
		})
	}
}

func queryPage(t *testing.T, values url.Values, page, size int) {
	t.Helper()
	if values.Get("page") != fmt.Sprint(page) {
		t.Errorf("page = %q, want %d", values.Get("page"), page)
	}
	if values.Get("size") != fmt.Sprint(size) {
		t.Errorf("size = %q, want %d", values.Get("size"), size)
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testTransport(handler func(*http.Request) ([]byte, error)) http.RoundTripper {
	return roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body, err := handler(r)
		if err != nil {
			return nil, err
		}
		return testResponse(http.StatusOK, body), nil
	})
}

func responseTransport(status int, body []byte) http.RoundTripper {
	return roundTripFunc(func(*http.Request) (*http.Response, error) {
		return testResponse(status, body), nil
	})
}

func testResponse(status int, body []byte) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		StatusCode: status,
	}
}

func prepareBankActionCaches(t *testing.T) {
	t.Helper()
	cache.CleanBank()
	cache.CleanBankItems()
	cache.CleanCharacters()
	t.Cleanup(cache.CleanBank)
	t.Cleanup(cache.CleanBankItems)
	t.Cleanup(cache.CleanCharacters)
	cache.SaveBank(schemas.BankSchema{})
}

func bankActionClient(t *testing.T, path string, requestBody, responseBody []byte) *Client {
	t.Helper()
	prepareBankActionCaches(t)
	return newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		if r.Method != http.MethodPost || r.URL.Path != path {
			t.Errorf("request = %s %s, want POST %s", r.Method, r.URL.Path, path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(body, requestBody) {
			t.Errorf("body = %s, want %s", body, requestBody)
		}
		return responseBody, nil
	}))
}

type actionTestCase struct {
	name string
	call func(*Client) error
}

type actionBodyTestCase struct {
	name string
	path string
	body []byte
	call func(*Client) error
}

func testActionErrors(t *testing.T, tests []actionTestCase) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.name+" request error", func(t *testing.T) {
			client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"action failed"}}`)))
			err := test.call(client)
			if err == nil {
				t.Fatal("expected request error")
			}
		})
		t.Run(test.name+" JSON error", func(t *testing.T) {
			client := newTestClient(responseTransport(http.StatusOK, []byte("invalid json")))
			err := test.call(client)
			if err == nil {
				t.Fatal("expected JSON error")
			}
		})
	}
}

func testActionBody(t *testing.T, test actionBodyTestCase) {
	t.Helper()
	cache.CleanCharacters()
	t.Cleanup(cache.CleanCharacters)
	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		if r.Method != http.MethodPost || r.URL.Path != test.path {
			t.Errorf("request = %s %s, want POST %s", r.Method, r.URL.Path, test.path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(body, test.body) {
			t.Errorf("body = %s, want %s", body, test.body)
		}
		return []byte(`{"data":{}}`), nil
	}))
	err := test.call(client)
	if err != nil {
		t.Fatal(err)
	}
}
