package api

import "testing"

func TestMyActionMove(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "move",
		path: "/my/hero/action/move",
		body: []byte(`{"x":3,"y":4}`),
		call: func(c *Client) error {
			_, err := c.MyActionMove("hero", 3, 4)
			return err
		},
	})
}

func TestMyActionMoveReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "move",
		call: func(c *Client) error {
			_, err := c.MyActionMove("hero", 1, 2)
			return err
		},
	}})
}
