package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/spf13/cobra"
)

var reactionEnableCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "enable <code> [name...]",
	Short: "Enable an event reaction",
	Long: `Enable an event reaction

Arguments:
  code   The event content code.
  name   Names of your character.`,
	ValidArgsFunction: completion.Custom(1, catalog.EventContentCodes).Custom(0, cache.GetCharacters).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateReactionCode(args[0])
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return cache.SetReactionsEnabled(args[0], args[1:], true)
	},
}

func init() {
	reactionCmd.AddCommand(reactionEnableCmd)
}
