package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/stock"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type stockFlags struct {
	Target bool `flag:"target" desc:"Use safety stock targets when calculating requirements"`
}

var stockCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "stock",
	Short: "Show stock overview",
	RunE:  stockRun,
}

func stockRun(cmd *cobra.Command, args []string) error {
	flags, err := utils.ReadFlags[stockFlags](cmd)
	if err != nil {
		return err
	}
	cmd.SilenceUsage = true
	overview, err := stock.Overview(cache.ListStocks(), args, flags.Target)
	if err != nil {
		return err
	}
	return console.Auto(overview)
}

func init() {
	rootCmd.AddCommand(stockCmd)
	err := utils.RegisterFlags[stockFlags](stockCmd)
	if err != nil {
		panic(err)
	}
}
