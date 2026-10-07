package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type stockExcessFlags struct {
	All  bool   `flag:"all" desc:"Include all exact stock quantities"`
	Sell string `flag:"sell" desc:"Show sellable items; optionally limit by NPC"`
}

type ItemExcess struct {
	Current  int `json:"current"`
	Quantity int `json:"quantity"`
	Excess   int `json:"excess"`
}

var stockExcessCmd = &cobra.Command{
	Use:   "excess",
	Short: "Show excess stock items",
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
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		excess := stockExcess(cache.ListStocks(), stockItemQuantities(bankItems, characters), flags.All, flags.Sell, stockSellable)
		return console.Auto(excess)
	},
}

func stockSellable(code, npc string) bool {
	item, ok := catalog.NpcsItems.Get(code)
	return ok && item.SellPrice != nil && *item.SellPrice > 0 && (npc == "" || npc == "all" || item.Npc == npc)
}

func stockExcess(stocks []models.Stock, quantities map[string]int, all bool, sell string, sellable func(string, string) bool) map[string]ItemExcess {
	result := map[string]ItemExcess{}
	for _, stock := range stocks {
		if !stock.Enabled || stock.Mode != models.StockModeExact || (!all && sell == "" && stock.Quantity != 0) {
			continue
		}
		if (sell != "" && !sellable(stock.Code, sell)) || (sell == "" && !all && sellable(stock.Code, "")) {
			continue
		}
		current := quantities[stock.Code]
		if current <= stock.Quantity {
			continue
		}
		result[stock.Code] = ItemExcess{
			Current:  current,
			Quantity: stock.Quantity,
			Excess:   current - stock.Quantity,
		}
	}
	return result
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
