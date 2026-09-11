package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
)

func TestAPIURL(t *testing.T) {
	tests := []struct {
		environment string
		want        string
	}{
		{environment: "sandbox", want: "https://api.sandbox.artifactsmmo.com"},
		{environment: "beta", want: "https://api.beta.artifactsmmo.com"},
		{environment: "", want: "https://api.artifactsmmo.com"},
		{environment: "standard", want: "https://api.artifactsmmo.com"},
	}
	for _, test := range tests {
		t.Run(test.environment, func(t *testing.T) {
			got := apiURL(test.environment)
			if got != test.want {
				t.Fatalf("apiURL(%q) = %q, want %q", test.environment, got, test.want)
			}
		})
	}
}

func TestInitializeWithToken(t *testing.T) {
	oldToken := token
	oldBaseURL := baseURL
	t.Cleanup(func() {
		token = oldToken
		baseURL = oldBaseURL
	})

	apiConfig := config.API{Environment: "sandbox", Token: "secret"}
	got, prompted := Initialize(apiConfig)
	if got != "secret" {
		t.Fatalf("Initialize() token = %q, want secret", got)
	}
	if prompted {
		t.Fatal("Initialize() prompted with a configured token")
	}
	wantURL := "https://api.sandbox.artifactsmmo.com"
	if baseURL != wantURL {
		t.Fatalf("baseURL = %q, want sandbox URL", baseURL)
	}
}

func TestInitializeReadsMissingToken(t *testing.T) {
	oldToken := token
	oldBaseURL := baseURL
	oldStdin := console.Stdin
	oldStdout := console.Stdout
	t.Cleanup(func() {
		token = oldToken
		baseURL = oldBaseURL
		console.Stdin = oldStdin
		console.Stdout = oldStdout
	})
	console.Stdin = strings.NewReader("entered-token\n")
	console.Stdout = &bytes.Buffer{}

	got, prompted := Initialize(config.API{Environment: "beta"})
	if got != "entered-token" {
		t.Fatalf("Initialize() token = %q, want entered-token", got)
	}
	if !prompted {
		t.Fatal("Initialize() did not report prompting")
	}
	if baseURL != "https://api.beta.artifactsmmo.com" {
		t.Fatalf("baseURL = %q, want beta URL", baseURL)
	}
}
