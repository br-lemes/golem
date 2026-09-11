package api

import (
	"net/http"
	"testing"
)

func TestMyDetailsReturnsJSONError(t *testing.T) {
	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusOK, []byte("invalid json"))

	_, err := MyDetails()
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
