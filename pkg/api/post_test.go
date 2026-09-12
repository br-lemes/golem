package api

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
)

type failingJSON struct{}

func (failingJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("cannot marshal payload")
}

func TestPostWrapperReturnsResponseWithoutCooldown(t *testing.T) {
	previousClient := defaultClient
	t.Cleanup(func() { defaultClient = previousClient })
	defaultClient = newTestClient(responseTransport(http.StatusOK, []byte(`{"data":{"cooldown":{"total_seconds":0}}}`)))

	got, err := Post("/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"data":{"cooldown":{"total_seconds":0}}}` {
		t.Fatalf("response = %s", got)
	}
}

func TestPostWrapperReturnsRequestError(t *testing.T) {
	previousClient := defaultClient
	t.Cleanup(func() { defaultClient = previousClient })
	defaultClient = newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"post failed"}}`)))

	_, err := Post("/test", nil)
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestPostSendsPayloadAndReturnsResponse(t *testing.T) {
	client := newTestClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("path = %s, want /test", r.URL.Path)
		}
		if r.Body == nil {
			t.Fatal("POST body is nil")
		}
		body := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(body)
		if !bytes.Equal(body, []byte(`{"value":"test"}`)) {
			t.Errorf("body = %s, want JSON payload", body)
		}
		return testResponse(http.StatusOK, []byte(`{"data":{"cooldown":{"total_seconds":0}}}`)), nil
	}))

	got, err := client.Post("/test", map[string]string{"value": "test"})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"data":{"cooldown":{"total_seconds":0}}}` {
		t.Fatalf("response = %s", got)
	}
}

func TestPostAcceptsNilPayload(t *testing.T) {
	client := newTestClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Body == nil {
			t.Fatal("POST body is nil")
		}
		return testResponse(http.StatusOK, []byte(`{"data":{"cooldown":{"total_seconds":0}}}`)), nil
	}))

	_, err := client.Post("/test", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostReturnsRequestError(t *testing.T) {
	client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"post failed"}}`)))

	_, err := client.Post("/test", nil)
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestPostReturnsMarshalError(t *testing.T) {
	client := newTestClient(http.DefaultTransport)
	_, err := client.Post("/test", failingJSON{})
	if err == nil {
		t.Fatal("expected marshal error")
	}
}
