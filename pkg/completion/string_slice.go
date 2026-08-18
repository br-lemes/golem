package completion

import (
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func StringSlice(fetcher func() []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		idx := strings.LastIndex(toComplete, ",")
		if idx == -1 {
			return fetcher(), cobra.ShellCompDirectiveNoFileComp
		}

		prefix := toComplete[:idx]
		used := strings.Split(prefix, ",")
		var suggestions []string
		for _, value := range fetcher() {
			if slices.Contains(used, value) {
				continue
			}
			suggestions = append(suggestions, prefix+","+value)
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}
