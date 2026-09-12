package api

import "testing"

func TestMyActionTaskNew(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "task new",
		path: "/my/hero/action/task/new",
		call: func(c *Client) error {
			_, err := c.MyActionTaskNew("hero")
			return err
		},
	})
}

func TestMyActionTaskNewReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "task new",
		call: func(c *Client) error {
			_, err := c.MyActionTaskNew("hero")
			return err
		},
	}})
}
