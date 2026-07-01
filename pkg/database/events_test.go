package database

import "testing"

func TestEventsSmoke(t *testing.T) {
	if len(GetEvents()) == 0 {
		t.Error("events list is empty")
	}
	if len(GetEventCodes()) == 0 {
		t.Error("event codes are empty")
	}
	if len(GetEventContentCodes()) == 0 {
		t.Error("event content codes are empty")
	}
}
