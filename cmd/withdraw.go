package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/parser"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var withdrawCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(2),
	Use:   "withdraw <name> <code>[=<qtd>]...",
	Short: "Withdraw items from the bank",
	Long: `Withdraw items from the bank

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).Item(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		items, err := parser.Items(args[1:])
		if err != nil {
			return err
		}
		withdraw, err := items.SimpleItemSchemas()
		if err != nil {
			return err
		}
		movement, err := movementOptions(cmd)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		_, err = routine.Withdraw(character, withdraw, routine.WithdrawOptions{
			Movement: movement,
		})
		return err
	},
}

func init() {
	rootCmd.AddCommand(withdrawCmd)
}
