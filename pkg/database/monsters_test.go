package database

import "testing"

func TestMonstersSmoke(t *testing.T) {
	if len(Monsters.All()) == 0 {
		t.Fatal("monsters catalog is empty")
	}
	if len(Bosses.All()) == 0 {
		t.Fatal("bosses view is empty")
	}
}
