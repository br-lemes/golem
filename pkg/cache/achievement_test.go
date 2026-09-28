package cache

import "testing"

func TestAchievementCache(t *testing.T) {
	initializeTestCache(t)
	SaveAchievements([]string{"second", "third"})
	if !HasAchievement("third") {
		t.Fatal("third achievement was not saved")
	}
	if HasAchievement("first") {
		t.Fatal("stale achievement was found")
	}
	if HasAchievement("missing") {
		t.Fatal("missing achievement was found")
	}
	CleanAchievements()
	if HasAchievement("first") {
		t.Fatal("achievement cache was not cleaned")
	}
}
