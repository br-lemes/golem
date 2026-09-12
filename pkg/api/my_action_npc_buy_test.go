package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionNPCBuy(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "npc buy",
		path: "/my/hero/action/npc/buy",
		body: []byte(`{"code":"potion","quantity":1}`),
		call: func(c *Client) error {
			_, err := c.MyActionNPCBuy("hero", schemas.SimpleItemSchema{
				Code:     "potion",
				Quantity: 1,
			})
			return err
		},
	})
}

func TestMyActionNPCBuyReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "npc buy",
		call: func(c *Client) error {
			_, err := c.MyActionNPCBuy("hero", schemas.SimpleItemSchema{})
			return err
		},
	}})
}
