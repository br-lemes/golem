package logs

import (
	"testing"
	"time"

	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/models"
)

func initializeTestLogs(t *testing.T) {
	t.Helper()
	database = nil
	databaseDay = ""
	err := Initialize(config.Storage{Logs: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRecordStoresEvent(t *testing.T) {
	initializeTestLogs(t)
	event := Event{Method: "GET", Path: "/test", Body: "{}"}
	event.Response = "ok"
	event.Status = 200
	Record(event)
	var request models.Request
	err := database.First(&request).Error
	if err != nil {
		t.Fatal(err)
	}
	if request.Method != "GET" || request.Path != "/test" || request.Response != "ok" {
		t.Fatalf("recorded request = %#v", request)
	}
}

func TestRecordUsesMessageWhenResponseIsEmpty(t *testing.T) {
	initializeTestLogs(t)
	Record(Event{Method: "POST", Message: "failed", Level: Error})
	var request models.Request
	err := database.First(&request).Error
	if err != nil {
		t.Fatal(err)
	}
	if request.Response != "failed" {
		t.Fatalf("fallback response = %q", request.Response)
	}
}

func TestRecordHandlesDebugLevel(t *testing.T) {
	initializeTestLogs(t)
	Record(Event{Message: "debug message", Level: Debug})
}

func TestOpenDatabaseReusesDatabaseForSameDay(t *testing.T) {
	initializeTestLogs(t)
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	err := openDatabase(now)
	if err != nil {
		t.Fatal(err)
	}
	first := database
	err = openDatabase(now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if database != first {
		t.Fatal("same-day database was reopened")
	}
}

func TestOpenDatabaseRotatesDatabaseOnNewDay(t *testing.T) {
	initializeTestLogs(t)
	now := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	err := openDatabase(now)
	if err != nil {
		t.Fatal(err)
	}
	first := database
	err = openDatabase(now.Add(24 * time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if database == first || databaseDay != "2026-01-02" {
		t.Fatal("database was not rotated on a new day")
	}
}
