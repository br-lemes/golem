package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var foodCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "food <name>",
	Short: "Show the automatically selected food",
	Long: `Show the automatically selected food

Arguments:
  name   Name of your character.`,
	ValidArgsFunction: completion.CharacterName(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		character, err := api.Characters(args[0])
		if err != nil {
			return err
		}
		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		bankQty := map[string]int{}
		for _, item := range items {
			bankQty[item.Code] += item.Quantity
		}
		code := routine.SelectFood(character, bankQty)
		if code == "" {
			return fmt.Errorf("no suitable food available in bank")
		}
		return outputItemCounts([]string{code})
	},
}

func init() {
	rootCmd.AddCommand(foodCmd)
}
