package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "find <name> <code>",
	Short: "Find the closest coordinates containing a specific code",
	Long: `Find the closest coordinates containing a specific code

Arguments:
  name   Name of your character.
  code   The target identifier (e.g., chicken, ash_tree, copper_rocks, bank).`,
	ValidArgsFunction: completion.CharacterName(1).Map(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		codes := append(database.GetMapCodes(), database.GetEventContentCodes()...)
		if !slices.Contains(codes, code) {
			return fmt.Errorf("code '%s' not found", code)
		}
		character, err := api.Characters(name)
		if err != nil {
			return err
		}

		tile := database.FindClosest(character, code)
		if tile == nil {
			return fmt.Errorf("no coordinates found for code %s", code)
		}
		return console.Auto(tile)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
