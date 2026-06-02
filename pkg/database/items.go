package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/br-lemes/golem/pkg/schemas"
)

//go:embed items.json
var items []byte

var itemsCache map[string]ItemSchema

func GetItem(code string) (ItemSchema, bool) {
	initItemsCache()
	item, exists := itemsCache[code]
	return item, exists
}

var initItemsCache = sync.OnceFunc(func() {
	itemsCache = make(map[string]ItemSchema)
	var result []ItemSchema
	err := json.Unmarshal(items, &result)
	if err != nil {
		return
	}
	for _, item := range result {
		itemsCache[item.Code] = item
	}
})
