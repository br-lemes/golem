package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/tidwall/gjson"
)

var (
	//go:embed items.json
	items []byte

	itemsList     []schemas.ItemSchema
	itemsCache    map[string]schemas.ItemSchema
	itemCodesList []string
	itemTypesList []string
)

var Items = newStore(jsonLoader[schemas.ItemSchema](items), func(item *schemas.ItemSchema) string { return item.Code })

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
	initItemsCodes()
	return itemCodesList
}

func GetItemTypes() []string {
	initItemTypes()
	return itemTypesList
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

var initItemsCodes = sync.OnceFunc(func() {
	gjson.GetBytes(items, "#.code").ForEach(func(_, value gjson.Result) bool {
		itemCodesList = append(itemCodesList, value.String())
		return true
	})
})

var initItemTypes = sync.OnceFunc(func() {
	uniqueTypes := make(map[string]struct{})

	gjson.GetBytes(items, "#.type").ForEach(func(_, value gjson.Result) bool {
		s := value.String()
		if s != "" {
			uniqueTypes[s] = struct{}{}
		}
		return true
	})
	gjson.GetBytes(items, "#.subtype").ForEach(func(_, value gjson.Result) bool {
		s := value.String()
		if s != "" {
			uniqueTypes[s] = struct{}{}
		}
		return true
	})

	for t := range uniqueTypes {
		itemTypesList = append(itemTypesList, t)
	}
})
