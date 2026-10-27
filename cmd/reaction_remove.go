package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/spf13/cobra"
)

var reactionRemoveCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "remove <code> [name...]",
	Short: "Remove an event reaction",
	Long: `Remove an event reaction

Arguments:
  code   The event content code.
  name   Names of your characters.`,
	ValidArgsFunction: completion.Custom(1, catalog.EventContentCodes).Custom(0, cache.GetCharacters).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateReactionCode(args[0])
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return cache.RemoveReactions(args[0], args[1:])
	},
}

func init() {
	reactionCmd.AddCommand(reactionRemoveCmd)
}
