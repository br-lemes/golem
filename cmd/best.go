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

type bestFlags struct {
	AllowDuplicateAdeptRing bool `flag:"allow-duplicate-adept-ring" desc:"Allow recommending ring_of_the_adept twice for this character"`
}

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
		name := args[0]
		flags, err := utils.ReadFlags[bestFlags](cmd)
		if err != nil {
			return err
		}
		effects, err := utils.NormalizeBestPriorities(args[1:])
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return bestRun(name, flags, effects)
	},
}

func bestRun(name string, flags bestFlags, effects []string) error {
	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	items, err := utils.BestFinder(character, !flags.AllowDuplicateAdeptRing, effects...)
	if err != nil {
		return fmt.Errorf("find best equipment: %w", err)
	}
	return console.Auto(items)
}

func init() {
	rootCmd.AddCommand(bestCmd)
	err := utils.RegisterFlags[bestFlags](bestCmd)
	if err != nil {
		panic(err)
	}
}
