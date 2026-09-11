package utils

import (
	"strings"
	"testing"
)

func TestTemplateHelpers(t *testing.T) {
	args, format := extractPathArgs("/my/{name}/item/{id}")
	if strings.Join(args, ",") != "name,id" || format != "/my/%s/item/%s" {
		t.Fatalf("extractPathArgs() = %v, %q", args, format)
	}
	wants := map[string]string{
		"integer": "Int",
		"boolean": "Bool",
		"string":  "String",
		"":        "String",
	}
	for apiType, want := range wants {
		got := mapGoType(apiType)
		if got != want {
			t.Errorf("mapGoType(%q) = %q, want %q", apiType, got, want)
		}
	}
	collector := flagCollector{flagsCheck: make(map[string]TemplateParam)}
	err := collector.collectFlag(FieldInfo{Name: "size", Type: "integer"})
	if err != nil {
		t.Fatal(err)
	}
	err = collector.collectFlag(FieldInfo{Name: "size", Type: "integer"})
	if err != nil {
		t.Fatal(err)
	}
	err = collector.collectFlag(FieldInfo{Name: "size", Type: "boolean"})
	if err == nil {
		t.Fatal("collectFlag accepted a type conflict")
	}
}

func TestFormatArguments(t *testing.T) {
	data := TemplateData{
		Summary: "summary",
		Routes:  map[int]TemplateRoute{0: {}, 1: {PathArgs: []string{"id"}}},
	}
	data.formatArguments(map[string]string{"id": "the identifier"})
	if data.ArgUse != "[id]" || !strings.Contains(data.Long, "the identifier") {
		t.Fatalf("formatted arguments = %#v", data)
	}
	data = TemplateData{
		Routes: map[int]TemplateRoute{1: {PathArgs: []string{"name", "id"}}},
	}
	data.formatArguments(nil)
	if data.ArgUse != "<name id>" {
		t.Fatalf("required arguments = %q", data.ArgUse)
	}
}

func TestFormatArgumentsWithNoRoutes(t *testing.T) {
	data := TemplateData{}
	data.formatArguments(nil)
	if data.ArgUse != "" || data.Long != "" {
		t.Fatalf("formatArguments() = %#v, want unchanged data", data)
	}
}

func TestBuildAndRenderTemplate(t *testing.T) {
	data, err := BuildTemplateData("items")
	if err != nil {
		t.Fatal(err)
	}
	if data.CommandName != "items" || len(data.Routes) == 0 {
		t.Fatalf("BuildTemplateData() = %#v, want populated data", data)
	}
	_, err = RenderTemplate(data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = BuildTemplateData("missing-command")
	if err == nil {
		t.Fatal("unknown command accepted")
	}
}

func TestBuildTemplateDataHandlesRequestBody(t *testing.T) {
	data, err := BuildTemplateData("myActionEquip")
	if err != nil {
		t.Fatal(err)
	}
	if !data.HasBody {
		t.Fatalf("BuildTemplateData() HasBody = false, want true: %#v", data)
	}
	if len(data.Flags) == 0 {
		t.Fatal("BuildTemplateData() returned no request body fields")
	}
}

func TestRenderTemplateFormatsCode(t *testing.T) {
	data := TemplateData{
		CommandName: "testCommand",
		Summary:     "Test command",
		Method:      "GET",
		Routes: map[int]TemplateRoute{
			0: {Method: "get", Path: "/test", GoPathFormat: "/test"},
		},
	}

	code, err := RenderTemplate(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) == 0 {
		t.Fatal("RenderTemplate() returned empty code")
	}
	if !strings.Contains(string(code), "var testCommandCmd") {
		t.Fatalf("RenderTemplate() did not generate command: %s", code)
	}
}

func TestRenderTemplateReturnsFormattingError(t *testing.T) {
	_, err := RenderTemplate(TemplateData{CommandName: "invalid command"})
	if err == nil {
		t.Fatal("RenderTemplate() returned nil error for invalid generated Go")
	}
	if !strings.Contains(err.Error(), "failed to format source code") {
		t.Fatalf("RenderTemplate() error = %v, want formatting error", err)
	}
}
