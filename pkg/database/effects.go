package database

import (
	_ "embed"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed effects.json
var effects []byte

type effectCatalog struct {
	*store[schemas.EffectSchema, string]
	equipments *view[schemas.EffectSchema, string]
}

func Effects() *effectCatalog {
	return effectsCatalog
}

func (c *effectCatalog) Equipments() *view[schemas.EffectSchema, string] {
	return c.equipments
}

var effectsCatalog = func() *effectCatalog {
	store := newStore(jsonLoader[schemas.EffectSchema](effects), func(effect *schemas.EffectSchema) string {
		return effect.Code
	})
	return &effectCatalog{
		store: store,
		equipments: store.View(func(effect *schemas.EffectSchema) bool {
			return effect.Type == "equipment"
		}),
	}
}()
