package routine

import (
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

var reservedTravelPotionCodes = []string{"forest_bank_potion", "recall_potion"}

type InventoryItemsOptions struct {
	KeepTypes         []string
	KeepTravelPotions bool
}

func Inventory(character schemas.CharacterSchema, keepTypes []string, options MoveOptions) (schemas.CharacterSchema, error) {
	//+gocover:ignore:block production wrapper over tested implementation
	return inventory(defaultDeps, character, keepTypes, options)
}

func InventorySpace(character schemas.CharacterSchema) int {
	return max(0, character.InventoryMaxItems-totalItems(character))
}

func SpaceAfterDeposit(character schemas.CharacterSchema) int {
	items := GetInventoryItems(character, InventoryItemsOptions{
		KeepTravelPotions: true,
	})
	depositable := 0
	for _, item := range items {
		depositable += item.Quantity
	}
	return max(0, character.InventoryMaxItems-totalItems(character)+depositable)
}

func KeepTravelReserve(character schemas.CharacterSchema, items []schemas.SimpleItemSchema) []schemas.SimpleItemSchema {
	reserve := make(map[string]int, len(reservedTravelPotionCodes))
	for _, code := range reservedTravelPotionCodes {
		reserve[code] = min(1, inventoryItemQuantity(character, code))
	}
	result := make([]schemas.SimpleItemSchema, 0, len(items))
	for _, item := range items {
		quantity := item.Quantity
		if reserve[item.Code] > 0 {
			kept := min(reserve[item.Code], quantity)
			reserve[item.Code] -= kept
			quantity -= kept
		}
		if quantity > 0 {
			item.Quantity = quantity
			result = append(result, item)
		}
	}
	return result
}

func GetInventoryItems(character schemas.CharacterSchema, options InventoryItemsOptions) []schemas.SimpleItemSchema {
	items := []schemas.SimpleItemSchema{}
	if character.Inventory == nil {
		return items
	}
	for _, item := range *character.Inventory {
		if item.Code == "" || item.Quantity == 0 {
			continue
		}
		if shouldKeepItem(item.Code, options.KeepTypes) {
			continue
		}
		items = append(items, schemas.SimpleItemSchema{
			Code:     item.Code,
			Quantity: item.Quantity,
		})
	}
	if options.KeepTravelPotions {
		return KeepTravelReserve(character, items)
	}
	return items
}

func RestockTravelPotion(character schemas.CharacterSchema, options MoveOptions) (schemas.CharacterSchema, error) {
	//+gocover:ignore:block production wrapper over tested implementation
	return restockTravelPotion(defaultDeps, character, options)
}

func inventory(d deps, character schemas.CharacterSchema, keepTypes []string, options MoveOptions) (schemas.CharacterSchema, error) {
	total := totalItems(character)
	if total+5 < character.InventoryMaxItems {
		return character, nil
	}
	character, err := move(d, character, "bank", options)
	if err != nil {
		return character, err
	}
	items := GetInventoryItems(character, InventoryItemsOptions{
		KeepTypes:         keepTypes,
		KeepTravelPotions: true,
	})
	if len(items) > 0 {
		transaction, err := d.myActionBankDepositItem(character.Name, items)
		if err != nil {
			return character, err
		}
		character = transaction.Character
	}
	return restockTravelPotion(d, character, options)
}

func totalItems(character schemas.CharacterSchema) int {
	if character.Inventory == nil {
		return 0
	}
	total := 0
	for _, item := range *character.Inventory {
		total += item.Quantity
	}
	return total
}

func inventoryItemQuantity(character schemas.CharacterSchema, code string) int {
	quantity := 0
	if character.Inventory == nil {
		return quantity
	}
	for _, item := range *character.Inventory {
		if item.Code == code {
			quantity += item.Quantity
		}
	}
	return quantity
}

func restockTravelPotion(d deps, character schemas.CharacterSchema, options MoveOptions) (schemas.CharacterSchema, error) {
	if options.NoTeleport || !characterAtBank(character) || d.myBankItems == nil || d.myActionBankWithdrawItem == nil {
		return character, nil
	}
	bankItems, err := d.myBankItems()
	if err != nil {
		return character, err
	}
	bankQuantity := make(map[string]int, len(bankItems))
	for _, item := range bankItems {
		bankQuantity[item.Code] += item.Quantity
	}
	freeSpace := character.InventoryMaxItems - totalItems(character)
	for _, code := range reservedTravelPotionCodes {
		missing := 1 - inventoryItemQuantity(character, code)
		quantity := min(missing, bankQuantity[code], freeSpace)
		if quantity <= 0 {
			continue
		}
		withdraw, err := d.myActionBankWithdrawItem(character.Name, []schemas.SimpleItemSchema{
			{Code: code, Quantity: quantity},
		})
		if err != nil {
			return character, err
		}
		character = withdraw.Character
		bankQuantity[code] -= quantity
		freeSpace -= quantity
	}
	return character, nil
}
func shouldKeepItem(code string, keepTypes []string) bool {
	if len(keepTypes) == 0 {
		return false
	}
	item, found := catalog.Items().Get(code)
	if !found {
		return false
	}
	for _, keepType := range keepTypes {
		if item.Type == keepType || item.Subtype == keepType {
			return true
		}
	}
	return false
}
