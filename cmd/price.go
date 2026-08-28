package cmd

import (
	"fmt"
	"sort"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

type priceResult struct {
	Code        string              `json:"code"`
	Buy         priceMarketSummary  `json:"buy"`
	Sell        priceMarketSummary  `json:"sell"`
	History     priceHistorySummary `json:"history"`
	Suggestions []string            `json:"suggestions,omitempty"`
}

type priceMarketSummary struct {
	Orders   int         `json:"orders"`
	Quantity int         `json:"quantity"`
	Prices   *priceStats `json:"prices"`
}

type priceHistorySummary struct {
	Sales  int         `json:"sales"`
	Prices *priceStats `json:"prices"`
}

type priceStats struct {
	Min     int `json:"min"`
	Average int `json:"average"`
	Max     int `json:"max"`
}

var priceCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "price <item>",
	Short: "Estimate useful Grand Exchange buy and sell prices",

	ValidArgsFunction: completion.Tradeables(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		includeOwnOrders, err := cmd.Flags().GetBool("include-own-orders")
		if err != nil {
			return err
		}
		item, exists := database.Items().Tradeables().Get(code)
		if !exists {
			return fmt.Errorf("item %q not found", code)
		}
		orders, err := api.GrandexchangeOrders(api.GrandexchangeOrdersOptions{
			Code: code,
		})
		if err != nil {
			return err
		}
		if !includeOwnOrders {
			ownOrders, err := api.MyGrandexchangeOrders(api.MyGrandexchangeOrdersOptions{
				Code: code,
			})
			if err != nil {
				return err
			}
			orders = excludeOrders(orders, ownOrders)
		}
		var buys, sells []schemas.GEOrderSchema
		for _, order := range orders {
			switch order.Type {
			case schemas.Buy:
				buys = append(buys, order)
			case schemas.Sell:
				sells = append(sells, order)
			}
		}
		history, err := api.GrandexchangeHistory(code, api.GrandexchangeHistoryOptions{})
		if err != nil {
			return err
		}
		result := priceResult{
			Code:    code,
			Buy:     summarizeOrders(buys),
			Sell:    summarizeOrders(sells),
			History: summarizeHistory(history),
		}
		if result.Buy.Orders == 0 && result.Sell.Orders == 0 && result.History.Sales == 0 {
			result.Suggestions = similarPriceItems(item, database.Items().Tradeables().All())
		}
		return console.Auto(result)
	},
}

func summarizeOrders(orders []schemas.GEOrderSchema) priceMarketSummary {
	result := priceMarketSummary{Orders: len(orders)}
	if len(orders) == 0 {
		return result
	}
	prices := make([]int, 0, len(orders))
	total := 0
	for _, order := range orders {
		prices = append(prices, order.Price)
		result.Quantity += order.Quantity
		total += order.Price * order.Quantity
	}
	sort.Ints(prices)
	result.Prices = &priceStats{
		Min:     prices[0],
		Average: total / result.Quantity,
		Max:     prices[len(prices)-1],
	}
	return result
}

func excludeOrders(orders, excluded []schemas.GEOrderSchema) []schemas.GEOrderSchema {
	excludedIDs := make(map[string]struct{}, len(excluded))
	for _, order := range excluded {
		excludedIDs[order.Id] = struct{}{}
	}
	filtered := make([]schemas.GEOrderSchema, 0, len(orders))
	for _, order := range orders {
		_, excluded := excludedIDs[order.Id]
		if !excluded {
			filtered = append(filtered, order)
		}
	}
	return filtered
}

func summarizeHistory(history []schemas.GEOrderHistorySchema) priceHistorySummary {
	result := priceHistorySummary{Sales: len(history)}
	if len(history) == 0 {
		return result
	}
	prices := make([]int, 0, len(history))
	total, quantity := 0, 0
	for _, sale := range history {
		prices = append(prices, sale.Price)
		total += sale.Price * sale.Quantity
		quantity += sale.Quantity
	}
	sort.Ints(prices)
	if quantity > 0 {
		result.Prices = &priceStats{
			Min:     prices[0],
			Average: total / quantity,
			Max:     prices[len(prices)-1],
		}
	}
	return result
}

func similarPriceItems(item *schemas.ItemSchema, items []*schemas.ItemSchema) []string {
	var candidates []string
	for _, other := range items {
		if other.Code == item.Code || other.Type != item.Type || other.Subtype != item.Subtype || other.Level != item.Level {
			continue
		}
		candidates = append(candidates, other.Code)
	}
	sort.Strings(candidates)
	return candidates
}

func init() {
	rootCmd.AddCommand(priceCmd)
	priceCmd.Flags().Bool("include-own-orders", false, "Include your own open orders in the market statistics.")
}
