package cache

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
)

func TestUsageCombatCache(t *testing.T) {
	initializeTestCache(t)
	want := models.UsageCombat{
		Type:    UsageCombatType,
		Key:     "combat",
		Version: UsageCombatVersion,
		Results: "{}",
	}
	SaveUsageCombat(want)
	SaveUsageCombat(models.UsageCombat{
		Type:    PotentialCombatType,
		Key:     "combat",
		Version: UsageCombatVersion,
		Results: "[]",
	})
	got, found := GetUsageCombat(UsageCombatType, "combat")
	if !found || got.Results != want.Results {
		t.Fatalf("GetUsageCombat() = %#v, %t", got, found)
	}
	_, found = GetUsageCombat(PotentialCombatType, "combat")
	if !found {
		t.Fatal("potential usage combat was not found")
	}
	_, found = GetUsageCombat(UsageCombatType, "missing")
	if found {
		t.Fatal("missing usage combat was found")
	}
	SaveUsageCombat(models.UsageCombat{
		Type:    UsageCombatType,
		Key:     "old",
		Version: UsageCombatVersion - 1,
		Results: "{}",
	})
	_, found = GetUsageCombat(UsageCombatType, "old")
	if found {
		t.Fatal("usage combat with wrong version was found")
	}
}
