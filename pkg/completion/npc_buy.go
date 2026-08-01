package completion

import (
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func NPCBuy(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.NPCBuy(count)
}

func (b *CompletionBuilder) NPCBuy(count int) *CompletionBuilder {
	return b.Custom(count, GetNPCBuyItems)
}

func GetNPCBuyItems() []string {
	items := database.NpcsItems.Filter(func(item *schemas.NPCItemSchema) bool {
		if item.BuyPrice != nil && *item.BuyPrice > 0 {
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
