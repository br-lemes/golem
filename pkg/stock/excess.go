package stock

import (
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/models"
)

type ItemExcess struct {
	Current  int `json:"current"`
	Quantity int `json:"quantity"`
	Excess   int `json:"excess"`
}

func Excess(stocks []models.Stock, all bool, sell string) (map[string]ItemExcess, error) {
	quantities, err := heldQuantities()
	if err != nil {
		return nil, err
	}
	return excess(stocks, quantities, all, sell, sellable), nil
}

func excess(stocks []models.Stock, quantities map[string]int, all bool, sell string, sellable func(string, string) bool) map[string]ItemExcess {
	result := map[string]ItemExcess{}
	for _, stock := range stocks {
		if !stock.Enabled || stock.Mode != models.StockModeExact || (!all && sell == "" && stock.Quantity != 0) {
			continue
		}
		if (sell != "" && !sellable(stock.Code, sell)) || (sell == "" && !all && sellable(stock.Code, "")) {
			continue
		}
		current := quantities[stock.Code]
		if current > stock.Quantity {
			result[stock.Code] = ItemExcess{Current: current, Quantity: stock.Quantity, Excess: current - stock.Quantity}
		}
	}
	return result
}

func sellable(code, npc string) bool {
	item, ok := catalog.NpcsItems.Get(code)
	return ok && item.SellPrice != nil && *item.SellPrice > 0 && (npc == "" || npc == "all" || item.Npc == npc)
}
