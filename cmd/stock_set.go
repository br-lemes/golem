package cmd

import (
	"fmt"
	"strconv"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type stockSetFlags struct {
	Safety bool `flag:"safety" desc:"Use the quantity as a replenishment threshold"`
}

var stockSetCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "set <code> <quantity>",
	Short: "Set an item stock quantity",
	Long: `Set an item stock quantity

Arguments:
  code       The code of the item.
  quantity   The exact quantity to keep.`,
	ValidArgsFunction: completion.Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[stockSetFlags](cmd)
		if err != nil {
			return err
		}
		_, ok := catalog.Items().Get(args[0])
		if !ok {
			return fmt.Errorf("unknown item %q", args[0])
		}
		quantity, err := strconv.Atoi(args[1])
		if err != nil || quantity < 0 {
			return fmt.Errorf("quantity must be a non-negative integer")
		}
		if flags.Safety && quantity == 0 {
			return fmt.Errorf("safety quantity must be positive")
		}
		mode := models.StockModeExact
		if flags.Safety {
			mode = models.StockModeSafety
		}
		cmd.SilenceUsage = true
		return cache.SetStock(args[0], quantity, mode)
	},
}

func init() {
	stockCmd.AddCommand(stockSetCmd)
	err := utils.RegisterFlags[stockSetFlags](stockSetCmd)
	if err != nil {
		panic(err)
	}
}
