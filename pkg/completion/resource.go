package completion

import "github.com/br-lemes/golem/pkg/catalog"

func Resource(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Resource(count)
}

func (b *CompletionBuilder) Resource(count int) *CompletionBuilder {
	return b.Custom(count, catalog.Resources.Keys)
}
