package completion

import "github.com/br-lemes/golem/pkg/catalog"

func Weapon(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Weapon(count)
}

func (b *CompletionBuilder) Weapon(count int) *CompletionBuilder {
	return b.Custom(count, catalog.Items().Weapons().Keys)
}
