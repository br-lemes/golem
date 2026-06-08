package cmd

import "github.com/spf13/cobra"

var moveCmd = &cobra.Command{
	Use:   "move <name> <code>",
	Short: "Move to the closest coordinates with a specific code",
	Long: `Move to the closest coordinates with a specific code

Arguments:
  name   Name of your character.
  code   The target identifier (e.g., chicken, ash_tree, copper_rocks, bank).`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		character, err := apiCharacters(name)
		if err != nil {
			return err
		}
		_, err = handleMap(character, code)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
