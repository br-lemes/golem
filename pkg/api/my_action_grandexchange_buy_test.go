package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionGrandexchangeBuy(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "grandexchange buy",
		path: "/my/hero/action/grandexchange/buy",
		body: []byte(`{"id":"order","quantity":2}`),
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeBuy("hero", schemas.GEBuyOrderSchema{
				Id:       "order",
				Quantity: 2,
			})
			return err
		},
	})
}

func TestMyActionGrandexchangeBuyReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "grandexchange buy",
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeBuy("hero", schemas.GEBuyOrderSchema{})
			return err
		},
	}})
}
