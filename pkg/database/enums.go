package database

import (
	_ "embed"
	"encoding/json"
	"sort"
	"sync"
)

//go:embed enums.json
var enums []byte

type enumCatalog map[string][]string

func Enums() enumCatalog {
	return enumsCatalog()
}

func (c enumCatalog) All() map[string][]string {
	return c
}

func (c enumCatalog) Get(name string) ([]string, bool) {
	values, exists := c[name]
	return values, exists
}

func (c enumCatalog) Keys() []string {
	keys := make([]string, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

var enumsCatalog = sync.OnceValue(func() enumCatalog {
	var catalog enumCatalog
	err := json.Unmarshal(enums, &catalog)
	if err != nil {
		panic(err)
	}
	return catalog
})
