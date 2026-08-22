package routine

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSelectFoodMatchesAutomaticRestockOrder(t *testing.T) {
	character := schemas.CharacterSchema{Level: 10}

	got := SelectFood(character, map[string]int{
		"apple":          10,
		"cooked_chicken": 10,
	})

	if got != "cooked_chicken" {
		t.Fatalf("SelectFood() = %q, want %q", got, "cooked_chicken")
	}
}

func TestSelectFoodIgnoresUnavailableOrTooHighLevelFood(t *testing.T) {
	character := schemas.CharacterSchema{Level: 1}

	got := SelectFood(character, map[string]int{
		"cooked_chicken": 0,
		"cooked_beef":    10,
	})

	if got != "" {
		t.Fatalf("SelectFood() = %q, want no food", got)
	}
}
