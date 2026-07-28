package routine

import "github.com/br-lemes/golem/pkg/schemas"

type BankOptions struct {
	Utility1 string
	Utility2 string
	Food     string // "" = disabled; "auto" = best available; otherwise, code
}

func Bank(character schemas.CharacterSchema, opts BankOptions) (schemas.CharacterSchema, error) {
	return bank(defaultDeps, character, opts)
}

func bank(d deps, character schemas.CharacterSchema, opts BankOptions) (schemas.CharacterSchema, error) {
	bankItems, err := d.myBankItems()
	if err != nil {
		return character, err
	}
	bankQty := map[string]int{}
	for _, item := range bankItems {
		bankQty[item.Code] += item.Quantity
	}
	needsUtility, err :=
		utilityCheck(character, opts.Utility1, opts.Utility2, bankQty)
	if err != nil {
		return character, err
	}
	needsFood := foodCheck(character, opts.Food, bankQty)
	needsSpace := totalItems(character)+5 >= character.InventoryMaxItems
	if !needsSpace && !needsUtility && !needsFood {
		return character, nil
	}
	character, err = Move(character, "bank")
	if err != nil {
		return character, err
	}
	items := GetInventoryItems(character, nil)
	if len(items) > 0 {
		depositData, err := d.myActionBankDepositItem(character.Name, items)
		if err != nil {
			return character, err
		}
		character = depositData.Character
	}
	character, err =
		utilityRestock(d, character, opts.Utility1, opts.Utility2, bankQty)
	if err != nil {
		return character, err
	}
	character, err = foodRestock(d, character, opts.Food, bankQty)
	if err != nil {
		return character, err
	}
	return character, nil
}
