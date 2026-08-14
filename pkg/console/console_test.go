package console

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

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
