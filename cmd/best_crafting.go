package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestCraftingCmd = &cobra.Command{
	Use:   "crafting <name>",
	Short: "Find the best equipment for crafting",
	Long: `Find the best equipment for crafting

Arguments:
  name   Name of your character.
`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		items, err := utils.BestFinder(character, "",
			map[string]int{"wisdom": 4, "inventory_space": 2})
		if err != nil {
			return err
		}
		console.Auto(items)
		return nil
	},
}

func init() {
	bestCmd.AddCommand(bestCraftingCmd)
}
