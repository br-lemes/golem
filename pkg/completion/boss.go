package completion

import "github.com/br-lemes/golem/pkg/catalog"

func Boss(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Boss(count)
}

func (b *CompletionBuilder) Boss(count int) *CompletionBuilder {
	return b.Custom(count, catalog.Bosses.Keys)
}
