package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var recyclingCmd = &cobra.Command{
	Use:   "recycling <name> <code>",
	Short: "Start the recycling bot loop",
	Long: `Start the recycling bot loop

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			codes := database.GetItemCodes()
			return codes, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		qty, err := cmd.Flags().GetInt("quantity")
		if err != nil {
			return err
		}

		err = StartRecyclingBot(name, code, qty)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	recyclingCmd.Flags().IntP("quantity", "q", 0,
		"Amount of items to recycle (0 for all available)")
	rootCmd.AddCommand(recyclingCmd)
}

func StartRecyclingBot(name string, code string, qty int) error {
	item, found := database.GetItem(code)
	if !found {
		return fmt.Errorf("item not found: %s", code)
	}

	if !isRecyclable(item) {
		return fmt.Errorf("item cannot be recycled")
	}

	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	routine.Cooldown(character)

	bankInventory, err := fetchAllBankItems()
	if err != nil {
		return err
	}

	totalItemsAvailable := bankInventory[item.Code]
	for _, invItem := range *character.Inventory {
		if invItem.Code == item.Code {
			totalItemsAvailable = totalItemsAvailable + invItem.Quantity
		}
	}

	if totalItemsAvailable == 0 {
		return fmt.Errorf("no items available to recycle: %s", item.Code)
	}

	targetQty := qty
	if targetQty == 0 {
		targetQty = totalItemsAvailable
	}

	if targetQty > totalItemsAvailable {
		return fmt.Errorf("requested to recycle %d items, but only %d are available", targetQty, totalItemsAvailable)
	}

	totalItemsInRecipe := 0
	for _, req := range *item.Craft.Items {
		totalItemsInRecipe = totalItemsInRecipe + req.Quantity
	}

	resourcesReturnedPerItem := ((totalItemsInRecipe - 1) / 5) + 1
	if resourcesReturnedPerItem < 1 {
		resourcesReturnedPerItem = 1
	}

	recycledSoFar := 0
	for recycledSoFar < targetQty {
		character, err = api.Characters(name)
		if err != nil {
			return err
		}

		currentInventoryCount := 0
		for _, invItem := range *character.Inventory {
			if invItem.Code != "" {
				currentInventoryCount = currentInventoryCount + invItem.Quantity
			}
		}

		remainingToRecycle := targetQty - recycledSoFar
		freeSpace := character.InventoryMaxItems - currentInventoryCount

		currentInventoryHas := 0
		for _, invItem := range *character.Inventory {
			if invItem.Code == item.Code {
				currentInventoryHas = invItem.Quantity
			}
		}

		if currentInventoryHas > 0 {
			character, err = routine.Move(character, string(*item.Craft.Skill))
			if err != nil {
				return err
			}

			batchToRecycle := currentInventoryHas
			if batchToRecycle > remainingToRecycle {
				batchToRecycle = remainingToRecycle
			}

			_, err = api.MyActionRecycling(name, schemas.SimpleItemSchema{
				Code:     item.Code,
				Quantity: batchToRecycle,
			})
			if err != nil {
				return err
			}

			recycledSoFar = recycledSoFar + batchToRecycle
			continue
		}

		character, err = routine.Move(character, "bank")
		if err != nil {
			return err
		}

		var depositList []schemas.SimpleItemSchema
		for _, invItem := range *character.Inventory {
			if invItem.Code != "" {
				depositList = append(depositList, schemas.SimpleItemSchema{
					Code:     invItem.Code,
					Quantity: invItem.Quantity,
				})
			}
		}

		if len(depositList) > 0 {
			_, err = api.MyActionBankDepositItem(name, depositList)
			if err != nil {
				return err
			}
			character, err = api.Characters(name)
			if err != nil {
				return err
			}
		}

		freeSpace = character.InventoryMaxItems

		bankInventory, err = fetchAllBankItems()
		if err != nil {
			return err
		}

		bankAvailable := bankInventory[item.Code]
		batchSize := remainingToRecycle
		if batchSize > bankAvailable {
			batchSize = bankAvailable
		}

		maxSafeBatch := freeSpace / resourcesReturnedPerItem
		if maxSafeBatch < 1 {
			maxSafeBatch = 1
		}

		if batchSize > maxSafeBatch {
			batchSize = maxSafeBatch
		}

		if batchSize > freeSpace {
			batchSize = freeSpace
		}

		if batchSize <= 0 {
			return fmt.Errorf("inventory is completely full or no items left in bank")
		}

		var withdrawList []schemas.SimpleItemSchema
		withdrawList = append(withdrawList, schemas.SimpleItemSchema{
			Code:     item.Code,
			Quantity: batchSize,
		})

		_, err = api.MyActionBankWithdrawItem(name, withdrawList)
		if err != nil {
			return err
		}
	}

	character, err = api.Characters(name)
	if err != nil {
		return err
	}
	character, err = routine.Move(character, "bank")
	if err != nil {
		return err
	}

	var finalDepositList []schemas.SimpleItemSchema
	for _, invItem := range *character.Inventory {
		if invItem.Code != "" {
			finalDepositList = append(finalDepositList,
				schemas.SimpleItemSchema{
					Code:     invItem.Code,
					Quantity: invItem.Quantity,
				})
		}
	}

	if len(finalDepositList) > 0 {
		_, err = api.MyActionBankDepositItem(name, finalDepositList)
		if err != nil {
			return err
		}
	}

	return nil
}

func isRecyclable(item schemas.ItemSchema) bool {
	if item.Craft == nil || item.Craft.Skill == nil {
		return false
	}
	switch string(*item.Craft.Skill) {
	case "alchemy", "gearcrafting", "jewelrycrafting", "weaponcrafting":
		return true
	default:
		return false
	}
}
