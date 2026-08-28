package completion

import "github.com/br-lemes/golem/pkg/database"

func Tradeables(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Tradeables(count)
}

func (b *CompletionBuilder) Tradeables(count int) *CompletionBuilder {
	return b.Custom(count, database.Items().Tradeables().Keys)
}
