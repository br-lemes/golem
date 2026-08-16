package completion

import "github.com/br-lemes/golem/pkg/database"

func Map(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.Map(count)
}

func (b *CompletionBuilder) Map(count int) *CompletionBuilder {
	return b.Custom(count, func() []string {
		return append(database.MapCodes(), database.EventContentCodes()...)
	})
}
