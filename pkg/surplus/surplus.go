package surplus

import (
	"reflect"
	"strings"

	"github.com/br-lemes/golem/pkg/schemas"
)

type Input struct {
	BankItems  []schemas.SimpleItemSchema
	Characters []schemas.CharacterSchema
}

type Result struct {
	Item        schemas.ItemSchema
	Total       int
	Surplus     int
	DominatedBy []string
}

func toolSkill(item schemas.ItemSchema) string {
	if item.Effects == nil {
		//+gocover:ignore:block cannot occur in real catalog
		return ""
	}
	for _, effect := range *item.Effects {
		switch effect.Code {
		case "woodcutting", "mining", "fishing", "alchemy":
			return effect.Code
		}
	}
	//+gocover:ignore:block cannot occur in real catalog
	return ""
}

func equippedItems(character schemas.CharacterSchema) map[string]int {
	result := map[string]int{}
	value := reflect.ValueOf(character)
	typeOf := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := typeOf.Field(i)
		if !strings.HasSuffix(field.Name, "Slot") || field.Type.Kind() != reflect.String {
			continue
		}
		code := value.Field(i).String()
		if code == "" {
			continue
		}
		result[code]++
	}
	return result
}
