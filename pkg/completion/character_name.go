package completion

import (
	"github.com/br-lemes/golem/pkg/utils"
)

func CharacterName(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.CharacterName(count)
}

func (b *CompletionBuilder) CharacterName(count int) *CompletionBuilder {
	return b.Custom(count, utils.GetCharacters)
}
