package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type reactionListFlags struct {
	Name []string `flag:"name" desc:"Filter by character names"`
	Code []string `flag:"code" desc:"Filter by event content codes"`
}

var reactionListCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "list",
	Short: "List event reactions",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[reactionListFlags](cmd)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return console.Auto(cache.ListReactions(flags.Name, flags.Code))
	},
}

func init() {
	reactionCmd.AddCommand(reactionListCmd)
	err := utils.RegisterFlags[reactionListFlags](reactionListCmd)
	if err != nil {
		panic(err)
	}
	err = reactionListCmd.RegisterFlagCompletionFunc("name", completion.StringSlice(cache.GetCharacters))
	if err != nil {
		panic(err)
	}
	err = reactionListCmd.RegisterFlagCompletionFunc("code", completion.StringSlice(catalog.EventContentCodes))
	if err != nil {
		panic(err)
	}
}
