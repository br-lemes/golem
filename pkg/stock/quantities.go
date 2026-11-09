package stock

import (
	"reflect"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

func ItemQuantities(bankItems []schemas.SimpleItemSchema, characters []schemas.CharacterSchema) map[string]int {
	return itemQuantities(bankItems, characters)
}

func GoldQuantity(bankGold int, characters []schemas.CharacterSchema) int {
	return goldQuantity(bankGold, characters)
}

func heldQuantities() (map[string]int, error) {
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return nil, err
	}
	bankItems, err := api.MyBankItems()
	if err != nil {
		return nil, err
	}
	return itemQuantities(bankItems, characters), nil
}

func itemQuantities(bankItems []schemas.SimpleItemSchema, characters []schemas.CharacterSchema) map[string]int {
	quantities := map[string]int{}
	for _, item := range bankItems {
		quantities[item.Code] += item.Quantity
	}
	for _, character := range characters {
		if character.Inventory != nil {
			for _, slot := range *character.Inventory {
				quantities[slot.Code] += slot.Quantity
			}
		}
		for code, quantity := range countInSlots(character) {
			quantities[code] += quantity
		}
	}
	return quantities
}

func goldQuantity(bankGold int, characters []schemas.CharacterSchema) int {
	quantity := bankGold
	for _, character := range characters {
		quantity += character.Gold
	}
	return quantity
}

func countInSlots(character schemas.CharacterSchema) map[string]int {
	quantities := map[string]int{}
	value := reflect.ValueOf(character)
	typeOfValue := value.Type()
	for index := 0; index < value.NumField(); index++ {
		field := typeOfValue.Field(index)
		if !strings.HasSuffix(field.Name, "Slot") || field.Type.Kind() != reflect.String {
			continue
		}
		code := value.Field(index).String()
		if code == "" {
			continue
		}
		quantity := value.FieldByName(field.Name + "Quantity")
		if quantity.IsValid() && quantity.Kind() == reflect.Int {
			quantities[code] += int(quantity.Int())
		} else {
			quantities[code]++
		}
	}
	return quantities
}
