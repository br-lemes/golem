package completion

import "github.com/br-lemes/golem/pkg/database"

func Resource(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Resource(count)
}

func (b *CompletionBuilder) Resource(count int) *CompletionBuilder {
	return b.Custom(count, database.Resources.Keys)
}
