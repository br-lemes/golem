package cache

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
)

func TestUsageCombatCache(t *testing.T) {
	initializeTestCache(t)
	want := models.UsageCombat{Key: "combat", Version: 2, Results: "{}"}
	SaveUsageCombat(want)
	got, found := GetUsageCombat("combat", 2)
	if !found || got.Results != want.Results {
		t.Fatalf("GetUsageCombat() = %#v, %t", got, found)
	}
	_, found = GetUsageCombat("missing", 2)
	if found {
		t.Fatal("missing usage combat was found")
	}
	_, found = GetUsageCombat("combat", 1)
	if found {
		t.Fatal("usage combat with wrong version was found")
	}
}
