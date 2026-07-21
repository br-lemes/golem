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
	myActionRest             func(name string) (schemas.CharacterRestDataSchema, error)
	myActionTransition       func(name string) (schemas.CharacterTransitionDataSchema, error)
	myActionUnequip          func(name string, unequips []schemas.UnequipSchema) (schemas.EquipmentTransactionSchema, error)
	myActionUse              func(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error)
	myBankItems              func() ([]schemas.SimpleItemSchema, error)
}

var defaultDeps = deps{
	characters:               api.Characters,
	myActionBankWithdrawItem: api.MyActionBankWithdrawItem,
	myActionEquip:            api.MyActionEquip,
	myActionMove:             api.MyActionMove,
	myActionRest:             api.MyActionRest,
	myActionTransition:       api.MyActionTransition,
	myActionUnequip:          api.MyActionUnequip,
	myActionUse:              api.MyActionUse,
	myBankItems:              api.MyBankItems,
}
