package routine

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

type deps struct {
	characters               func(name string) (schemas.CharacterSchema, error)
	myActionBankWithdrawItem func(name string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error)
	myActionEquip            func(name string, equips []schemas.EquipSchema) (schemas.EquipmentTransactionSchema, error)
	myActionMove             func(name string, x, y int) (schemas.CharacterMovementDataSchema, error)
	myActionTransition       func(name string) (schemas.CharacterTransitionDataSchema, error)
	myBankItems              func() ([]schemas.SimpleItemSchema, error)
}

var defaultDeps = deps{
	characters:               api.Characters,
	myActionBankWithdrawItem: api.MyActionBankWithdrawItem,
	myActionEquip:            api.MyActionEquip,
	myActionMove:             api.MyActionMove,
	myActionTransition:       api.MyActionTransition,
	myBankItems:              api.MyBankItems,
}

func Equip(name string, equipments []schemas.EquipSchema) (schemas.CharacterSchema, error) {
	return equip(defaultDeps, name, equipments)
}

func Move(character schemas.CharacterSchema, code string) (schemas.CharacterSchema, error) {
	return move(defaultDeps, character, code)
}
