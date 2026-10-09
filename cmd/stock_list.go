package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/spf13/cobra"
)

var stockListCmd = &cobra.Command{
	Aliases: []string{"ls"},
	Args:    cobra.ArbitraryArgs,
	Use:     "list [code...]",
	Short:   "List configured stock quantities",
	Long: `List configured stock quantities

Arguments:
  code   The code of the item.`,
	ValidArgsFunction: completion.Item(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, code := range args {
			_, ok := catalog.Items().Get(code)
			if !ok {
				return fmt.Errorf("unknown item %q", code)
			}
		}
		cmd.SilenceUsage = true
		stocks := cache.ListStocks()
		if len(args) > 0 {
			stocks = slices.DeleteFunc(stocks, func(stock models.Stock) bool {
				return !slices.Contains(args, stock.Code)
			})
		}
		return console.Auto(stocks)
	},
}

func init() {
	stockCmd.AddCommand(stockListCmd)
}
