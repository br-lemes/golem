package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var equipCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "equip <name> <code>",
	Short: "Equip one or more items on a character",
	Long: `Equip one or more items on a character

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).Equipment(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		var equipments []schemas.EquipSchema
		for _, arg := range args[1:] {
			equipment, err := parseCliEquipment(arg)
			if err != nil {
				return err
			}
			equipments = append(equipments, equipment)
		}
		_, err := routine.Equip(args[0], equipments)
		return err
	},
}

func parseCliEquipment(arg string) (schemas.EquipSchema, error) {
	qty := 1
	slotNum := 0
	itemCode := arg
	hasSlot := strings.Contains(arg, "@")
	if hasSlot {
		parts := strings.Split(arg, "@")
		itemCode = parts[0]
		rightPart := parts[1]
		if strings.Contains(rightPart, "x") {
			subParts := strings.Split(rightPart, "x")
			slotNum, _ = strconv.Atoi(subParts[0])
			qty, _ = strconv.Atoi(subParts[1])
		} else {
			slotNum, _ = strconv.Atoi(rightPart)
		}
	}

	item, exists := database.GetItem(itemCode)
	if !exists {
		return schemas.EquipSchema{}, fmt.Errorf("item not found in database: %s", itemCode)
	}
	slots, ok := database.EquipmentTypeToSlots[item.Type]
	if !ok {
		return schemas.EquipSchema{}, fmt.Errorf("item type cannot be equipped: %s", item.Type)
	}
	if hasSlot && slotNum <= 0 {
		return schemas.EquipSchema{}, fmt.Errorf("invalid slot number %d for item %s", slotNum, itemCode)
	}
	if len(slots) > 1 {
		if !hasSlot {
			return schemas.EquipSchema{}, fmt.Errorf("must specify slot (e.g. @1 or @2) for item with multiple slots: %s", itemCode)
		}
		if slotNum > len(slots) {
			return schemas.EquipSchema{}, fmt.Errorf("invalid slot number %d for item %s", slotNum, itemCode)
		}
	}
	if len(slots) == 1 && hasSlot {
		return schemas.EquipSchema{}, fmt.Errorf("cannot specify slot number for item with single slot: %s", itemCode)
	}
	targetSlotStr := slots[0]
	if slotNum > 0 && slotNum <= len(slots) {
		targetSlotStr = slots[slotNum-1]
	}
	if qty > 1 && !strings.HasPrefix(targetSlotStr, "utility") {
		return schemas.EquipSchema{}, fmt.Errorf("cannot specify quantity for non-utility slot: %s", targetSlotStr)
	}

	return schemas.EquipSchema{
		Code:     itemCode,
		Quantity: &qty,
		Slot:     schemas.ItemSlot(targetSlotStr),
	}, nil
}

func init() {
	rootCmd.AddCommand(equipCmd)
}
