package best

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

type deps struct {
	characters  func(string) (schemas.CharacterSchema, error)
	myBankItems func() ([]schemas.SimpleItemSchema, error)
}

var defaultDeps = deps{characters: api.Characters, myBankItems: api.MyBankItems}
