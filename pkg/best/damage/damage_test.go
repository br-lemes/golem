package damage

import "testing"

func TestElementDamageRoundsAndClamps(t *testing.T) {
	rounded := elementDamage(6, 0, -30)
	if rounded != 8 {
		t.Fatalf("damage = %d, want 8", rounded)
	}
	clamped := elementDamage(16, 0, 115)
	if clamped != 0 {
		t.Fatalf("damage = %d, want 0", clamped)
	}
	documented := elementDamage(100, 30, 30)
	if documented != 91 {
		t.Fatalf("damage = %d, want 91", documented)
	}
}
