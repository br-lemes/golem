package completion

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
)

func equipmentArgs() []string {
	var result []string
	for _, item := range database.Items().Equipments().All() {
		slots := database.EquipmentTypeToSlots[item.Type]
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
