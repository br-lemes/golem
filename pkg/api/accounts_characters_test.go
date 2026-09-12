package api

import (
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestAccountsCharactersUsesCachedCharacters(t *testing.T) {
	cleanAccountCharactersCache(t)
	cache.SaveAccount("account")
	characters := []schemas.CharacterSchema{
		{Account: "account", Name: "cached"},
	}
	cache.SaveCharacters(characters)

	client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"request not expected"}}`)))

	got, err := client.AccountsCharacters("")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 || got[0].Name != "cached" {
		t.Fatalf("got cached characters %#v", got)
	}
}

func TestAccountsCharactersFetchesExplicitAccount(t *testing.T) {
	cleanAccountCharactersCache(t)
	cache.SaveAccount("account")

	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		if r.URL.Path != "/accounts/other/characters" {
			t.Fatalf("path = %s, want other account path", r.URL.Path)
		}
		return []byte(`{"data":[{"name":"other"}]}`), nil
	}))

	got, err := client.AccountsCharacters("other")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "other" {
		t.Fatalf("got characters %#v", got)
	}
}

func TestAccountsCharactersLoadsAccountFromDetails(t *testing.T) {
	cleanAccountCharactersCache(t)

	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		switch r.URL.Path {
		case "/my/details":
			return []byte(`{"data":{"username":"loaded"}}`), nil
		case "/accounts/loaded/characters":
			return []byte(`{"data":[]}`), nil
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
			return nil, nil
		}
	}))

	got, err := client.AccountsCharacters("")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d characters, want 0", len(got))
	}
}

func TestAccountsCharactersReturnsDetailsError(t *testing.T) {
	cleanAccountCharactersCache(t)
	client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"details unavailable"}}`)))

	_, err := client.AccountsCharacters("")
	if err == nil {
		t.Fatal("expected details error")
	}
}

func TestAccountsCharactersReturnsRequestError(t *testing.T) {
	cleanAccountCharactersCache(t)
	cache.SaveAccount("account")
	client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"characters unavailable"}}`)))

	_, err := client.AccountsCharacters("other")
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestAccountsCharactersReturnsJSONError(t *testing.T) {
	cleanAccountCharactersCache(t)
	cache.SaveAccount("account")
	client := newTestClient(responseTransport(http.StatusOK, []byte("invalid json")))

	_, err := client.AccountsCharacters("other")
	if err == nil {
		t.Fatal("expected JSON error")
	}
}

func cleanAccountCharactersCache(t *testing.T) {
	t.Helper()
	cache.CleanAccount()
	cache.CleanCharacters()
	t.Cleanup(func() {
		cache.CleanAccount()
		cache.CleanCharacters()
	})
}
