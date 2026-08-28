package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestGatheringCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "gathering <name> <code>",
	Short: "Find the best equipment for gathering",
	Long: `Find the best equipment for gathering

Arguments:
  name   Name of your character.
  code   The code of the resource.`,
	ValidArgsFunction: completion.CharacterName(1).Resource(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]
		flags, err := utils.ReadFlags[bestFlags](cmd)
		if err != nil {
			return err
		}
		err = bestGatheringValidate(code)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return bestGatheringRun(name, code, flags)
	},
}

func bestGatheringValidate(code string) error {
	_, ok := database.Resources.Get(code)
	if !ok {
		return fmt.Errorf("resource not found: %s", code)
	}
	return nil
}

func bestGatheringRun(name, code string, flags bestFlags) error {
	resource, _ := database.Resources.Get(code)
	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	priorities := best.GatheringPriorities(character, resource)
	items, err := best.FindEquipment(character, best.EquipmentOptions{
		UniqueAdeptRing: !flags.AllowDuplicateAdeptRing,
		Priorities:      priorities,
	})
	if err != nil {
		return err
	}
	return console.Auto(items)
}

func init() {
	bestCmd.AddCommand(bestGatheringCmd)
	err := utils.RegisterFlags[bestFlags](bestGatheringCmd)
	if err != nil {
		panic(err)
	}
}
