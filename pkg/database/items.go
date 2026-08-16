package database

import (
	_ "embed"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed items.json
var items []byte

var Items = newStore(jsonLoader[schemas.ItemSchema](items), func(item *schemas.ItemSchema) string {
	return item.Code
})

var itemTypes = sync.OnceValue(func() []string {
	seen := make(map[string]struct{})
	var types []string

	for _, item := range Items.All() {
		for _, value := range []string{item.Type, item.Subtype} {
			if value == "" {
				continue
			}
			_, exists := seen[value]
			if exists {
				continue
			}
			seen[value] = struct{}{}
			types = append(types, value)
		}
	}
	return types
})

func ItemTypes() []string {
	return itemTypes()
}
