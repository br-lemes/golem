package database

import "testing"

func TestEventsSmoke(t *testing.T) {
	if len(Events.All()) == 0 {
		t.Fatal("events catalog is empty")
	}
}
