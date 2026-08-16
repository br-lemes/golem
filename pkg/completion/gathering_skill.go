package completion

import "github.com/br-lemes/golem/pkg/database"

func GatheringSkill(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.GatheringSkill(count)
}

func (b *CompletionBuilder) GatheringSkill(count int) *CompletionBuilder {
	return b.Custom(count, func() []string {
		return database.Enum("GatheringSkill")
	})
}
