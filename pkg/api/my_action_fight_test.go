package api

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
)

func TestMyActionFight(t *testing.T) {
	tests := []struct {
		name   string
		result string
		multi  bool
		want   string
	}{
		{
			name:   "win",
			result: "win",
			want:   "XP gained: 10, Gold gained: 5, Drops: 2x iron\n",
		},
		{
			name:   "loss",
			result: "loss",
			want:   "💀 Fight lost!\n",
		},
		{
			name:   "win with multiple characters",
			result: "win",
			multi:  true,
			want:   "[one] XP gained: 10, Gold gained: 5, Drops: 2x iron\n[two] XP gained: 10, Gold gained: 5\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache.CleanCharacters()
			t.Cleanup(cache.CleanCharacters)
			oldClient := defaultClient
			oldStdout := console.Stdout
			t.Cleanup(func() {
				defaultClient = oldClient
				console.Stdout = oldStdout
			})
			output := &bytes.Buffer{}
			console.Stdout = output
			defaultClient = testClient(func(r *http.Request) ([]byte, error) {
				if r.Method != http.MethodPost || r.URL.Path != "/my/hero/action/fight" {
					t.Errorf("request = %s %s, want POST /my/hero/action/fight", r.Method, r.URL.Path)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				wantBody := `{"participants":["ally"]}`
				if string(body) != wantBody {
					t.Errorf("body = %s, want %s", body, wantBody)
				}
				fightCharacters := `{"character_name":"one","xp":10,"gold":5,"drops":[{"code":"iron","quantity":2}]}`
				if test.multi {
					fightCharacters += `,{"character_name":"two","xp":10,"gold":5,"drops":[]}`
				}
				response := `{"data":{"characters":[],"fight":{"result":"` + test.result + `","characters":[` + fightCharacters + `]}}}`
				return []byte(response), nil
			})

			participants := []string{"ally"}
			_, err := MyActionFight("hero", participants)
			if err != nil {
				t.Fatal(err)
			}
			if output.String() != test.want {
				t.Errorf("output = %q, want %q", output.String(), test.want)
			}
		})
	}
}

func TestMyActionFightReturnsErrors(t *testing.T) {
	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"fight unavailable"}}`))
	_, err := MyActionFight("hero", nil)
	if err == nil {
		t.Fatal("expected request error")
	}

	defaultClient = responseClient(http.StatusOK, []byte("invalid json"))
	_, err = MyActionFight("hero", nil)
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
