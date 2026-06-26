package config

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/console"
)

func TestGetCharacters(t *testing.T) {
	Characters = map[string][]string{
		"Hero":  {"mining", "woodcutting"},
		"Mage":  {"alchemy"},
		"Empty": {},
	}

	got := GetCharacters()

	if len(got) != 3 {
		t.Errorf("GetCharacters() returned %d items, want 3", len(got))
	}

	sort.Strings(got)

	expectedHero := "Hero\tmining, woodcutting"
	hasHero := false
	for _, v := range got {
		if strings.Contains(v, "Hero") {
			hasHero = true
			if v != expectedHero {
				t.Errorf("GetCharacters() format wrong. Got %q, want %q",
					v, expectedHero)
			}
		}
	}

	if !hasHero {
		t.Error("GetCharacters() missing data for Hero")
	}
}

func TestConfirmSkill(t *testing.T) {
	Characters = map[string][]string{
		"Hero": {"mining", "woodcutting"},
		"Mage": {},
	}

	tests := []struct {
		name       string
		character  string
		skill      string
		userInput  string
		wantResult bool
	}{
		{
			name:       "Character with no configured skills returns true",
			character:  "Mage",
			skill:      "alchemy",
			userInput:  "",
			wantResult: true,
		},
		{
			name:       "Character with matching skill returns true",
			character:  "Hero",
			skill:      "mining",
			userInput:  "",
			wantResult: true,
		},
		{
			name:       "Character without skill declines prompt",
			character:  "Hero",
			skill:      "alchemy",
			userInput:  "n\n",
			wantResult: false,
		},
		{
			name:       "Character without skill accepts prompt",
			character:  "Hero",
			skill:      "alchemy",
			userInput:  "y\n",
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.userInput != "" {
				buf := bytes.NewBufferString(tt.userInput)
				console.Stdin = buf
			}

			gotResult := ConfirmSkill(tt.character, tt.skill)

			if gotResult != tt.wantResult {
				t.Errorf("ConfirmSkill() got = %v, want = %v",
					gotResult, tt.wantResult)
			}
		})
	}
}
