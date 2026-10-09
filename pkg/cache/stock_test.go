package cache

import (
	"slices"
	"testing"

	"github.com/br-lemes/golem/pkg/models"
)

func TestStockCache(t *testing.T) {
	initializeTestCache(t)
	err := SetStock("ash_wood", 100, models.StockModeSafety)
	if err != nil {
		t.Fatal(err)
	}
	err = SetStock("raw_chicken", 0, models.StockModeExact)
	if err != nil {
		t.Fatal(err)
	}
	err = SetStockEnabled("copper", false)
	if err != nil {
		t.Fatal(err)
	}
	err = SetStockEnabled("ash_wood", false)
	if err != nil {
		t.Fatal(err)
	}
	err = SetStock("ash_wood", 200, models.StockModeSafety)
	if err != nil {
		t.Fatal(err)
	}
	err = SetStockDiscard("raw_chicken", true)
	if err != nil {
		t.Fatal(err)
	}
	err = SetStock("raw_chicken", 1, models.StockModeExact)
	if err != nil {
		t.Fatal(err)
	}
	err = SetStock("raw_chicken", 0, models.StockModeExact)
	if err != nil {
		t.Fatal(err)
	}
	err = RemoveStock("copper")
	if err != nil {
		t.Fatal(err)
	}

	stocks := ListStocks()
	if len(stocks) != 2 {
		t.Fatalf("ListStocks() = %#v", stocks)
	}
	if stocks[0].Code != "ash_wood" || stocks[0].Quantity != 200 || stocks[0].Mode != models.StockModeSafety || !stocks[0].Enabled {
		t.Fatalf("ash_wood stock = %#v", stocks[0])
	}
	if stocks[1].Code != "raw_chicken" || stocks[1].Quantity != 0 || stocks[1].Mode != models.StockModeExact || !stocks[1].Enabled || stocks[1].Discard {
		t.Fatalf("raw_chicken stock = %#v", stocks[1])
	}
}

func TestExactZeroStockCodes(t *testing.T) {
	initializeTestCache(t)
	for _, stock := range []models.Stock{
		{
			Code:     "enabled_zero",
			Quantity: 0,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		{
			Code:     "disabled_zero",
			Quantity: 0,
			Mode:     models.StockModeExact,
		},
		{
			Code:     "enabled_one",
			Quantity: 1,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		{
			Code:     "safety_zero",
			Quantity: 0,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
	} {
		err := cache.Create(&stock).Error
		if err != nil {
			t.Fatal(err)
		}
	}
	codes := ExactZeroStockCodes()
	if !slices.Equal(codes, []string{"enabled_zero"}) {
		t.Fatalf("ExactZeroStockCodes() = %#v", codes)
	}
}
