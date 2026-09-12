package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionUnequip(t *testing.T) {
	q := 1
	testActionBody(t, actionBodyTestCase{
		name: "unequip",
		path: "/my/hero/action/unequip",
		body: []byte(`[{"quantity":1,"slot":"weapon"}]`),
		call: func(c *Client) error {
			_, err := c.MyActionUnequip("hero", []schemas.UnequipSchema{{
				Quantity: &q,
				Slot:     "weapon",
			}})
			return err
		},
	})
}

func TestMyActionUnequipReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "unequip",
		call: func(c *Client) error {
			_, err := c.MyActionUnequip("hero", nil)
			return err
		},
	}})
}
