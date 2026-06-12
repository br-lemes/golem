package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

var (
	//go:embed effects.json
	effects []byte

	effectsList  []schemas.EffectSchema
	effectsCache map[string]schemas.EffectSchema
)

func GetEffects() []schemas.EffectSchema {
	initEffectsCache()
	return effectsList
}

func GetEffect(code string) (schemas.EffectSchema, bool) {
	initEffectsCache()
	effect, exists := effectsCache[code]
	return effect, exists
}

func GetEffectCodes() []string {
	initEffectsCache()
	var codes []string
	for code := range effectsCache {
		codes = append(codes, code)
	}
	return codes
}

var initEffectsCache = sync.OnceFunc(func() {
	effectsCache = make(map[string]schemas.EffectSchema)
	err := json.Unmarshal(effects, &effectsList)
	if err != nil {
		panic("failed to unmarshal effects: " + err.Error())
	}
	for _, effect := range effectsList {
		effectsCache[effect.Code] = effect
	}
})
