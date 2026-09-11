package api

import (
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
)

func TestAccountsAchievementsUsesCachedAccount(t *testing.T) {
	cache.CleanAccount()
	cache.SaveAccount("cached")
	t.Cleanup(cache.CleanAccount)

	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = testClient(func(r *http.Request) ([]byte, error) {
		if r.URL.Path != "/accounts/cached/achievements" {
			t.Fatalf("path = %s, want cached account path", r.URL.Path)
		}
		return []byte(`{"data":[],"pages":1}`), nil
	})

	got, err := AccountsAchievements("")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d achievements, want 0", len(got))
	}
}

func TestAccountsAchievementsLoadsAccountFromDetails(t *testing.T) {
	cache.CleanAccount()
	t.Cleanup(cache.CleanAccount)

	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = testClient(func(r *http.Request) ([]byte, error) {
		switch r.URL.Path {
		case "/my/details":
			return []byte(`{"data":{"username":"loaded"}}`), nil
		case "/accounts/loaded/achievements":
			return []byte(`{"data":[],"pages":1}`), nil
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
			return nil, nil
		}
	})

	got, err := AccountsAchievements("")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d achievements, want 0", len(got))
	}
}

func TestAccountsAchievementsReturnsDetailsError(t *testing.T) {
	cache.CleanAccount()
	t.Cleanup(cache.CleanAccount)

	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"details unavailable"}}`))

	_, err := AccountsAchievements("")
	if err == nil {
		t.Fatal("expected details error")
	}
}
