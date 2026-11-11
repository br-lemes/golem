package stock

import (
	"github.com/br-lemes/golem/pkg/models"
)

func Unconfigured(stocks []models.Stock) (map[string]int, error) {
	quantities, err := heldQuantities()
	if err != nil {
		return nil, err
	}
	return unconfigured(stocks, quantities), nil
}

func unconfigured(stocks []models.Stock, quantities map[string]int) map[string]int {
	configured := configured(stocks)
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
