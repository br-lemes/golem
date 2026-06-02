package cmd

import (
	"errors"
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	. "github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var craftingCmd = &cobra.Command{
	Use:   "crafting [name] [code]",
	Short: "Start the crafting bot loop",
	Args:  cobra.ExactArgs(2),
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
	craftingCmd.Flags().IntP("quantity", "q", 0, "Amount of items to craft (0 for infinite/max available)")
	botCmd.AddCommand(craftingCmd)
}

func StartCraftingBot(name string, code string, qty int) error {
	item, found := database.GetItem(code)
	if !found {
		return fmt.Errorf("item not found: %s", code)
	}
	if !isCraftable(item) {
		return fmt.Errorf("item is not craftable: %s", code)
	}
	if !confirmSkill(name, string(*item.Craft.Skill)) {
		return fmt.Errorf("operation cancelled")
	}
	character, err := apiCharacters(name)
	if err != nil {
		return err
	}
	characterSkillLevel := getCraftSkill(character, *item.Craft.Skill)
	if characterSkillLevel < *item.Craft.Level {
		return fmt.Errorf("character level in %s is %d, but requires %d",
			*item.Craft.Skill, characterSkillLevel, *item.Craft.Level)
	}

	maxPossible, err := EvaluateCraftingFeasibility(character, item)
	if err != nil {
		return err
	}

	if qty == 0 {
		qty = maxPossible
	}

	if qty > maxPossible {
		return fmt.Errorf("requested to craft %d, but you only have materials for %d", qty, maxPossible)
	}

	fmt.Printf("Starting crafting loop for %d units of %s\n", qty, code)
	handleMap(character, string(*item.Craft.Skill))
	_, err = apiCraft(name, code, qty)
	if err != nil {
		return err
	}
	return nil
}

func getCraftSkill(character CharacterSchema, skill CraftSkill) int {
	switch skill {
	case CraftSkill("alchemy"):
		return character.AlchemyLevel
	case CraftSkill("cooking"):
		return character.CookingLevel
	case CraftSkill("gearcrafting"):
		return character.GearcraftingLevel
	case CraftSkill("jewelrycrafting"):
		return character.JewelrycraftingLevel
	case CraftSkill("weaponcrafting"):
		return character.WeaponcraftingLevel
	case CraftSkill("woodcutting"):
		return character.WoodcuttingLevel
	case CraftSkill("mining"):
		return character.MiningLevel
	case CraftSkill("fishing"):
		return character.FishingLevel
	default:
		return 0
	}
}

func getInventory(character CharacterSchema) map[string]int {
	inventoryMap := make(map[string]int)
	if character.Inventory == nil {
		return inventoryMap
	}
	for _, slot := range *character.Inventory {
		if slot.Quantity > 0 {
			inventoryMap[slot.Code] = inventoryMap[slot.Code] + slot.Quantity
		}
	}
	return inventoryMap
}

func isCraftable(item ItemSchema) bool {
	return item.Craft != nil &&
		item.Craft.Items != nil &&
		len(*item.Craft.Items) > 0 &&
		item.Craft.Level != nil &&
		*item.Craft.Quantity == 1 &&
		item.Craft.Skill != nil
}

func EvaluateCraftingFeasibility(character CharacterSchema, item ItemSchema) (int, error) {
	inventoryMap := getInventory(character)

	maxPossible := -1
	for _, reqItem := range *item.Craft.Items {
		available := inventoryMap[reqItem.Code]
		if available == 0 {
			maxPossible = 0
			break
		}

		possibleWithThisIngredient := available / reqItem.Quantity
		if maxPossible == -1 {
			maxPossible = possibleWithThisIngredient
		}

		if possibleWithThisIngredient < maxPossible {
			maxPossible = possibleWithThisIngredient
		}
	}

	if maxPossible == 0 {
		return 0, errors.New("insufficient materials in inventory to craft this item")
	}

	return maxPossible, nil
}
