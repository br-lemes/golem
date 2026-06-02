package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/br-lemes/golem/pkg/schemas"
)

//go:embed monsters.json
var monsters []byte

var monstersCache map[string]MonsterSchema

func GetMonster(code string) (MonsterSchema, bool) {
	initMonstersCache()
	monster, exists := monstersCache[code]
	return monster, exists
}

var initMonstersCache = sync.OnceFunc(func() {
	monstersCache = make(map[string]MonsterSchema)
	var result []MonsterSchema
	err := json.Unmarshal(monsters, &result)
	if err != nil {
		return
	}
	for _, monster := range result {
		monstersCache[monster.Code] = monster
	}
})
