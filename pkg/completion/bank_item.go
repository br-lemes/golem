package completion

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/database"
)

func BankItem(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.BankItem(count)
}

func (b *CompletionBuilder) BankItem(count int) *CompletionBuilder {
	return b.Custom(count, GetBankItems)
}

func GetBankItems() []string {
	items := cache.GetBankItems()
	if items != nil {
		result := make([]string, len(items))
		for i, item := range items {
			result[i] = item.Code
		}
		return result
	}
	return database.GetItemCodes()
}
