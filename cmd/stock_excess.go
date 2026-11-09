package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/stock"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type stockExcessFlags struct {
	All  bool   `flag:"all" desc:"Include all exact stock quantities"`
	Sell string `flag:"sell" desc:"Include sellable items; optionally limit by NPC"`
}

var stockExcessCmd = &cobra.Command{
	Use:   "excess",
	Short: "List excess stock items",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[stockExcessFlags](cmd)
		if err != nil {
			return err
		}
		if flags.Sell != "" && flags.Sell != "all" {
			_, ok := catalog.NpcsDetails.Get(flags.Sell)
			if !ok {
				return fmt.Errorf("unknown NPC %q", flags.Sell)
			}
		}
		cmd.SilenceUsage = true
		excess, err := stock.Excess(cache.ListStocks(), flags.All, flags.Sell)
		if err != nil {
			return err
		}
		return console.Auto(excess)
	},
}

func init() {
	stockCmd.AddCommand(stockExcessCmd)
	err := utils.RegisterFlags[stockExcessFlags](stockExcessCmd)
	if err != nil {
		panic(err)
	}
	stockExcessCmd.Flags().Lookup("sell").NoOptDefVal = "all"
	err = stockExcessCmd.RegisterFlagCompletionFunc("sell", func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return append([]string{"all"}, catalog.NpcsDetails.Keys()...), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}
