package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var moveCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "move <name> <code>",
	Short: "Move to the closest coordinates with a specific code",
	Long: `Move to the closest coordinates with a specific code

Arguments:
  name   Name of your character.
  code   The target identifier (e.g., chicken, ash_tree, copper_rocks, bank).`,
	ValidArgsFunction: completion.CharacterName(1).Map(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		codes := append(database.MapCodes(), database.EventContentCodes()...)
		if !slices.Contains(codes, code) {
			return fmt.Errorf("code '%s' not found", code)
		}
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		routine.Cooldown(character)
		_, err = routine.Move(character, code)
		return err
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
