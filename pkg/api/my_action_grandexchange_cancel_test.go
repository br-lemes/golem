package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionGrandexchangeCancel(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "grandexchange cancel",
		path: "/my/hero/action/grandexchange/cancel",
		body: []byte(`{"id":"order"}`),
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeCancel("hero", schemas.GECancelOrderSchema{
				Id: "order",
			})
			return err
		},
	})
}

func TestMyActionGrandexchangeCancelReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "grandexchange cancel",
		call: func(c *Client) error {
			_, err := c.MyActionGrandexchangeCancel("hero", schemas.GECancelOrderSchema{})
			return err
		},
	}})
}
