package completion

import "github.com/br-lemes/golem/pkg/catalog"

func Tradeable(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Tradeable(count)
}

func (b *CompletionBuilder) Tradeable(count int) *CompletionBuilder {
	return b.Custom(count, catalog.Items().Tradeables().Keys)
}
