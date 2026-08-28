package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
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
		flags, err := utils.ReadFlags[bestFlags](cmd)
		if err != nil {
			return err
		}
		err = bestCraftingValidate(code)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return bestCraftingRun(name, code, flags)
	},
}

func bestCraftingValidate(code string) error {
	item, ok := database.Items().Get(code)
	if !ok || item.Craft == nil || item.Craft.Skill == nil {
		return fmt.Errorf("item not found or is not craftable: %s", code)
	}
	return nil
}

func bestCraftingRun(name, code string, flags bestFlags) error {
	item, _ := database.Items().Get(code)
	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	items, err := best.FindEquipment(character, best.EquipmentOptions{
		UniqueAdeptRing: !flags.AllowDuplicateAdeptRing,
		Priorities:      best.CraftingPriorities(character, item),
	})
	if err != nil {
		return err
	}
	return console.Auto(items)
}

func init() {
	bestCmd.AddCommand(bestCraftingCmd)
	err := utils.RegisterFlags[bestFlags](bestCraftingCmd)
	if err != nil {
		panic(err)
	}
}
