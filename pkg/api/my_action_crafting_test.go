package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionCrafting(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "crafting",
		path: "/my/hero/action/crafting",
		body: []byte(`{"code":"iron","quantity":1}`),
		call: func(c *Client) error {
			_, err := c.MyActionCrafting("hero", schemas.SimpleItemSchema{
				Code:     "iron",
				Quantity: 1,
			})
			return err
		},
	})
}

func TestMyActionCraftingReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "crafting",
		call: func(c *Client) error {
			_, err := c.MyActionCrafting("hero", schemas.SimpleItemSchema{})
			return err
		},
	}})
}
