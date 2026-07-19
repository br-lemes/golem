package completion

import "github.com/br-lemes/golem/pkg/database"

func Item(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Item(count)
}

func (b *CompletionBuilder) Item(count int) *CompletionBuilder {
	return b.Custom(count, database.GetItemCodes)
}
