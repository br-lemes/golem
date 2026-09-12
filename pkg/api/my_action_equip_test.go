package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionEquip(t *testing.T) {
	q := 1
	testActionBody(t, actionBodyTestCase{
		name: "equip",
		path: "/my/hero/action/equip",
		body: []byte(`[{"code":"sword","quantity":1,"slot":"weapon"}]`),
		call: func(c *Client) error {
			_, err := c.MyActionEquip("hero", []schemas.EquipSchema{{
				Code:     "sword",
				Quantity: &q,
				Slot:     "weapon",
			}})
			return err
		},
	})
}

func TestMyActionEquipReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "equip",
		call: func(c *Client) error {
			_, err := c.MyActionEquip("hero", nil)
			return err
		},
	}})
}
