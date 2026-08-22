package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var foodCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "food [name]",
	Short: "Show the automatically selected food",
	Long: `Show the automatically selected food

Arguments:
  name   Name of your character.`,
	ValidArgsFunction: completion.CharacterName(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		characters := []schemas.CharacterSchema{}
		if len(args) == 1 {
			character, err := api.Characters(args[0])
			if err != nil {
				return err
			}
			characters = append(characters, character)
		} else {
			var err error
			characters, err = api.AccountsCharacters("")
			if err != nil {
				return err
			}
		}
		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		bankQty := map[string]int{}
		for _, item := range items {
			bankQty[item.Code] += item.Quantity
		}
		codes := []string{}
		for _, character := range characters {
			code := routine.SelectFood(character, bankQty)
			if code != "" {
				codes = append(codes, code)
			}
		}
		if len(codes) == 0 {
			return fmt.Errorf("no suitable food available in bank")
		}
		return outputItemCounts(uniqueStrings(codes))
	},
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(foodCmd)
}
