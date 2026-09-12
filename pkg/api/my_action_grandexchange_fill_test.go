package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionGrandexchangeFill(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "fill buy order",
		path: "/my/hero/action/grandexchange/fill",
		body: []byte(`{"id":"order","quantity":2}`),
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeFill("hero", schemas.GEFillBuyOrderSchema{
				Id:       "order",
				Quantity: 2,
			})
			return err
		},
	})
}

func TestMyActionGrandexchangeFillReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "fill buy order",
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeFill("hero", schemas.GEFillBuyOrderSchema{})
			return err
		},
	}})
}
