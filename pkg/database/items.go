package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

var (
	//go:embed items.json
	items []byte

	itemsList  []schemas.ItemSchema
	itemsCache map[string]schemas.ItemSchema
)

func GetItems() []schemas.ItemSchema {
	initItemsCache()
	return itemsList
}

func GetItem(code string) (schemas.ItemSchema, bool) {
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

func GetItemTypes() []string {
	initItemsCache()
	uniqueTypes := make(map[string]bool)
	for _, item := range itemsList {
		if item.Type != "" {
			uniqueTypes[item.Type] = true
		}
		if item.Subtype != "" {
			uniqueTypes[item.Subtype] = true
		}
	}
	types := make([]string, 0, len(uniqueTypes))
	for itemType := range uniqueTypes {
		types = append(types, itemType)
	}
	return types
}

var initItemsCache = sync.OnceFunc(func() {
	itemsCache = make(map[string]schemas.ItemSchema)
	err := json.Unmarshal(items, &itemsList)
	if err != nil {
		panic("failed to unmarshal items: " + err.Error())
	}
	for _, item := range itemsList {
		itemsCache[item.Code] = item
	}
})
