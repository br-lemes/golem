package api

import "testing"

func TestMyActionTaskComplete(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "task complete",
		path: "/my/hero/action/task/complete",
		call: func(c *Client) error {
			_, err := c.MyActionTaskComplete("hero")
			return err
		},
	})
}

func TestMyActionTaskCompleteReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "task complete",
		call: func(c *Client) error {
			_, err := c.MyActionTaskComplete("hero")
			return err
		},
	}})
}
