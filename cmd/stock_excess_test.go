package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
)

func TestStockExcess(t *testing.T) {
	stocks := []models.Stock{
		{
			Code:     "raw_chicken",
			Quantity: 0,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		{
			Code:     "iron_sword",
			Quantity: 5,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		{
			Code:     "golden_egg",
			Quantity: 0,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		{
			Code:     "wolf_ears",
			Quantity: 5,
			Mode:     models.StockModeExact,
			Enabled:  true,
		},
		{
			Code:     "ash_wood",
			Quantity: 100,
			Mode:     models.StockModeSafety,
			Enabled:  true,
		},
		{
			Code:     "copper",
			Quantity: 0,
			Mode:     models.StockModeExact,
			Enabled:  false,
		},
	}
	quantities := map[string]int{
		"raw_chicken": 4,
		"iron_sword":  7,
		"golden_egg":  2,
		"wolf_ears":   8,
		"ash_wood":    200,
		"copper":      2,
	}
	sellable := func(code, npc string) bool {
		return (code == "golden_egg" || code == "wolf_ears") && (npc == "" || npc == "all" || npc == "monk")
	}

	excess := stockExcess(stocks, quantities, false, "", sellable)
	if len(excess) != 1 || excess["raw_chicken"].Excess != 4 {
		t.Fatalf("stockExcess() = %#v", excess)
	}

	excess = stockExcess(stocks, quantities, false, "all", sellable)
	if len(excess) != 2 || excess["golden_egg"].Excess != 2 || excess["wolf_ears"].Excess != 3 {
		t.Fatalf("sell stockExcess() = %#v", excess)
	}

	excess = stockExcess(stocks, quantities, true, "", sellable)
	if len(excess) != 4 || excess["iron_sword"].Excess != 2 {
		t.Fatalf("all stockExcess() = %#v", excess)
	}
}
