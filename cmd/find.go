package cmd

import (
	"github.com/br-lemes/golem/pkg/database"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <name> <code>",
	Short: "Find the closest coordinates containing a specific code",
	Long: `Find the closest coordinates containing a specific code

Arguments:
  name   Name of your character.
  code   The code of the monster.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		character, err := apiCharacters(name)
		if err != nil {
			return err
		}

		tile := database.FindClosest(character, code)
		output(tile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
