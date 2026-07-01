package database

import "testing"

func TestMonstersSmoke(t *testing.T) {
	if len(GetMonsters()) == 0 {
		t.Error("monsters list is empty")
	}
	if len(GetMonsterCodes()) == 0 {
		t.Error("monster codes are empty")
	}
}
