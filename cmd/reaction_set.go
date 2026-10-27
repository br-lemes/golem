package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type reactionSetFlags struct {
	Priority int  `flag:"priority" desc:"Priority of the reaction; lower values run first"`
	Disabled bool `flag:"disabled" desc:"Disable the reaction"`
}

var reactionSetCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "set <code> [name...]",
	Short: "Set an event reaction",
	Long: `Set an event reaction

Arguments:
  code   The event content code.
  name   Names of your character.`,
	ValidArgsFunction: completion.Custom(1, catalog.EventContentCodes).Custom(0, cache.GetCharacters).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[reactionSetFlags](cmd)
		if err != nil {
			return err
		}
		err = validateReactionCode(args[0])
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		if len(args[1:]) != 1 && cmd.Flags().Changed("priority") {
			return fmt.Errorf("priority can only be specified for one character")
		}
		var priority *int
		if cmd.Flags().Changed("priority") {
			priority = &flags.Priority
		}
		var enabled *bool
		if cmd.Flags().Changed("disabled") {
			value := !flags.Disabled
			enabled = &value
		}
		names := args[1:]
		if len(names) == 0 {
			names = cache.GetCharacters()
		}
		for _, name := range names {
			err = cache.SetReaction(name, args[0], priority, enabled)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	reactionCmd.AddCommand(reactionSetCmd)
	err := utils.RegisterFlags[reactionSetFlags](reactionSetCmd)
	if err != nil {
		panic(err)
	}
}
