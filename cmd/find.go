package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find <name> <code>",
	Short: "Find the closest coordinates containing a specific code",
	Long: `Find the closest coordinates containing a specific code

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

		tile := database.FindClosest(character, code)
		if tile == nil {
			return fmt.Errorf("no coordinates found for code %s", code)
		}
		output(tile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
