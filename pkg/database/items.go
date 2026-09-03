package database

import (
	_ "embed"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed items.json
var items []byte

type itemCatalog struct {
	*store[schemas.ItemSchema, string]
	equipments *view[schemas.ItemSchema, string]
	potions    *view[schemas.ItemSchema, string]
	tradeables *view[schemas.ItemSchema, string]
}

func Items() *itemCatalog {
	return itemsCatalog
}

func (c *itemCatalog) Equipments() *view[schemas.ItemSchema, string] {
	return c.equipments
}

func (c *itemCatalog) Potions() *view[schemas.ItemSchema, string] {
	return c.potions
}

func (c *itemCatalog) Tradeables() *view[schemas.ItemSchema, string] {
	return c.tradeables
}

func (c *itemCatalog) Types() []string {
	return itemTypes()
}

var itemsCatalog = func() *itemCatalog {
	store := newStore(jsonLoader[schemas.ItemSchema](items), func(item *schemas.ItemSchema) string {
		return item.Code
	})
	return &itemCatalog{
		store: store,
		equipments: store.View(func(item *schemas.ItemSchema) bool {
			_, ok := EquipmentTypeToSlots[item.Type]
			return ok
		}),
		potions: store.View(func(item *schemas.ItemSchema) bool {
			return item.Subtype == "potion"
		}),
		tradeables: store.View(func(item *schemas.ItemSchema) bool {
			return item.Tradeable
		}),
	}
}()

var itemTypes = sync.OnceValue(func() []string {
	seen := make(map[string]struct{})
	var types []string

	for _, item := range Items().All() {
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
