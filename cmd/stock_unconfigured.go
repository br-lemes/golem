package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/stock"
	"github.com/spf13/cobra"
)

var stockUnconfiguredCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "unconfigured",
	Short: "List held items without configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		unconfigured, err := stock.Unconfigured(cache.ListStocks())
		if err != nil {
			return err
		}
		return console.Auto(unconfigured)
	},
}

func init() {
	stockCmd.AddCommand(stockUnconfiguredCmd)
}
