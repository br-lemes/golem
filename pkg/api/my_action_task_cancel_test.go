package api

import "testing"

func TestMyActionTaskCancel(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "task cancel",
		path: "/my/hero/action/task/cancel",
		call: func(c *Client) error {
			_, err := c.MyActionTaskCancel("hero")
			return err
		},
	})
}

func TestMyActionTaskCancelReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "task cancel",
		call: func(c *Client) error {
			_, err := c.MyActionTaskCancel("hero")
			return err
		},
	}})
}
