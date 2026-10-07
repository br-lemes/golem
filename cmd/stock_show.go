package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var stockShowCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "show <code>...",
	Short: "Show requirements for stock items",
	RunE:  stockRun,

	ValidArgsFunction: completion.Item(0).Build(),
}

func init() {
	stockCmd.AddCommand(stockShowCmd)
	err := utils.RegisterFlags[stockFlags](stockShowCmd)
	if err != nil {
		panic(err)
	}
}
