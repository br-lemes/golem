package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/task"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Use:   "move <name> <code>",
	Short: "Move to the closest coordinates with a specific code",
	Long: `Move to the closest coordinates with a specific code

Arguments:
  name   Name of your character.
  code   The target identifier (e.g., chicken, ash_tree, copper_rocks, bank).`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			characters := getCharacters()
			return characters, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			codes := database.GetMapCodes()
			return codes, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		task.Cooldown(character)
		_, err = task.Move(character, code)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
