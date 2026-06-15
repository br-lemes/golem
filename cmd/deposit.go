package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/task"
	"github.com/spf13/cobra"
)

var depositCmd = &cobra.Command{
	Use:   "deposit <name>",
	Short: "Deposit all items and gold to the bank",
	Long: `Deposit all items and gold to the bank

Arguments:
  name   Name of your character.`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			characters := config.GetCharacters()
			return characters, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		task.Cooldown(character)
		character, err = task.Move(character, "bank")
		if err != nil {
			return err
		}
		items := task.GetInventoryItems(character, []string{})
		if len(items) > 0 {
			_, err = api.MyActionBankDepositItem(character.Name,
				task.GetInventoryItems(character, []string{}))
			if err != nil {
				return err
			}
		}
		if character.Gold > 0 {
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
}
