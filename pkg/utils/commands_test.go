package utils

import "testing"

func TestCommandsAndCompletion(t *testing.T) {
	commands, err := GetCommands()
	if err != nil || len(commands) == 0 {
		t.Fatalf("GetCommands() = %d commands, %v", len(commands), err)
	}
	for command, routes := range commands {
		if len(routes) == 0 {
			t.Fatalf("GetCommands returned no routes for %q", command)
		}
	}
	if len(GetCommandsCompletion()) != len(commands) {
		t.Fatal("GetCommandsCompletion returned the wrong number of commands")
	}
}

func TestBuildCompactMap(t *testing.T) {
	routes := []map[string]string{
		{"get": "/items"},
		{"post": "/items"},
		{"get": "/missing"},
		{"delete": "/items"},
	}
	compact := BuildCompactMap(routes)
	if compact["get"]["/items"] == "" {
		t.Fatalf("BuildCompactMap() = %#v", compact)
	}
	if compact["get"]["/missing"] != "" {
		t.Fatalf("unknown route return type = %q", compact["get"]["/missing"])
	}
}

func TestFetchReturnTypeGuards(t *testing.T) {
	got := fetchReturnTypeFromSpec("get", "/missing")
	if got != "" {
		t.Fatalf("missing path return type = %q", got)
	}
	got = fetchReturnTypeFromSpec("patch", "/items")
	if got != "" && got != "void" {
		t.Fatalf("missing operation return type = %q", got)
	}
}
