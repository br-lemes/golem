package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/br-lemes/golem/pkg/schemas"
)

var (
	//go:embed items.json
	items []byte

	itemsList  []ItemSchema
	itemsCache map[string]ItemSchema
)

func GetItems() []ItemSchema {
	initItemsCache()
	return itemsList
}

func GetItem(code string) (ItemSchema, bool) {
	initItemsCache()
	item, exists := itemsCache[code]
	return item, exists
}

func GetItemCodes() []string {
	initItemsCache()
	codes := make([]string, 0, len(itemsCache))
	for code := range itemsCache {
		codes = append(codes, code)
	}
	return codes
}

var initItemsCache = sync.OnceFunc(func() {
	itemsCache = make(map[string]ItemSchema)
	err := json.Unmarshal(items, &itemsList)
	if err != nil {
		panic("failed to unmarshal items: " + err.Error())
	}
	for _, item := range itemsList {
		itemsCache[item.Code] = item
	}
})
