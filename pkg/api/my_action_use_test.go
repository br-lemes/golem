package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionUse(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "use",
		path: "/my/hero/action/use",
		body: []byte(`{"code":"potion","quantity":1}`),
		call: func(c *Client) error {
			_, err := c.MyActionUse("hero", schemas.SimpleItemSchema{
				Code:     "potion",
				Quantity: 1,
			})
			return err
		},
	})
}

func TestMyActionUseReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "use",
		call: func(c *Client) error {
			_, err := c.MyActionUse("hero", schemas.SimpleItemSchema{})
			return err
		},
	}})
}
