package api

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

type actionTestCase struct {
	name string
	call func() error
}

func testActionErrors(t *testing.T, tests []actionTestCase) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.name+" request error", func(t *testing.T) {
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"action failed"}}`))
			err := test.call()
			if err == nil {
				t.Fatal("expected request error")
			}
		})
		t.Run(test.name+" JSON error", func(t *testing.T) {
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = responseClient(http.StatusOK, []byte("invalid json"))
			err := test.call()
			if err == nil {
				t.Fatal("expected JSON error")
			}
		})
	}
}

func TestActionsWithoutPayload(t *testing.T) {
	tests := []struct {
		name string
		path string
		call func() error
	}{
		{
			name: "gathering",
			path: "/my/hero/action/gathering",
			call: func() error { _, err := MyActionGathering("hero"); return err },
		},
		{
			name: "rest",
			path: "/my/hero/action/rest",
			call: func() error { _, err := MyActionRest("hero"); return err },
		},
		{
			name: "transition",
			path: "/my/hero/action/transition",
			call: func() error { _, err := MyActionTransition("hero"); return err },
		},
		{
			name: "task new",
			path: "/my/hero/action/task/new",
			call: func() error { _, err := MyActionTaskNew("hero"); return err },
		},
		{
			name: "task complete",
			path: "/my/hero/action/task/complete",
			call: func() error { _, err := MyActionTaskComplete("hero"); return err },
		},
		{
			name: "task cancel",
			path: "/my/hero/action/task/cancel",
			call: func() error { _, err := MyActionTaskCancel("hero"); return err },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache.CleanCharacters()
			t.Cleanup(cache.CleanCharacters)
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = testClient(func(r *http.Request) ([]byte, error) {
				if r.Method != http.MethodPost || r.URL.Path != test.path {
					t.Errorf("request = %s %s, want POST %s", r.Method, r.URL.Path, test.path)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil || len(body) != 0 {
					t.Errorf("request body = %q, want empty", body)
				}
				return []byte(`{"data":{}}`), nil
			})

			err := test.call()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestActionsWithoutPayloadReturnErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{
		{
			name: "gathering",
			call: func() error { _, err := MyActionGathering("hero"); return err },
		},
		{
			name: "rest",
			call: func() error { _, err := MyActionRest("hero"); return err },
		},
		{
			name: "transition",
			call: func() error { _, err := MyActionTransition("hero"); return err },
		},
		{
			name: "task new",
			call: func() error { _, err := MyActionTaskNew("hero"); return err },
		},
		{
			name: "task complete",
			call: func() error { _, err := MyActionTaskComplete("hero"); return err },
		},
		{
			name: "task cancel",
			call: func() error { _, err := MyActionTaskCancel("hero"); return err },
		},
	})
}

func TestActionsWithItemPayload(t *testing.T) {
	quantity := 2
	tests := []actionBodyTestCase{
		{
			name: "crafting",
			path: "/my/hero/action/crafting",
			body: []byte(`{"code":"iron","quantity":1}`),
			call: func() error {
				item := schemas.SimpleItemSchema{Code: "iron", Quantity: 1}
				_, err := MyActionCrafting("hero", item)
				return err
			},
		},
		{
			name: "use",
			path: "/my/hero/action/use",
			body: []byte(`{"code":"potion","quantity":1}`),
			call: func() error {
				item := schemas.SimpleItemSchema{Code: "potion", Quantity: 1}
				_, err := MyActionUse("hero", item)
				return err
			},
		},
		{
			name: "npc buy",
			path: "/my/hero/action/npc/buy",
			body: []byte(`{"code":"potion","quantity":1}`),
			call: func() error {
				item := schemas.SimpleItemSchema{Code: "potion", Quantity: 1}
				_, err := MyActionNPCBuy("hero", item)
				return err
			},
		},
		{
			name: "npc sell",
			path: "/my/hero/action/npc/sell",
			body: []byte(`{"code":"potion","quantity":1}`),
			call: func() error {
				item := schemas.SimpleItemSchema{Code: "potion", Quantity: 1}
				_, err := MyActionNPCSell("hero", item)
				return err
			},
		},
		{
			name: "recycling",
			path: "/my/hero/action/recycling",
			body: []byte(`{"code":"iron","quantity":2}`),
			call: func() error {
				item := schemas.RecyclingSchema{
					Code:     "iron",
					Quantity: &quantity,
				}
				_, err := MyActionRecycling("hero", item)
				return err
			},
		},
	}
	testActionsWithBody(t, tests)
}

func TestActionsWithItemPayloadReturnErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{
		{
			name: "crafting",
			call: func() error { _, err := MyActionCrafting("hero", schemas.SimpleItemSchema{}); return err },
		},
		{
			name: "use",
			call: func() error { _, err := MyActionUse("hero", schemas.SimpleItemSchema{}); return err },
		},
		{
			name: "npc buy",
			call: func() error { _, err := MyActionNPCBuy("hero", schemas.SimpleItemSchema{}); return err },
		},
		{
			name: "npc sell",
			call: func() error { _, err := MyActionNPCSell("hero", schemas.SimpleItemSchema{}); return err },
		},
		{
			name: "recycling",
			call: func() error { _, err := MyActionRecycling("hero", schemas.RecyclingSchema{}); return err },
		},
	})
}

func TestActionsWithSpecificPayload(t *testing.T) {
	tests := []actionBodyTestCase{
		{
			name: "move",
			path: "/my/hero/action/move",
			body: []byte(`{"x":3,"y":4}`),
			call: func() error { _, err := MyActionMove("hero", 3, 4); return err },
		},
		{
			name: "task trade",
			path: "/my/hero/action/task/trade",
			body: []byte(`{"code":"iron","quantity":2}`),
			call: func() error {
				item := schemas.SimpleItemSchema{Code: "iron", Quantity: 2}
				_, err := MyActionTaskTrade("hero", item)
				return err
			},
		},
		{
			name: "grandexchange buy",
			path: "/my/hero/action/grandexchange/buy",
			body: []byte(`{"id":"order","quantity":2}`),
			call: func() error {
				order := schemas.GEBuyOrderSchema{Id: "order", Quantity: 2}
				_, err := MyActionGrandexchangeBuy("hero", order)
				return err
			},
		},
		{
			name: "grandexchange cancel",
			path: "/my/hero/action/grandexchange/cancel",
			body: []byte(`{"id":"order"}`),
			call: func() error {
				order := schemas.GECancelOrderSchema{Id: "order"}
				_, err := MyActionGrandexchangeCancel("hero", order)
				return err
			},
		},
		{
			name: "create sell order",
			path: "/my/hero/action/grandexchange/create_sell_order",
			body: []byte(`{"code":"iron","price":10,"quantity":2}`),
			call: func() error {
				order := schemas.GEOrderCreationSchema{
					Code:     "iron",
					Price:    10,
					Quantity: 2,
				}
				_, err := MyActionGrandexchangeCreateSellOrder("hero", order)
				return err
			},
		},
		{
			name: "fill buy order",
			path: "/my/hero/action/grandexchange/fill",
			body: []byte(`{"id":"order","quantity":2}`),
			call: func() error {
				order := schemas.GEFillBuyOrderSchema{Id: "order", Quantity: 2}
				_, err := MyActionGrandexchangeFill("hero", order)
				return err
			},
		},
	}
	testActionsWithBody(t, tests)
}

func TestActionsWithSpecificPayloadReturnErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{
		{
			name: "move",
			call: func() error { _, err := MyActionMove("hero", 1, 2); return err },
		},
		{
			name: "task trade",
			call: func() error { _, err := MyActionTaskTrade("hero", schemas.SimpleItemSchema{}); return err },
		},
		{
			name: "grandexchange buy",
			call: func() error { _, err := MyActionGrandexchangeBuy("hero", schemas.GEBuyOrderSchema{}); return err },
		},
		{
			name: "grandexchange cancel",
			call: func() error { _, err := MyActionGrandexchangeCancel("hero", schemas.GECancelOrderSchema{}); return err },
		},
		{
			name: "create sell order",
			call: func() error {
				_, err := MyActionGrandexchangeCreateSellOrder("hero", schemas.GEOrderCreationSchema{})
				return err
			},
		},
		{
			name: "fill buy order",
			call: func() error { _, err := MyActionGrandexchangeFill("hero", schemas.GEFillBuyOrderSchema{}); return err },
		},
	})
}

type actionBodyTestCase struct {
	name string
	path string
	body []byte
	call func() error
}

func testActionsWithBody(t *testing.T, tests []actionBodyTestCase) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache.CleanCharacters()
			t.Cleanup(cache.CleanCharacters)
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = testClient(func(r *http.Request) ([]byte, error) {
				if r.Method != http.MethodPost || r.URL.Path != test.path {
					t.Errorf("request = %s %s, want POST %s", r.Method, r.URL.Path, test.path)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(body, test.body) {
					t.Errorf("body = %s, want %s", body, test.body)
				}
				return []byte(`{"data":{}}`), nil
			})
			err := test.call()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestActionsWithEquipmentPayload(t *testing.T) {
	quantity := 1
	tests := []actionBodyTestCase{
		{
			name: "equip",
			path: "/my/hero/action/equip",
			body: []byte(`[{"code":"sword","quantity":1,"slot":"weapon"}]`),
			call: func() error {
				items := []schemas.EquipSchema{
					{Code: "sword", Quantity: &quantity, Slot: "weapon"},
				}
				_, err := MyActionEquip("hero", items)
				return err
			},
		},
		{
			name: "unequip",
			path: "/my/hero/action/unequip",
			body: []byte(`[{"quantity":1,"slot":"weapon"}]`),
			call: func() error {
				items := []schemas.UnequipSchema{
					{Quantity: &quantity, Slot: "weapon"},
				}
				_, err := MyActionUnequip("hero", items)
				return err
			},
		},
	}
	testActionsWithBody(t, tests)
}

func TestActionsWithEquipmentPayloadReturnErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{
		{
			name: "equip",
			call: func() error {
				_, err := MyActionEquip("hero", nil)
				return err
			},
		},
		{
			name: "unequip",
			call: func() error {
				_, err := MyActionUnequip("hero", nil)
				return err
			},
		},
	})
}
