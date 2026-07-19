package completion

import "github.com/spf13/cobra"

func Custom(count int, fetcher func() []string) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Custom(count, fetcher)
}

func (b *CompletionBuilder) Custom(count int, fetcher func() []string) *CompletionBuilder {
	validator := func(cmd *cobra.Command, args []string, toComplete string, offset int) ([]string, int, bool) {
		currentArgIndex := len(args)
		if currentArgIndex < offset {
			return nil, count, false
		}
		if count > 0 {
			targetIndex := offset + count
			if currentArgIndex >= targetIndex {
				return nil, count, false
			}
		}
		results := fetcher()
		return results, count, true
	}
	b.validators = append(b.validators, validator)
	return b
}
