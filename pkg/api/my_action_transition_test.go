package api

import "testing"

func TestMyActionTransition(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "transition",
		path: "/my/hero/action/transition",
		call: func(c *Client) error {
			_, err := c.MyActionTransition("hero")
			return err
		},
	})
}

func TestMyActionTransitionReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "transition",
		call: func(c *Client) error {
			_, err := c.MyActionTransition("hero")
			return err
		},
	}})
}
