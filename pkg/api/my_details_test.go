package api

import (
	"net/http"
	"testing"
)

func TestMyDetailsReturnsJSONError(t *testing.T) {
	client := newTestClient(responseTransport(http.StatusOK, []byte("invalid json")))

	_, err := client.MyDetails()
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
