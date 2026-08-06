package completion

import (
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func NPCSell(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.NPCSell(count)
}

func (b *CompletionBuilder) NPCSell(count int) *CompletionBuilder {
	return b.Custom(count, GetNPCSellItems)
}

func GetNPCSellItems() []string {
	items := database.NpcsItems.Filter(func(item *schemas.NPCItemSchema) bool {
		if item.SellPrice != nil && *item.SellPrice > 0 {
			return true
		}
		return false
	})
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.Code
	}
	return result
}
