package cmd

import (
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
		_, err = routine.Deposit(character, keepTypes)
		return err
	},
}

func init() {
	rootCmd.AddCommand(depositCmd)
	depositCmd.Flags().StringSliceVarP(&keepTypes, "keep", "k", []string{}, "Types of items to keep in inventory")
	err := depositCmd.RegisterFlagCompletionFunc("keep", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		types := database.Items().Types()
		types = append(types, "gold")
		return types, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}
