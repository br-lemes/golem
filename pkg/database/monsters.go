package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

var (
	//go:embed monsters.json
	monsters []byte

	monstersList  []schemas.MonsterSchema
	monstersCache map[string]schemas.MonsterSchema
)

func GetMonsters() []schemas.MonsterSchema {
	initMonstersCache()
	return monstersList
}

func GetMonster(code string) (schemas.MonsterSchema, bool) {
	initMonstersCache()
	monster, exists := monstersCache[code]
	return monster, exists
}

func GetMonsterCodes() []string {
	initMonstersCache()
	codes := make([]string, 0, len(monstersCache))
	for code := range monstersCache {
		codes = append(codes, code)
	}
	return codes
}

var initMonstersCache = sync.OnceFunc(func() {
	monstersCache = make(map[string]schemas.MonsterSchema)
	err := json.Unmarshal(monsters, &monstersList)
	if err != nil {
		panic("failed to unmarshal monsters: " + err.Error())
	}
	for _, monster := range monstersList {
		monstersCache[monster.Code] = monster
	}
})
