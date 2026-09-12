package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionTaskTrade(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "task trade",
		path: "/my/hero/action/task/trade",
		body: []byte(`{"code":"iron","quantity":2}`),
		call: func(c *Client) error {
			_, err := c.MyActionTaskTrade("hero", schemas.SimpleItemSchema{
				Code:     "iron",
				Quantity: 2,
			})
			return err
		},
	})
}

func TestMyActionTaskTradeReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "task trade",
		call: func(c *Client) error {
			_, err := c.MyActionTaskTrade("hero", schemas.SimpleItemSchema{})
			return err
		},
	}})
}
