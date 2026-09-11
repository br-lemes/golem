package completion

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCatalogCompletionsReturnValues(t *testing.T) {
	checks := []struct {
		name string
		fn   func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective)
	}{
		{"items", Item(1).Build()},
		{"maps", Map(1).Build()},
		{"monsters", Monster(1).Build()},
		{"bosses", Boss(1).Build()},
		{"resources", Resource(1).Build()},
		{"tradeables", Tradeable(1).Build()},
		{"skills", GatheringSkill(1).Build()},
	}
	for _, check := range checks {
		got, _ := check.fn(nil, nil, "")
		if len(got) == 0 {
			t.Errorf("%s completion returned no values", check.name)
		}
	}
}
