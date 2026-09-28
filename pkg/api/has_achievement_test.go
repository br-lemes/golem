package api

import (
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
)

func TestHasAchievementCachesCompletedAchievements(t *testing.T) {
	cache.CleanAccount()
	cache.CleanAchievements()
	cache.SaveAccount("account")
	t.Cleanup(func() {
		cache.CleanAccount()
		cache.CleanAchievements()
	})
	requests := 0
	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		requests++
		if r.URL.Path != "/accounts/account/achievements" {
			t.Fatalf("path = %s, want account achievements", r.URL.Path)
		}
		if r.URL.Query().Get("completed") != "true" {
			t.Fatalf("completed = %q, want true", r.URL.Query().Get("completed"))
		}
		switch r.URL.Query().Get("page") {
		case "1":
			return []byte(`{"data":[{"code":"first","completed_at":"2026-01-01T00:00:00Z"}],"pages":2}`), nil
		case "2":
			return []byte(`{"data":[{"code":"second","completed_at":"2026-01-01T00:00:00Z"}],"pages":2}`), nil
		default:
			t.Fatalf("page = %q, want 1 or 2", r.URL.Query().Get("page"))
			return nil, nil
		}
	}))

	unlocked, err := client.HasAchievement("second")
	if err != nil {
		t.Fatal(err)
	}
	if !unlocked {
		t.Fatal("HasAchievement() = false, want true")
	}
	if requests != 2 {
		t.Fatalf("requests = %d, want 2", requests)
	}

	unlocked, err = client.HasAchievement("first")
	if err != nil {
		t.Fatal(err)
	}
	if !unlocked {
		t.Fatal("cached HasAchievement() = false, want true")
	}
	if requests != 2 {
		t.Fatalf("requests = %d after cache hit, want 2", requests)
	}
}

func TestHasAchievementReturnsFalseForUncompletedAchievement(t *testing.T) {
	cache.CleanAccount()
	cache.CleanAchievements()
	cache.SaveAccount("account")
	t.Cleanup(func() {
		cache.CleanAccount()
		cache.CleanAchievements()
	})
	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		return []byte(`{"data":[{"code":"completed","completed_at":"2026-01-01T00:00:00Z"}],"pages":1}`), nil
	}))

	unlocked, err := client.HasAchievement("uncompleted")
	if err != nil {
		t.Fatal(err)
	}
	if unlocked {
		t.Fatal("HasAchievement() = true, want false")
	}
}
