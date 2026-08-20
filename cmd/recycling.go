package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var recyclingCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "recycling <name> <code>",
	Short: "Start the recycling bot loop",
	Long: `Start the recycling bot loop

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		qty, err := cmd.Flags().GetInt("quantity")
		if err != nil {
			return err
		}
		enhanced, err := cmd.Flags().GetBool("enhanced")
		if err != nil {
			return err
		}

		err = StartRecyclingBot(name, code, qty, enhanced)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	recyclingCmd.Flags().IntP("quantity", "q", 0, "Amount of items to recycle (0 for all available)")
	recyclingCmd.Flags().Bool("enhanced", false, "Use enhanced recycling (costs gold for 30% more resources)")
	rootCmd.AddCommand(recyclingCmd)
}

func StartRecyclingBot(name string, code string, qty int, enhanced bool) error {
	item, found := database.Items.Get(code)
	if !found {
		return fmt.Errorf("item not found: %s", code)
	}

	if !isRecyclable(*item) {
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
	if enhanced {
		// Enhanced recycling returns 30% more resources, rounded to the
		// nearest whole resource and capped by the recipe's ingredients.
		resourcesReturnedPerItem = (resourcesReturnedPerItem*13 + 5) / 10
		resourcesReturnedPerItem = min(resourcesReturnedPerItem, totalItemsInRecipe)
	}

	goldPerIngredient := 0
	if enhanced {
		switch {
		case item.Level <= 20:
			goldPerIngredient = 5
		case item.Level <= 30:
			goldPerIngredient = 10
		case item.Level <= 40:
			goldPerIngredient = 15
		case item.Level <= 45:
			goldPerIngredient = 20
		default:
			goldPerIngredient = 25
		}
		requiredGold := totalItemsInRecipe * targetQty * goldPerIngredient
		if character.Gold < requiredGold {
			bank, bankErr := api.MyBank()
			if bankErr != nil {
				return bankErr
			}
			if character.Gold+bank.Gold < requiredGold {
				return fmt.Errorf("not enough gold for enhanced recycling: required %d, available %d", requiredGold, character.Gold+bank.Gold)
			}
			character, err = routine.Move(character, "bank")
			if err != nil {
				return err
			}
			_, err = api.MyActionBankWithdrawGold(name, requiredGold-character.Gold)
			if err != nil {
				return fmt.Errorf("failed to withdraw gold for enhanced recycling: %w", err)
			}
			character, err = api.Characters(name)
			if err != nil {
				return err
			}
			if character.Gold < requiredGold {
				return fmt.Errorf("not enough gold for enhanced recycling: required %d, available %d", requiredGold, character.Gold)
			}
		}
	}

	recycledSoFar := 0
	for recycledSoFar < targetQty {
		character, err = api.Characters(name)
		if err != nil {
			return err
		}

		currentInventoryCount := 0
		currentInventoryHas := 0
		for _, invItem := range *character.Inventory {
			if invItem.Code != "" {
				currentInventoryCount += invItem.Quantity
				if invItem.Code == item.Code {
					currentInventoryHas += invItem.Quantity
				}
			}
		}

		remainingToRecycle := targetQty - recycledSoFar
		freeSpace := character.InventoryMaxItems - currentInventoryCount

		if currentInventoryHas > 0 {
			batchToRecycle := currentInventoryHas
			if batchToRecycle > remainingToRecycle {
				batchToRecycle = remainingToRecycle
			}

			maxMaterialsReturned := batchToRecycle * resourcesReturnedPerItem
			netSpaceNeeded := maxMaterialsReturned - batchToRecycle

			if freeSpace >= netSpaceNeeded {
				character, err = routine.Move(character, string(*item.Craft.Skill))
				if err != nil {
					return err
				}

				enhancedPayload := enhanced
				quantityPayload := batchToRecycle
				_, err = api.MyActionRecycling(name, schemas.RecyclingSchema{
					Code:     item.Code,
					Quantity: &quantityPayload,
					Enhanced: &enhancedPayload,
				})
				if err != nil {
					return err
				}

				recycledSoFar += batchToRecycle
				continue
			}
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
			finalDepositList = append(finalDepositList, schemas.SimpleItemSchema{
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
