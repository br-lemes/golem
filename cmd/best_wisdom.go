package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestWisdomCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "wisdom <name>",
	Short: "Find the best equipment for wisdom",
	Long: `Find the best equipment for wisdom

Arguments:
  name   Name of your character.`,
	ValidArgsFunction: completion.CharacterName(1).Build(),
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
		return console.Auto(items)
	},
}

func init() {
	bestCmd.AddCommand(bestWisdomCmd)
}
