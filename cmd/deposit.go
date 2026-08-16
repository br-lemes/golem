package cmd

import (
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var keepTypes []string

var depositCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "deposit <name>",
	Short: "Deposit all items and gold to the bank",
	Long: `Deposit all items and gold to the bank

Arguments:
  name   Name of your character.`,
	ValidArgsFunction: completion.CharacterName(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		routine.Cooldown(character)
		character, err = routine.Move(character, "bank")
		if err != nil {
			return err
		}
		items := routine.GetInventoryItems(character, keepTypes)
		if len(items) > 0 {
			_, err = api.MyActionBankDepositItem(character.Name, routine.GetInventoryItems(character, keepTypes))
			if err != nil {
				return err
			}
		}
		if character.Gold > 0 && !slices.Contains(keepTypes, "gold") {
			_, err = api.MyActionBankDepositGold(character.Name, character.Gold)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(depositCmd)
	depositCmd.Flags().StringSliceVarP(&keepTypes, "keep", "k", []string{}, "Types of items to keep in inventory")
	err := depositCmd.RegisterFlagCompletionFunc("keep", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		types := database.ItemTypes()
		types = append(types, "gold")
		return types, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}
