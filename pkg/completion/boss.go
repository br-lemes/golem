package completion

import "github.com/br-lemes/golem/pkg/database"

func Boss(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Boss(count)
}

func (b *CompletionBuilder) Boss(count int) *CompletionBuilder {
	return b.Custom(count, database.GetBossCodes)
}
