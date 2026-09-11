package completion

import "github.com/br-lemes/golem/pkg/catalog"

func Item(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Item(count)
}

func (b *CompletionBuilder) Item(count int) *CompletionBuilder {
	return b.Custom(count, catalog.Items().Keys)
}
