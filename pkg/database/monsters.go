package database

import (
	_ "embed"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed monsters.json
var monsters []byte

var Monsters = newStore(jsonLoader[schemas.MonsterSchema](monsters), func(item *schemas.MonsterSchema) string { return item.Code })

var Bosses = Monsters.View(func(monster *schemas.MonsterSchema) bool {
	return monster.Type == "boss"
})
