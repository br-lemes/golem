package completion

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/catalog"
)

func equipmentArgs() []string {
	var result []string
	for _, item := range catalog.Items().Equipments().All() {
		slots := catalog.EquipmentTypeToSlots[item.Type]
		for slot := 1; slot <= len(slots); slot++ {
			if len(slots) == 1 {
				result = append(result, item.Code)
				break
			}
			result = append(result, fmt.Sprintf("%s@%d", item.Code, slot))
		}
	}
	return result
}

func EquipmentArgs(count int) *CompletionBuilder {
	builder := &CompletionBuilder{}
	return builder.EquipmentArgs(count)
}

func (b *CompletionBuilder) EquipmentArgs(count int) *CompletionBuilder {
	return b.Custom(count, equipmentArgs)
}
