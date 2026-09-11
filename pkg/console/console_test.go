package console

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestAuto(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		input        any
		wantInOutput string
		wantErr      bool
	}{
		{
			name:         "Outputs YAML when format is yaml",
			format:       "yaml",
			input:        map[string]string{"name": "Go"},
			wantInOutput: "name: Go",
			wantErr:      false,
		},
		{
			name:         "Decodes JSON bytes and outputs JSON string",
			format:       "json",
			input:        []byte(`{"alive":true}`),
			wantInOutput: `{"alive":true}`,
			wantErr:      false,
		},
		{
			name:         "Returns error for invalid JSON bytes",
			format:       "json",
			input:        []byte(`{"invalid"-json}`),
			wantInOutput: "",
			wantErr:      true,
		},
		{
			name:         "Outputs JSON automatically when stdout is not a terminal",
			format:       "auto",
			input:        map[string]bool{"active": true},
			wantInOutput: `{"active":true}`,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Format = tt.format

			buf := &bytes.Buffer{}
			Stdout = buf

			err := Auto(tt.input)

			hasError := err != nil
			if hasError != tt.wantErr {
				t.Errorf("Auto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotOutput := strings.TrimSpace(buf.String())
				if !strings.Contains(gotOutput, tt.wantInOutput) {
					t.Errorf("Auto() got = %q, want to contain %q", gotOutput, tt.wantInOutput)
				}
			}
		})
	}
}

