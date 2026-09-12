package api

import "testing"

func TestMyActionGathering(t *testing.T) {
	testActionBody(t, actionBodyTestCase{
		name: "gathering",
		path: "/my/hero/action/gathering",
		call: func(c *Client) error {
			_, err := c.MyActionGathering("hero")
			return err
		},
	})
}

func TestMyActionGatheringReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "gathering",
		call: func(c *Client) error {
			_, err := c.MyActionGathering("hero")
			return err
		},
	}})
}
