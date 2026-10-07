package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/spf13/cobra"
)

var stockDiscardCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "discard <code>...",
	Short: "Discard items from inventory",
	Long: `Discard items from inventory

Arguments:
  code   The code of the item.`,
	ValidArgsFunction: completion.Item(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, code := range args {
			_, ok := catalog.Items().Get(code)
			if !ok {
				return fmt.Errorf("unknown item %q", code)
			}
			stock, exists := cache.GetStock(code)
			if !exists || stock.Mode != models.StockModeExact || stock.Quantity != 0 {
				return fmt.Errorf("item %q must have an exact zero stock quantity", code)
			}
		}
		cmd.SilenceUsage = true
		for _, code := range args {
			err := cache.SetStockDiscard(code, true)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	stockCmd.AddCommand(stockDiscardCmd)
}
