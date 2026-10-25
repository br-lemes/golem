package routine

import "github.com/br-lemes/golem/pkg/schemas"

type BankOptions struct {
	Movement        MoveOptions
	Utility1        string
	Utility2        string
	FoodReserve     int
	Utility1Reserve int
	Utility2Reserve int
	Food            string // "" = disabled; "auto" = best available; otherwise, code
	FoodOnly        bool
	NoFood          bool
}

func Bank(character schemas.CharacterSchema, opts BankOptions) (schemas.CharacterSchema, error) {
	//+gocover:ignore:block production wrapper over tested implementation
	return bank(defaultDeps, character, opts)
}

func bank(d deps, character schemas.CharacterSchema, opts BankOptions) (schemas.CharacterSchema, error) {
	var err error
	character, err = utilityInventoryRestock(d, character, opts.Utility1, opts.Utility2)
	if err != nil {
		return character, err
	}
	foodReserve := opts.FoodReserve
	if opts.Food == "" || opts.NoFood {
		foodReserve = 0
	} else if foodReserve == 0 {
		foodReserve = DefaultFoodReserve
	}
	plan, err := PlanInventory(character, opts.Movement.NoTeleport, foodReserve, opts.Utility1Reserve, opts.Utility2Reserve)
	if err != nil {
		return character, err
	}
	bankItems, err := d.myBankItems()
	if err != nil {
		return character, err
	}
	bankQty := map[string]int{}
	for _, item := range bankItems {
		bankQty[item.Code] += item.Quantity
	}
	needsUtility, err := utilityCheck(character, opts.Utility1, opts.Utility2, bankQty)
	if err != nil {
		return character, err
	}
	needsFood := !opts.NoFood && foodCheck(character, opts.Food, bankQty)
	needsSpace := totalItems(character)+5 >= character.InventoryMaxItems
	needsReserve := needsUtilityReserve(character, opts, plan)
	if !needsSpace && !needsUtility && !needsFood && !needsReserve && (opts.Movement.NoTeleport || !needsTravelPotions(character, bankQty)) {
		return character, nil
	}
	character, err = move(d, character, "bank", opts.Movement)
	if err != nil {
		return character, err
	}
	items := GetInventoryItems(character, InventoryItemsOptions{
		KeepTravelPotions: true,
	})
	if len(items) > 0 {
		depositData, err := d.myActionBankDepositItem(character.Name, items)
		if err != nil {
			return character, err
		}
		character = depositData.Character
		for _, item := range items {
			bankQty[item.Code] += item.Quantity
		}
	}
	foodTarget := plan.FoodQuantity
	reserve1 := plan.Utility1ReserveQuantity
	reserve2 := plan.Utility2ReserveQuantity
	character, err = utilityRestock(d, character, opts.Utility1, opts.Utility2, bankQty, reserve1, reserve2)
	if err != nil {
		return character, err
	}
	if !opts.NoFood {
		character, err = foodRestock(d, character, opts.Food, bankQty, opts.FoodOnly, foodTarget)
		if err != nil {
			return character, err
		}
	}
	return restockTravelPotion(d, character, opts.Movement)
}

func needsTravelPotions(character schemas.CharacterSchema, bankQty map[string]int) bool {
	for _, code := range reservedTravelPotionCodes {
		if inventoryItemQuantity(character, code) < 1 && bankQty[code] > 0 {
			return true
		}
	}
	return false
}

func needsUtilityReserve(character schemas.CharacterSchema, opts BankOptions, plan InventoryPlan) bool {
	if !characterAtBank(character) {
		return false
	}
	utility1Missing := opts.Utility1 != "" && inventoryItemQuantity(character, opts.Utility1) < plan.Utility1ReserveQuantity
	utility2Missing := opts.Utility2 != "" && inventoryItemQuantity(character, opts.Utility2) < plan.Utility2ReserveQuantity
	return utility1Missing || utility2Missing
}
