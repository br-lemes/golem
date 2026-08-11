package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/tidwall/gjson"
)

var (
	//go:embed monsters.json
	monsters []byte

	monstersList     []schemas.MonsterSchema
	monstersCache    map[string]schemas.MonsterSchema
	monsterCodesList []string
)

var Monsters = newStore(jsonLoader[schemas.MonsterSchema](monsters), func(item *schemas.MonsterSchema) string { return item.Code })

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
	initMonsterCodes()
	return monsterCodesList
}

func GetBossCodes() []string {
	initMonstersCache()
	var result []string
	for _, monster := range monstersCache {
		if monster.Type == "boss" {
			result = append(result, monster.Code)
		}
	}
	return result
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

var initMonsterCodes = sync.OnceFunc(func() {
	gjson.GetBytes(monsters, "#.code").ForEach(func(_, value gjson.Result) bool {
		monsterCodesList = append(monsterCodesList, value.String())
		return true
	})
})
