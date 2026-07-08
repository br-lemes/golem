package cmd

import (
	"fmt"
	"maps"
	"slices"
	"sort"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

type ItemResult struct {
	Item  string `json:"item"`
	Value string `json:"value"`
}

var bestCraftingCmd = &cobra.Command{
	Use:   "crafting <name>",
	Short: "Find the best equipment for crafting",
	Long: `Find the best equipment for crafting

Arguments:
  name   Name of your character.
`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		character, err := api.Characters(name)
		if err != nil {
			return fmt.Errorf("failed to get character: %w", err)
		}
		bankItems, err := api.MyBankItems()
		if err != nil {
			return fmt.Errorf("failed to get bank items: %w", err)
		}
		owned := make(map[string]int)
		for _, bItem := range bankItems {
			owned[bItem.Code] += bItem.Quantity
		}
		if character.Inventory != nil {
			for _, invItem := range *character.Inventory {
				owned[invItem.Code] += invItem.Quantity
			}
		}
		equippedSlots := []string{
			character.AmuletSlot,
			character.Artifact1Slot,
			character.Artifact2Slot,
			character.Artifact3Slot,
			character.BagSlot,
			character.BodyArmorSlot,
			character.BootsSlot,
			character.HelmetSlot,
			character.LegArmorSlot,
			character.Ring1Slot,
			character.Ring2Slot,
			character.RuneSlot,
			character.ShieldSlot,
			character.Utility1Slot,
			character.Utility2Slot,
			character.WeaponSlot,
		}
		for _, code := range equippedSlots {
			if code != "" {
				owned[code]++
			}
		}
		allItems := database.GetItems()
		var validItems []schemas.ItemSchema
		itemValues := make(map[string]int)
		for _, item := range allItems {
			if item.Level > character.Level {
				continue
			}
			targetEffect := "wisdom"
			if item.Type == "bag" {
				targetEffect = "inventory_space"
			}
			hasEffect := false
			value := 0
			if item.Effects != nil {
				for _, effect := range *item.Effects {
					if effect.Code == targetEffect {
						hasEffect = true
						value = effect.Value
					}
				}
			}
			if !hasEffect {
				continue
			}
			validItems = append(validItems, item)
			itemValues[item.Code] = value
		}
		sort.Slice(validItems, func(i, j int) bool {
			return itemValues[validItems[i].Code] >
				itemValues[validItems[j].Code]
		})
		slots := slices.Collect(maps.Keys(database.EquipmentSlotToTypes))
		slices.Sort(slots)

		result := make(map[string]ItemResult)
		for _, slot := range slots {
			itemType := database.EquipmentSlotToTypes[slot]
			for _, item := range validItems {
				if item.Type == itemType && owned[item.Code] > 0 {
					val := itemValues[item.Code]
					result[slot] = ItemResult{
						Item:  item.Code,
						Value: formatDisplayValue(item.Type, val),
					}
					owned[item.Code]--
					break
				}
			}
		}
		console.Auto(result)
		return nil
	},
}

func formatDisplayValue(itemType string, value int) string {
	if itemType == "bag" {
		return fmt.Sprintf("+%d inventory space", value)
	}
	return fmt.Sprintf("+%d wisdom", value)
}

func init() {
	bestCmd.AddCommand(bestCraftingCmd)
}
