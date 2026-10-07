package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var stockListCmd = &cobra.Command{
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Use:     "list",
	Short:   "List configured stock quantities",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return console.Auto(cache.ListStocks())
	},
}

func init() {
	stockCmd.AddCommand(stockListCmd)
}
