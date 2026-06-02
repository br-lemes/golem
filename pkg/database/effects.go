package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/br-lemes/golem/pkg/schemas"
)

//go:embed effects.json
var effects []byte

var effectsCache map[string]EffectSchema

func GetEffect(code string) (EffectSchema, bool) {
	initEffectsCache()
	effect, exists := effectsCache[code]
	return effect, exists
}

var initEffectsCache = sync.OnceFunc(func() {
	effectsCache = make(map[string]EffectSchema)
	var result []EffectSchema
	err := json.Unmarshal(effects, &result)
	if err != nil {
		return
	}
	for _, effect := range result {
		effectsCache[effect.Code] = effect
	}
})
