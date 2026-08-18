package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestCraftingCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "crafting <name> <code>",
	Short: "Find the best equipment for crafting",
	Long: `Find the best equipment for crafting

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		item, ok := database.Items.Get(code)
		if !ok || item.Craft == nil || item.Craft.Skill == nil {
			return fmt.Errorf("item not found or is not craftable: %s", args[1])
		}
		character, err := api.Characters(name)
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true
		skill := string(*item.Craft.Skill)
		skillLevel := craftSkillLevel(character, skill)
		priorities := []string{"inventory_space"}
		if skillLevel-item.Level <= 10 && skillLevel > 0 {
			priorities = []string{"wisdom", "inventory_space"}
		}
		items, err := utils.BestFinder(character, priorities...)
		if err != nil {
			return err
		}
		return console.Auto(items)
	},
}

func craftSkillLevel(c schemas.CharacterSchema, skill string) int {
	switch skill {
	case "alchemy":
		return c.AlchemyLevel
	case "cooking":
		return c.CookingLevel
	case "gearcrafting":
		return c.GearcraftingLevel
	case "jewelrycrafting":
		return c.JewelrycraftingLevel
	case "mining":
		return c.MiningLevel
	case "weaponcrafting":
		return c.WeaponcraftingLevel
	case "woodcutting":
		return c.WoodcuttingLevel
	}
	return 0
}

func init() {
	bestCmd.AddCommand(bestCraftingCmd)
}
