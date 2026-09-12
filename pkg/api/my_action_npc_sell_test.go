package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionNPCSell(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "npc sell",
		path: "/my/hero/action/npc/sell",
		body: []byte(`{"code":"potion","quantity":1}`),
		call: func(c *Client) error {
			_, err := c.MyActionNPCSell("hero", schemas.SimpleItemSchema{
				Code:     "potion",
				Quantity: 1,
			})
			return err
		},
	})
}

func TestMyActionNPCSellReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "npc sell",
		call: func(c *Client) error {
			_, err := c.MyActionNPCSell("hero", schemas.SimpleItemSchema{})
			return err
		},
	}})
}
