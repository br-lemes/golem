package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/spf13/cobra"
)

var stockDisableCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "disable <code>...",
	Short: "Disable item stock tracking",
	Long: `Disable item stock tracking

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
		for _, code := range args {
			err := cache.SetStockEnabled(code, false)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	stockCmd.AddCommand(stockDisableCmd)
}
