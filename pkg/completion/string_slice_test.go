package completion

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestStringSlice(t *testing.T) {
	values := func() []string {
		return []string{"alchemy", "mining", "woodcutting"}
	}
	complete := StringSlice(values)

	tests := []struct {
		name       string
		toComplete string
		want       []string
	}{
		{
			name:       "first value returns all options",
			toComplete: "m",
			want:       []string{"alchemy", "mining", "woodcutting"},
		},
		{
			name:       "after comma excludes values already selected",
			toComplete: "alchemy,",
			want:       []string{"alchemy,mining", "alchemy,woodcutting"},
		},
		{
			name:       "value after comma returns remaining options",
			toComplete: "alchemy,m",
			want:       []string{"alchemy,mining", "alchemy,woodcutting"},
		},
		{
			name:       "preserves multiple previous values",
			toComplete: "alchemy,mining,w",
			want:       []string{"alchemy,mining,woodcutting"},
		},
		{
			name:       "partial value is left for the shell to filter",
			toComplete: "alchemy,z",
			want:       []string{"alchemy,mining", "alchemy,woodcutting"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, directive := complete(&cobra.Command{}, nil, test.toComplete)
			if directive != cobra.ShellCompDirectiveNoFileComp {
				t.Fatalf("completion directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("completion suggestions = %#v, want %#v", got, test.want)
			}
		})
	}
}
