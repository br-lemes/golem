package completion

import (
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func SingleValue(fetcher func() []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fetcher(), cobra.ShellCompDirectiveNoFileComp
	}
}

func UniqueList(fetcher func() []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return stringSlice(fetcher, true)
}

func RepeatableList(fetcher func() []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return stringSlice(fetcher, false)
}

func stringSlice(fetcher func() []string, unique bool) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		idx := strings.LastIndex(toComplete, ",")
		if idx == -1 {
			return fetcher(), cobra.ShellCompDirectiveNoFileComp
		}

		prefix := toComplete[:idx]
		used := strings.Split(prefix, ",")
		var suggestions []string
		for _, value := range fetcher() {
			if unique && slices.Contains(used, value) {
				continue
			}
			suggestions = append(suggestions, prefix+","+value)
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	}
}
