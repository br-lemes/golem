package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/spf13/cobra"
)

var stockUnconfiguredCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "unconfigured",
	Short: "List held items without stock configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		return console.Auto(stockUnconfigured(cache.ListStocks(), stockItemQuantities(bankItems, characters)))
	},
}

func stockUnconfigured(stocks []models.Stock, quantities map[string]int) map[string]int {
	configured := stockConfigured(stocks)
	result := map[string]int{}
	for code, quantity := range quantities {
		if code != "" && quantity > 0 {
			_, found := configured[code]
			if !found {
				result[code] = quantity
			}
		}
	}
	return result
}

func init() {
	stockCmd.AddCommand(stockUnconfiguredCmd)
}
