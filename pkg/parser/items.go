package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

type Item struct {
	Code     string
	Slot     *string
	Quantity *int
}

type ItemList []Item

func Items(args []string) (ItemList, error) {
	result := make(ItemList, 0, len(args))
	for _, arg := range args {
		item, err := item(arg)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func item(arg string) (Item, error) {
	if arg == "" {
		return Item{}, fmt.Errorf("empty item argument")
	}

	codeEnd := strings.IndexAny(arg, "@=")
	if codeEnd < 0 {
		codeEnd = len(arg)
	}

	result := Item{Code: arg[:codeEnd]}
	if result.Code == "" {
		return Item{}, fmt.Errorf("missing item code: %s", arg)
	}
	_, exists := database.Items().Get(result.Code)
	if !exists {
		return Item{}, fmt.Errorf("item not found in catalog: %s", result.Code)
	}

	rest := arg[codeEnd:]
	if strings.HasPrefix(rest, "@") {
		rest = rest[1:]
		end := strings.IndexByte(rest, '=')
		if end < 0 {
			end = len(rest)
		}
		if end == 0 {
			return Item{}, fmt.Errorf("missing slot: %s", arg)
		}
		slot := rest[:end]
		_, err := strconv.Atoi(slot)
		if err != nil {
			return Item{}, fmt.Errorf("invalid slot number %s for item %s", slot, result.Code)
		}
		result.Slot = &slot
		rest = rest[end:]
	}

	if strings.HasPrefix(rest, "=") {
		quantityText := rest[1:]
		if quantityText == "" {
			return Item{}, fmt.Errorf("missing quantity: %s", arg)
		}
		quantity, err := strconv.Atoi(quantityText)
		if err != nil {
			return Item{}, fmt.Errorf("invalid quantity in %s: %w", arg, err)
		}
		if quantity < 1 {
			return Item{}, fmt.Errorf("quantity must be at least 1: %s", arg)
		}
		result.Quantity = &quantity
		rest = ""
	}

	return result, nil
}

func (item Item) EquipSchema() (schemas.EquipSchema, error) {
	catalogItem, exists := database.Items().Get(item.Code)
	if !exists {
		return schemas.EquipSchema{}, fmt.Errorf("item not found in catalog: %s", item.Code)
	}
	slots, isEquipment := database.EquipmentTypeToSlots[catalogItem.Type]
	if !isEquipment {
		return schemas.EquipSchema{}, fmt.Errorf("item is not equipment: %s", item.Code)
	}
	quantity := item.Quantity
	if quantity == nil {
		defaultQuantity := 1
		quantity = &defaultQuantity
	}
	equipment := schemas.EquipSchema{Code: item.Code, Quantity: quantity}
	if len(slots) == 1 {
		if item.Slot != nil {
			return schemas.EquipSchema{}, fmt.Errorf("cannot specify slot for item with single slot: %s", item.Code)
		}
		equipment.Slot = schemas.ItemSlot(slots[0])
	} else if item.Slot != nil {
		slot, _ := strconv.Atoi(*item.Slot)
		if slot < 1 || slot > len(slots) {
			return schemas.EquipSchema{}, fmt.Errorf("invalid slot number %s for item %s", *item.Slot, item.Code)
		}
		equipment.Slot = schemas.ItemSlot(slots[slot-1])
	}
	return equipment, nil
}

func (item Item) SimpleItemSchema() (schemas.SimpleItemSchema, error) {
	if item.Slot != nil {
		return schemas.SimpleItemSchema{}, fmt.Errorf("item %s specifies a slot", item.Code)
	}
	quantity := 1
	if item.Quantity != nil {
		quantity = *item.Quantity
	}
	return schemas.SimpleItemSchema{Code: item.Code, Quantity: quantity}, nil
}

func (items ItemList) EquipSchemas() ([]schemas.EquipSchema, error) {
	result := make([]schemas.EquipSchema, 0, len(items))
	for _, item := range items {
		equipment, err := item.EquipSchema()
		if err != nil {
			return nil, err
		}
		result = append(result, equipment)
	}
	return result, nil
}

func (items ItemList) SimpleItemSchemas() ([]schemas.SimpleItemSchema, error) {
	result := make([]schemas.SimpleItemSchema, 0, len(items))
	for _, item := range items {
		simpleItem, err := item.SimpleItemSchema()
		if err != nil {
			return nil, err
		}
		result = append(result, simpleItem)
	}
	return result, nil
}
