package cmd

import "testing"

func TestBestGroupMonster(t *testing.T) {
	_, err := bestGroupMonster("wolf")
	if err != nil {
		t.Fatalf("valid monster rejected: %v", err)
	}
	_, err = bestGroupMonster("missing")
	if err == nil {
		t.Fatal("invalid monster accepted")
	}
}
