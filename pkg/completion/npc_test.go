package completion

import "testing"

func TestNPCBuyAndSellItems(t *testing.T) {
	buy := GetNPCBuyItems()
	sell := GetNPCSellItems()
	if len(buy) == 0 {
		t.Fatal("NPC buy completion returned no items")
	}
	if len(sell) == 0 {
		t.Fatal("NPC sell completion returned no items")
	}
	buySuggestions, _ := NPCBuy(1).Build()(nil, nil, "")
	sellSuggestions, _ := NPCSell(1).Build()(nil, nil, "")
	if len(buySuggestions) == 0 || len(sellSuggestions) == 0 {
		t.Fatal("NPC completion builders were not created")
	}
}
