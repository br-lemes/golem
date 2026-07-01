package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/tidwall/gjson"
)

var (
	//go:embed effects.json
	effects []byte

	effectsList     []schemas.EffectSchema
	effectsCache    map[string]schemas.EffectSchema
	effectCodesList []string
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
	initEffectsCodes()
	return effectCodesList
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

var initEffectsCodes = sync.OnceFunc(func() {
	gjson.GetBytes(effects, "#.code").ForEach(func(_, value gjson.Result) bool {
		effectCodesList = append(effectCodesList, value.String())
		return true
	})
})
