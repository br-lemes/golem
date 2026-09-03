package database

import "testing"

func TestTasksSmoke(t *testing.T) {
	if len(Tasks().All()) == 0 {
		t.Fatal("tasks catalog is empty")
	}
}
