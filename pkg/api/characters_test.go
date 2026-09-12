package api

import (
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCharactersWildcardRequiresEnvironmentVariable(t *testing.T) {
	t.Setenv("GOLEM_NAME", "")
	client := &Client{}
	_, err := client.Characters(".")
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
	client := newTestClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		return testResponse(http.StatusOK, nil), nil
	}))

	got, err := client.Characters("gandalf")
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
	client := newTestClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/characters/gandalf" {
			t.Errorf("request path = %q", r.URL.Path)
		}
		return testResponse(http.StatusOK, []byte(`{"data":{"name":"gandalf","account":"account"},"status":"success"}`)), nil
	}))

	got, err := client.Characters("gandalf")
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
	client := newTestClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return testResponse(http.StatusNotFound, []byte(`{"error":{"message":"not found"}}`)), nil
	}))

	_, err := client.Characters("gandalf")
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestCharactersReturnsJSONError(t *testing.T) {
	cache.CleanCharacter("gandalf")
	client := newTestClient(roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return testResponse(http.StatusOK, []byte(`{"data":,"status":"success"}`)), nil
	}))

	_, err := client.Characters("gandalf")
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
