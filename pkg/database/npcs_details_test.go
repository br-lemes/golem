package database

import "testing"

func TestNPCsSmoke(t *testing.T) {
	if len(NpcsDetails.All()) == 0 {
		t.Fatal("NPC details catalog is empty")
	}
	if len(NpcsItems.All()) == 0 {
		t.Fatal("NPC items catalog is empty")
	}
}