func TestAutoAppliesPathFilters(t *testing.T) {
	tests := []struct {
		name    string
		exclude []string
		only    []string
		ifPaths []string
		want    string
	}{
		{
			name:    "exclude",
			exclude: []string{"secret"},
			want:    `{"name":"Go","active":true}`,
		},
		{
			name: "only",
			only: []string{"name"},
			want: `{"name":"Go"}`,
		},
		{
			name:    "exclude-if",
			ifPaths: []string{"stats.active == true"},
			want:    `{"name":"Go","secret":"hidden"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Format = "json"
			Exclude = test.exclude
			Only = test.only
			ExcludeIf = test.ifPaths
			Stdout = &bytes.Buffer{}
			data := map[string]any{
				"name":   "Go",
				"active": true,
				"secret": "hidden",
				"stats":  map[string]any{"active": true},
			}

			err := Auto(data)
			if err != nil {
				t.Fatal(err)
			}
			output := strings.TrimSpace(Stdout.(*bytes.Buffer).String())
			if test.name == "only" && output != test.want {
				t.Fatalf("Auto() = %q, want %q", output, test.want)
			}
			if test.name == "exclude" && strings.Contains(output, "secret") {
				t.Fatalf("Auto() = %q, want secret excluded", output)
			}
			if test.name == "exclude-if" && strings.Contains(output, "stats") {
				t.Fatalf("Auto() = %q, want stats excluded", output)
			}
		})
	}
	Exclude = nil
	Only = nil
	ExcludeIf = nil
}

func TestAutoReturnsPathFilterErrors(t *testing.T) {
	tests := []struct {
		name    string
		exclude []string
		only    []string
		ifPaths []string
	}{
		{name: "exclude", exclude: []string{"value"}},
		{name: "only", only: []string{"value"}},
		{name: "exclude-if", ifPaths: []string{"value.name == Go"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Format = "json"
			Exclude = test.exclude
			Only = test.only
			ExcludeIf = test.ifPaths
			Stdout = io.Discard

			err := Auto(map[string]any{"value": make(chan int)})
			if err == nil {
				t.Fatal("Auto() returned nil error for unserializable data")
			}
		})
	}
	Exclude = nil
	Only = nil
	ExcludeIf = nil
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name       string
		userInput  string
		wantResult bool
	}{
		{name: "Accepts lowercase y", userInput: "y\n", wantResult: true},
		{name: "Accepts full word yes", userInput: "yes\n", wantResult: true},
		{name: "Accepts uppercase input", userInput: "YES\n", wantResult: true},
		{name: "Rejects lowercase n", userInput: "n\n", wantResult: false},
		{name: "Rejects unknown text", userInput: "maybe\n", wantResult: false},
		{name: "Rejects empty input", userInput: "\n", wantResult: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.userInput)
			Stdin = buf
			Stdout = io.Discard

			gotResult := Confirm("Proceed?")

			if gotResult != tt.wantResult {
				t.Errorf("Confirm() for input %q got = %v, want = %v", tt.userInput, gotResult, tt.wantResult)
			}
		})
	}
}

func TestPathFilters(t *testing.T) {
	data := map[string]any{
		"name":   "Go",
		"active": true,
		"stats":  map[string]any{"level": 10, "score": 3},
		"items":  []any{map[string]any{"code": "ore", "qty": 2}},
	}

	filtered, err := excludePaths(data, []string{"stats.score"})
	if err != nil || filtered.(map[string]any)["stats"].(map[string]any)["score"] != nil {
		t.Fatalf("excludePaths() = %#v, %v", filtered, err)
	}
	only, err := onlyPaths(data, []string{"stats.level", "items.*.code"})
	if err != nil || only.(map[string]any)["stats"] == nil || only.(map[string]any)["items"] == nil {
		t.Fatalf("onlyPaths() = %#v, %v", only, err)
	}
	conditional, err := excludeIfPaths(data, []string{"stats.level >= 10"})
	if err != nil || conditional.(map[string]any)["stats"] != nil {
		t.Fatalf("excludeIfPaths() = %#v, %v", conditional, err)
	}
}

func TestExcludeIfComparisons(t *testing.T) {
	operators := []struct {
		operator string
		want     bool
	}{
		{">", false},
		{">=", true},
		{"<", false},
		{"<=", true},
		{"==", true},
		{"!=", false},
	}
	for _, test := range operators {
		got, err := compareValue(float64(10), test.operator, "10")
		if err != nil || got != test.want {
			t.Errorf("compareValue(%q) = %t, %v", test.operator, got, err)
		}
	}
	got, err := compareValue("ready", "==", "ready")
	if err != nil || !got {
		t.Fatalf("compareValue(string) = %t, %v", got, err)
	}
	got, err = compareValue("ready", "!=", "pending")
	if err != nil || !got {
		t.Fatalf("compareValue(string !=) = %t, %v", got, err)
	}
	_, err = compareValue(float64(1), ">", "invalid")
	if err == nil {
		t.Fatal("compareValue() returned nil error for invalid number")
	}
}

func TestExcludeIfReturnsExpressionErrors(t *testing.T) {
	tests := []struct {
		name       string
		expression string
	}{
		{name: "invalid expression", expression: "stats.level >="},
		{name: "invalid path", expression: "level == 10"},
		{name: "invalid pattern", expression: "stats.[ == 10"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := map[string]any{
				"level": 10,
				"stats": map[string]int{"level": 10},
			}
			_, err := excludeIfPaths(data, []string{test.expression})
			if err == nil {
				t.Fatalf("excludeIfPaths() returned nil error for %q", test.expression)
			}
		})
	}
}

func TestExcludeIfPathHandlesNonObjectValue(t *testing.T) {
	value := "ready"
	got, err := excludeIfPath(value, []string{"status"}, "==", "ready")
	if err != nil || got != value {
		t.Fatalf("excludeIfPath() = %#v, %v, want unchanged value", got, err)
	}
}

func TestExcludeIfPathReturnsInvalidPatternError(t *testing.T) {
	_, err := excludeIfPath(map[string]any{"status": "ready"}, []string{"["}, "==", "ready")
	if err == nil {
		t.Fatal("excludeIfPath() returned nil error for invalid pattern")
	}
}

func TestConditionMatchesArrayWildcard(t *testing.T) {
	value := []any{
		map[string]any{"active": false},
		map[string]any{"active": true},
	}
	matched, err := conditionMatches(value, []string{"*", "active"}, "==", "true")
	if err != nil || !matched {
		t.Fatalf("conditionMatches() = %t, %v, want true", matched, err)
	}
}

func TestConditionMatchesArrayWildcardReturnsError(t *testing.T) {
	value := []any{map[string]any{"active": true}}
	_, err := conditionMatches(value, []string{"*", "active"}, ">", "invalid")
	if err == nil {
		t.Fatal("conditionMatches() returned nil error for invalid comparison")
	}
}

func TestConditionMatchesArrayWithoutMatch(t *testing.T) {
	value := []any{map[string]any{"active": true}}
	matched, err := conditionMatches(value, []string{"status"}, "==", "ready")
	if err != nil || matched {
		t.Fatalf("conditionMatches() = %t, %v, want false", matched, err)
	}
}

func TestInputAndDebugOutput(t *testing.T) {
	Stdin = strings.NewReader("value\n")
	Stdout = io.Discard
	got := Input("Enter")
	if got != "value" {
		t.Fatalf("Input() = %q, want value", got)
	}
	Stderr = &bytes.Buffer{}
	Debug = true
	Debugf("debug %d\n", 1)
	Errorf("error %d\n", 2)
	if !strings.Contains(Stderr.(*bytes.Buffer).String(), "DEBUG: debug 1") {
		t.Fatal("Debugf() did not write debug output")
	}
	if !strings.Contains(Stderr.(*bytes.Buffer).String(), "error 2") {
		t.Fatal("Errorf() did not write error output")
	}
	Debug = false
}

func TestOutputReturnsWriterError(t *testing.T) {
	Stdout = failingWriter{}
	err := Json(map[string]string{"key": "value"}, false)
	if err == nil {
		t.Fatal("Json() returned nil error for failing writer")
	}
}

func TestOnlyPathHandlesArrayIndexes(t *testing.T) {
	value := []any{"first", "second"}
	got, kept := onlyPath(value, [][]string{{"1"}})
	if !kept || got.([]any)[0] != nil || got.([]any)[1] != "second" {
		t.Fatalf("onlyPath() = %#v, %t, want second element", got, kept)
	}
	got, kept = onlyPath(value, [][]string{{"9"}})
	if kept || got.([]any)[0] != nil || got.([]any)[1] != nil {
		t.Fatalf("onlyPath() = %#v, %t, want no elements", got, kept)
	}
}

func TestOnlyPathLeavesScalarUnchanged(t *testing.T) {
	value := "Go"
	got, kept := onlyPath(value, [][]string{{"name"}})
	if kept || got != value {
		t.Fatalf("onlyPath() = %#v, %t, want unchanged scalar", got, kept)
	}
}

func TestExcludePathHandlesArrayIndexesAndWildcards(t *testing.T) {
	value := []any{"first", "second"}
	got, removed := excludePath(value, []string{"0"})
	if removed || len(got.([]any)) != 1 || got.([]any)[0] != "second" {
		t.Fatalf("excludePath() = %#v, %t, want second element", got, removed)
	}
	got, removed = excludePath([]any{"first", "second"}, []string{"*"})
	if !removed || got != nil {
		t.Fatalf("excludePath() = %#v, %t, want removed array", got, removed)
	}
}

func TestExcludePathHandlesRemainingCases(t *testing.T) {
	got, removed := excludePath("value", nil)
	if !removed || got != nil {
		t.Fatalf("excludePath(empty) = %#v, %t, want removed", got, removed)
	}

	value := map[string]any{"stats": map[string]any{"score": 10}}
	got, removed = excludePath(value, []string{"stats", "score"})
	if removed || len(got.(map[string]any)["stats"].(map[string]any)) != 0 {
		t.Fatalf("excludePath(nested) = %#v, %t, want empty stats", got, removed)
	}

	value = map[string]any{
		"items": []any{
			map[string]any{"secret": "one"},
			map[string]any{"secret": "two"},
		},
	}
	got, removed = excludePath(value, []string{"items", "*", "secret"})
	if removed || got == nil {
		t.Fatalf("excludePath(array nested) = %#v, %t, want updated value", got, removed)
	}

	array := []any{"first", "second"}
	got, removed = excludePath(array, []string{"invalid"})
	if removed || got.([]any)[0] != "first" {
		t.Fatalf("excludePath(invalid index) = %#v, %t, want unchanged array", got, removed)
	}
	got, removed = excludePath(array, []string{"1", "value"})
	if removed || got == nil {
		t.Fatalf("excludePath(nested index) = %#v, %t, want unchanged array", got, removed)
	}
	got, removed = excludePath("value", []string{"field"})
	if removed || got != "value" {
		t.Fatalf("excludePath(scalar) = %#v, %t, want unchanged value", got, removed)
	}

	value = map[string]any{"items": []any{"first", "second"}}
	got, removed = excludePath(value, []string{"items", "*"})
	if removed || len(got.(map[string]any)) != 0 {
		t.Fatalf("excludePath(remove child) = %#v, %t, want empty map", got, removed)
	}
}

func TestJsonAndYamlReturnMarshalErrors(t *testing.T) {
	data := make(chan int)
	err := Json(data, false)
	if err == nil {
		t.Fatal("Json() returned nil error for channel")
	}
	err = Yaml(data, false)
	if err == nil {
		t.Fatal("Yaml() returned nil error for channel")
	}
}

func TestOutputHighlightsWhenColorEnabled(t *testing.T) {
	Stdout = &bytes.Buffer{}
	Color = true
	err := output([]byte("name: Go\n"), false)
	if err != nil {
		t.Fatal(err)
	}
	Color = false
}

func TestConfirmAndInputReturnErrors(t *testing.T) {
	Stdin = strings.NewReader("")
	if Confirm("Proceed?") {
		t.Fatal("Confirm() returned true for empty input")
	}
	if Input("Enter") != "" {
		t.Fatal("Input() returned value for empty input")
	}
}
