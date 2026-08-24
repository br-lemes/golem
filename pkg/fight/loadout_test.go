package fight

import "testing"

func TestFromLoadoutLevelOneBaseStats(t *testing.T) {
	fighter := FromLoadout(1, nil, nil)
	if fighter.Stats.HP != 120 {
		t.Fatalf("level-one HP = %d, want 120", fighter.Stats.HP)
	}
	if fighter.Stats.Initiative != 100 {
		t.Fatalf("level-one initiative = %d, want 100", fighter.Stats.Initiative)
	}
}
