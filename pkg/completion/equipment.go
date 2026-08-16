package completion

import "github.com/br-lemes/golem/pkg/database"

func Equipment(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Equipment(count)
}

func (b *CompletionBuilder) Equipment(count int) *CompletionBuilder {
	return b.Custom(count, database.Equipments)
}
