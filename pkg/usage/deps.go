package usage

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/schemas"
)

type deps struct {
	accountsCharacters func(string) ([]schemas.CharacterSchema, error)
	findEquipment      func(schemas.CharacterSchema, best.EquipmentOptions) (map[string]best.BestResult, error)
	findFight          func(schemas.CharacterSchema, schemas.MonsterSchema, map[string]int, bool, bool) (best.Result, error)
	myBankItems        func() ([]schemas.SimpleItemSchema, error)
}

var defaultDeps = deps{
	accountsCharacters: api.AccountsCharacters,
	findEquipment:      best.FindEquipment,
	findFight:          best.FindFightWithAvailable,
	myBankItems:        api.MyBankItems,
}
