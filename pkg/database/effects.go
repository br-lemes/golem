package database

import (
	_ "embed"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed effects.json
var effects []byte

var Effects = newStore(jsonLoader[schemas.EffectSchema](effects), func(effect *schemas.EffectSchema) string {
	return effect.Code
})
