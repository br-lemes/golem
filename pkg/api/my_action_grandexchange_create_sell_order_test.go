package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionGrandexchangeCreateSellOrder(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "create sell order",
		path: "/my/hero/action/grandexchange/create_sell_order",
		body: []byte(`{"code":"iron","price":10,"quantity":2}`),
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeCreateSellOrder("hero", schemas.GEOrderCreationSchema{
				Code:     "iron",
				Price:    10,
				Quantity: 2,
			})
			return err
		},
	})
}

func TestMyActionGrandexchangeCreateSellOrderReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "create sell order",
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeCreateSellOrder("hero", schemas.GEOrderCreationSchema{})
			return err
		},
	}})
}
