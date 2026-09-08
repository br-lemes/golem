package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCharactersWildcardRequiresEnvironmentVariable(t *testing.T) {
	t.Setenv("GOLEM_NAME", "")
	_, err := Characters(".")
	if err == nil {
		t.Fatal("expected GOLEM_NAME error")
	}
	if err.Error() != "GOLEM_NAME is not set" {
		t.Errorf("error = %q, want %q", err.Error(), "GOLEM_NAME is not set")
	}
}

func TestCharactersReturnsCachedCharacter(t *testing.T) {
	cache.SaveAccount("account")
	character := schemas.CharacterSchema{Name: "gandalf", Account: "account"}
	cache.SaveCharacter(character)
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer server.Close()
	previousURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = previousURL }()

	got, err := Characters("gandalf")
	if err != nil {
		t.Fatalf("get cached character: %v", err)
	}
	if got.Name != character.Name || called {
		t.Errorf("got %#v, API called = %v", got, called)
	}
}

func TestCharactersFetchesAndCachesCharacter(t *testing.T) {
	cache.SaveAccount("account")
	cache.CleanCharacter("gandalf")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/characters/gandalf" {
			t.Errorf("request path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"data":{"name":"gandalf","account":"account"},"status":"success"}`))
	}))
	defer server.Close()
	previousURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = previousURL }()

	got, err := Characters("gandalf")
	if err != nil {
		t.Fatalf("fetch character: %v", err)
	}
	if got.Name != "gandalf" {
		t.Errorf("character name = %q", got.Name)
	}
	if cache.GetCharacter("gandalf") == nil {
		t.Error("expected character to be cached")
	}
}

func TestCharactersReturnsRequestError(t *testing.T) {
	cache.CleanCharacter("gandalf")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"message":"not found"}}`))
	}))
	defer server.Close()
	previousURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = previousURL }()

	_, err := Characters("gandalf")
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestCharactersReturnsJSONError(t *testing.T) {
	cache.CleanCharacter("gandalf")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":,"status":"success"}`))
	}))
	defer server.Close()
	previousURL := baseURL
	baseURL = server.URL
	defer func() { baseURL = previousURL }()

	_, err := Characters("gandalf")
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
