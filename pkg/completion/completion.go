package completion

import "github.com/spf13/cobra"

type ArgValidator func(cmd *cobra.Command, args []string, toComplete string, offset int) ([]string, int, bool)

type CompletionBuilder struct {
	validators []ArgValidator
}

func (b *CompletionBuilder) Build() func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		currentOffset := 0
		for _, validator := range b.validators {
			suggestions, consumed, matched := validator(cmd, args, toComplete, currentOffset)
			if matched {
				return suggestions, cobra.ShellCompDirectiveNoFileComp
			}
			if consumed == 0 {
				currentOffset = len(args)
				continue
			}
			currentOffset = currentOffset + consumed
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}
