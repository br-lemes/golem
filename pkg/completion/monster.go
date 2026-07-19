package completion

import "github.com/br-lemes/golem/pkg/database"

func Monster(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Monster(count)
}

func (b *CompletionBuilder) Monster(count int) *CompletionBuilder {
	return b.Custom(count, database.GetMonsterCodes)
}
