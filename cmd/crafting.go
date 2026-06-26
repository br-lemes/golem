package cmd

import (
	"fmt"
	"math"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/task"
	"github.com/spf13/cobra"
)

var craftingCmd = &cobra.Command{
	Use:   "crafting <name> <code>",
	Short: "Start the crafting bot loop",
	Long: `Start the crafting bot loop

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			characters := config.GetCharacters()
			return characters, cobra.ShellCompDirectiveNoFileComp
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

		qty, _ := cmd.Flags().GetInt("quantity")

		err := StartCraftingBot(name, code, qty)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	craftingCmd.Flags().IntP("quantity", "q", 0,
		"Amount of items to craft (0 for infinite/max available)")
	rootCmd.AddCommand(craftingCmd)
}

func StartCraftingBot(name string, code string, qty int) error {
	item, found := database.GetItem(code)
	if !found {
		return fmt.Errorf("item not found: %s", code)
	}
	if !isCraftable(item) {
		return fmt.Errorf("item is not craftable: %s", code)
	}
	if !config.ConfirmSkill(name, string(*item.Craft.Skill)) {
		return fmt.Errorf("operation cancelled by user")
	}

	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	task.Cooldown(character)

	skillLevel := getCraftSkill(character, *item.Craft.Skill)
	if skillLevel < *item.Craft.Level {
		return fmt.Errorf(
			"character %s level too low. Required: %d, Current: %d",
			name, *item.Craft.Level, skillLevel)
	}

	bankInventory, err := fetchAllBankItems()
	if err != nil {
		return err
	}

	totalInventory := make(map[string]int)
	for _, invItem := range *character.Inventory {
		if invItem.Code != "" {
			totalInventory[invItem.Code] =
				totalInventory[invItem.Code] + invItem.Quantity
		}
	}
	for code, amount := range bankInventory {
		totalInventory[code] = totalInventory[code] + amount
	}

	maxPossible := math.MaxInt32
	for _, req := range *item.Craft.Items {
		available := totalInventory[req.Code]
		possibleWithThisIngredient := available / req.Quantity
		if possibleWithThisIngredient < maxPossible {
			maxPossible = possibleWithThisIngredient
		}
	}

	if maxPossible == 0 {
		return fmt.Errorf("insufficient materials in inventory and bank to craft even 1 %s", item.Code)
	}

	targetQty := qty
	if targetQty == 0 {
		targetQty = maxPossible
	}

	if targetQty > maxPossible {
		return fmt.Errorf("requested %d items, but combined inventory and bank can only produce %d", targetQty, maxPossible)
	}

	craftedSoFar := 0
	for craftedSoFar < targetQty {
		character, err = api.Characters(name)
		if err != nil {
			return err
		}

		currentInventory := make(map[string]int)
		for _, invItem := range *character.Inventory {
			if invItem.Code != "" {
				currentInventory[invItem.Code] = invItem.Quantity
			}
		}

		neededInInventory := targetQty - craftedSoFar
		batchPossible := math.MaxInt32
		for _, req := range *item.Craft.Items {
			currentHas := currentInventory[req.Code]
			possibleBatchWithThis := currentHas / req.Quantity
			if possibleBatchWithThis < batchPossible {
				batchPossible = possibleBatchWithThis
			}
		}

		if batchPossible > neededInInventory {
			batchPossible = neededInInventory
		}

		if batchPossible > 0 {
			character, err = task.Move(character, string(*item.Craft.Skill))
			if err != nil {
				return err
			}

			_, err = api.MyActionCrafting(name, schemas.SimpleItemSchema{
				Code:     item.Code,
				Quantity: batchPossible,
			})
			if err != nil {
				return err
			}

			craftedSoFar = craftedSoFar + batchPossible
			continue
		}

		character, err = task.Move(character, "bank")
		if err != nil {
			return err
		}

		var depositList []schemas.SimpleItemSchema
		for _, invItem := range *character.Inventory {
			if invItem.Code == "" {
				continue
			}
			isIngredient := false
			for _, req := range *item.Craft.Items {
				if req.Code == invItem.Code {
					isIngredient = true
					break
				}
			}
			if !isIngredient {
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

		currentInventory = make(map[string]int)
		totalItemsInInventory := 0
		slotsUsed := 0
		for _, invItem := range *character.Inventory {
			if invItem.Code != "" {
				currentInventory[invItem.Code] = invItem.Quantity
				totalItemsInInventory = totalItemsInInventory + invItem.Quantity
				slotsUsed = slotsUsed + 1
			}
		}

		bankInventory, err = fetchAllBankItems()
		if err != nil {
			return err
		}

		remainingToCraft := targetQty - craftedSoFar
		freeSpace := character.InventoryMaxItems - totalItemsInInventory
		freeSlots := 20 - slotsUsed

		neededIngredientsCount := 0
		for _, req := range *item.Craft.Items {
			alreadyHas := currentInventory[req.Code]
			if alreadyHas == 0 {
				neededIngredientsCount = neededIngredientsCount + 1
			}
		}

		if neededIngredientsCount > freeSlots {
			if freeSlots == 0 {
				return fmt.Errorf("inventory is completely full of other items")
			}
		}

		currentBatchSize := remainingToCraft
		for _, req := range *item.Craft.Items {
			alreadyHas := currentInventory[req.Code]
			totalNeededForBatch := req.Quantity * currentBatchSize
			stillNeeds := totalNeededForBatch - alreadyHas
			if stillNeeds > 0 {
				bankAvailable := bankInventory[req.Code]
				if bankAvailable < stillNeeds {
					currentBatchSize =
						(bankAvailable + alreadyHas) / req.Quantity
				}
			}
		}

		for {
			totalUnitsToWithdraw := 0
			for _, req := range *item.Craft.Items {
				totalNeededForBatch := req.Quantity * currentBatchSize
				alreadyHas := currentInventory[req.Code]
				stillNeeds := totalNeededForBatch - alreadyHas
				if stillNeeds > 0 {
					totalUnitsToWithdraw = totalUnitsToWithdraw + stillNeeds
				}
			}

			if totalUnitsToWithdraw <= freeSpace {
				break
			}

			currentBatchSize = currentBatchSize - 1
			if currentBatchSize == 0 {
				return fmt.Errorf("not enough inventory weight capacity to withdraw ingredients for even 1 item")
			}
		}

		if currentBatchSize == 0 {
			return fmt.Errorf("unexpected calculation error or missing materials in bank")
		}

		var withdrawList []schemas.SimpleItemSchema
		for _, req := range *item.Craft.Items {
			totalNeededForBatch := req.Quantity * currentBatchSize
			alreadyHas := currentInventory[req.Code]
			stillNeeds := totalNeededForBatch - alreadyHas
			if stillNeeds > 0 {
				withdrawList = append(withdrawList, schemas.SimpleItemSchema{
					Code:     req.Code,
					Quantity: stillNeeds,
				})
			}
		}

		if len(withdrawList) > 0 {
			_, err = api.MyActionBankWithdrawItem(name, withdrawList)
			if err != nil {
				return err
			}
		}
	}

	character, err = api.Characters(name)
	if err != nil {
		return err
	}
	character, err = task.Move(character, "bank")
	if err != nil {
		return err
	}

	var finalDepositList []schemas.SimpleItemSchema
	for _, invItem := range *character.Inventory {
		if invItem.Code == item.Code {
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

func fetchAllBankItems() (map[string]int, error) {
	result := make(map[string]int)
	items, err := api.MyBankItems()
	if err != nil {
		return result, err
	}
	for _, item := range items {
		result[item.Code] = result[item.Code] + item.Quantity
	}
	return result, nil
}

func getCraftSkill(character schemas.CharacterSchema, skill schemas.CraftSkill) int {
	switch skill {
	case schemas.CraftSkill("alchemy"):
		return character.AlchemyLevel
	case schemas.CraftSkill("cooking"):
		return character.CookingLevel
	case schemas.CraftSkill("gearcrafting"):
		return character.GearcraftingLevel
	case schemas.CraftSkill("jewelrycrafting"):
		return character.JewelrycraftingLevel
	case schemas.CraftSkill("weaponcrafting"):
		return character.WeaponcraftingLevel
	case schemas.CraftSkill("woodcutting"):
		return character.WoodcuttingLevel
	case schemas.CraftSkill("mining"):
		return character.MiningLevel
	case schemas.CraftSkill("fishing"):
		return character.FishingLevel
	default:
		return 0
	}
}

func isCraftable(item schemas.ItemSchema) bool {
	return item.Craft != nil && item.Craft.Items != nil &&
		len(*item.Craft.Items) > 0 && item.Craft.Level != nil &&
		*item.Craft.Quantity == 1 && item.Craft.Skill != nil
}
