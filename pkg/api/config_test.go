package api

import (
	"bytes"
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
)

func TestNewClientEnvironment(t *testing.T) {
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
			got := NewClient(config.API{Environment: test.environment}).BaseURL
			if got != test.want {
				t.Fatalf("NewClient(%q).BaseURL = %q, want %q", test.environment, got, test.want)
			}
		})
	}
}

func TestInitializeWithToken(t *testing.T) {
	oldClient := defaultClient
	t.Cleanup(func() {
		defaultClient = oldClient
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
	if defaultClient.BaseURL != wantURL {
		t.Fatalf("BaseURL = %q, want sandbox URL", defaultClient.BaseURL)
	}
}

func TestInitializeReadsMissingToken(t *testing.T) {
	oldClient := defaultClient
	oldStdin := console.Stdin
	oldStdout := console.Stdout
	t.Cleanup(func() {
		defaultClient = oldClient
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
	if defaultClient.BaseURL != "https://api.beta.artifactsmmo.com" {
		t.Fatalf("BaseURL = %q, want beta URL", defaultClient.BaseURL)
	}
}
