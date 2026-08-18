package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(2),
	Use:   "best <name> <effect> [effect...]",
	Short: "Find the best equipment for the specified effects",
	Long: `Find the best equipment for the specified effects

Arguments:
  name     Name of your character.
  effect   The code of the effect.`,
	ValidArgsFunction: completion.CharacterName(1).Custom(0, database.Effects().Equipments().Keys).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		effects, err := utils.NormalizeBestPriorities(args[1:])
		if err != nil {
			return err
		}
		character, err := api.Characters(args[0])
		if err != nil {
			return err
		}

		cmd.SilenceUsage = false
		items, err := utils.BestFinder(character, effects...)
		if err != nil {
			return fmt.Errorf("find best equipment: %w", err)
		}
		return console.Auto(items)
	},
}

func init() {
	rootCmd.AddCommand(bestCmd)
}
