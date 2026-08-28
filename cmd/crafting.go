package cmd

import (
	"fmt"
	"math"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var craftingCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "crafting <name> <code>",
	Short: "Start the crafting bot loop",
	Long: `Start the crafting bot loop

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).Item(1).Build(),
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
	rootCmd.AddCommand(craftingCmd)
	craftingCmd.Flags().IntP("quantity", "q", 0, "Amount of items to craft (0 for infinite/max available)")
}

func StartCraftingBot(name string, code string, qty int) error {
	item, found := database.Items().Get(code)
	if !found {
		return fmt.Errorf("item not found: %s", code)
	}
	if !isCraftable(*item) {
		return fmt.Errorf("item is not craftable: %s", code)
	}

	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	routine.Cooldown(character)

	skillLevel := getCraftSkill(character, *item.Craft.Skill)
	if skillLevel < *item.Craft.Level {
		return fmt.Errorf("character %s level too low. Required: %d, Current: %d", name, *item.Craft.Level, skillLevel)
	}

	equipments, err := best.FindEquipmentSchemas(character, best.EquipmentOptions{
		UniqueAdeptRing: true,
		Priorities:      best.CraftingPriorities(character, item),
	})
	if err != nil {
		return err
	}
	character, err = routine.Equip(name, equipments)
	if err != nil {
		return err
	}

	bankInventory, err := fetchAllBankItems()
	if err != nil {
		return err
	}

	totalInventory := make(map[string]int)
	for _, invItem := range *character.Inventory {
		if invItem.Code != "" {
			totalInventory[invItem.Code] = totalInventory[invItem.Code] + invItem.Quantity
		}
	}
	for code, amount := range bankInventory {
		totalInventory[code] = totalInventory[code] + amount
	}

	maxActionsPossible := math.MaxInt32
	for _, req := range *item.Craft.Items {
		available := totalInventory[req.Code]
		possibleWithThisIngredient := available / req.Quantity
		if possibleWithThisIngredient < maxActionsPossible {
			maxActionsPossible = possibleWithThisIngredient
		}
	}

	recipeOutputQty := *item.Craft.Quantity
	maxYieldPossible := maxActionsPossible * recipeOutputQty

	if maxYieldPossible == 0 {
		return fmt.Errorf("insufficient materials in inventory and bank to craft even %d %s", recipeOutputQty, item.Code)
	}

	targetQty := qty
	if targetQty == 0 {
		targetQty = maxYieldPossible
	}

	if targetQty > maxYieldPossible {
		return fmt.Errorf("requested %d items, but combined inventory and bank can only produce %d", targetQty, maxYieldPossible)
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

		remainingToCraft := targetQty - craftedSoFar
		neededActions := int(math.Ceil(float64(remainingToCraft) / float64(recipeOutputQty)))

		batchActionsPossible := math.MaxInt32
		for _, req := range *item.Craft.Items {
			currentHas := currentInventory[req.Code]
			possibleActionsWithThis := currentHas / req.Quantity
			if possibleActionsWithThis < batchActionsPossible {
				batchActionsPossible = possibleActionsWithThis
			}
		}

		if batchActionsPossible > neededActions {
			batchActionsPossible = neededActions
		}

		if batchActionsPossible > 0 {
			character, err = routine.Move(character, string(*item.Craft.Skill))
			if err != nil {
				return err
			}

			_, err = api.MyActionCrafting(name, schemas.SimpleItemSchema{
				Code:     item.Code,
				Quantity: batchActionsPossible,
			})
			if err != nil {
				return err
			}

			craftedSoFar = craftedSoFar + (batchActionsPossible * recipeOutputQty)
			continue
		}

		character, err = routine.Move(character, "bank")
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

		remainingToCraft = targetQty - craftedSoFar
		neededActions = int(math.Ceil(float64(remainingToCraft) / float64(recipeOutputQty)))

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

		currentBatchActions := neededActions
		for _, req := range *item.Craft.Items {
			alreadyHas := currentInventory[req.Code]
			totalNeededForBatch := req.Quantity * currentBatchActions
			stillNeeds := totalNeededForBatch - alreadyHas
			if stillNeeds > 0 {
				bankAvailable := bankInventory[req.Code]
				if bankAvailable < stillNeeds {
					currentBatchActions = (bankAvailable + alreadyHas) / req.Quantity
				}
			}
		}

		for {
			totalUnitsToWithdraw := 0
			for _, req := range *item.Craft.Items {
				totalNeededForBatch := req.Quantity * currentBatchActions
				alreadyHas := currentInventory[req.Code]
				stillNeeds := totalNeededForBatch - alreadyHas
				if stillNeeds > 0 {
					totalUnitsToWithdraw = totalUnitsToWithdraw + stillNeeds
				}
			}

			if totalUnitsToWithdraw <= freeSpace {
				break
			}

			currentBatchActions = currentBatchActions - 1
			if currentBatchActions == 0 {
				return fmt.Errorf("not enough inventory weight capacity to withdraw ingredients for even 1 item action")
			}
		}

		if currentBatchActions == 0 {
			return fmt.Errorf("unexpected calculation error or missing materials in bank")
		}

		var withdrawList []schemas.SimpleItemSchema
		for _, req := range *item.Craft.Items {
			totalNeededForBatch := req.Quantity * currentBatchActions
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
	character, err = routine.Move(character, "bank")
	if err != nil {
		return err
	}

	var finalDepositList []schemas.SimpleItemSchema
	for _, invItem := range *character.Inventory {
		if invItem.Code == item.Code {
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
	level, _ := utils.GetCharacterCraftingSkillLevel(character, string(skill))
	return level
}

func isCraftable(item schemas.ItemSchema) bool {
	return item.Craft != nil && item.Craft.Items != nil && len(*item.Craft.Items) > 0 && item.Craft.Level != nil && *item.Craft.Quantity >= 1 && item.Craft.Skill != nil
}
