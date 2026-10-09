package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
)

func TestStockUnconfigured(t *testing.T) {
	stocks := []models.Stock{
		{Code: "exact", Quantity: 5, Mode: models.StockModeExact},
		{Code: "zero", Quantity: 0, Mode: models.StockModeExact},
		{Code: "safety", Quantity: 10, Mode: models.StockModeSafety},
	}
	quantities := map[string]int{
		"":       1,
		"exact":  10,
		"zero":   2,
		"safety": 5,
		"other":  3,
		"empty":  0,
	}

	result := stockUnconfigured(stocks, quantities)
	if len(result) != 1 || result["other"] != 3 {
		t.Fatalf("stockUnconfigured() = %#v", result)
	}
}
