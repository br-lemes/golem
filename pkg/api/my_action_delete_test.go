package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionDelete(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "delete",
		path: "/my/hero/action/delete",
		body: []byte(`{"code":"potion","quantity":1}`),
		call: func(c *Client) error {
			_, err := c.MyActionDelete("hero", schemas.SimpleItemSchema{
				Code:     "potion",
				Quantity: 1,
			})
			return err
		},
	})
}

func TestMyActionDeleteReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "delete",
		call: func(c *Client) error {
			_, err := c.MyActionDelete("hero", schemas.SimpleItemSchema{})
			return err
		},
	}})
}
