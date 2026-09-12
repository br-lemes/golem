package api

import "testing"

func TestMyActionRest(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "rest",
		path: "/my/hero/action/rest",
		call: func(c *Client) error {
			_, err := c.MyActionRest("hero")
			return err
		},
	})
}

func TestMyActionRestReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "rest",
		call: func(c *Client) error {
			_, err := c.MyActionRest("hero")
			return err
		},
	}})
}
