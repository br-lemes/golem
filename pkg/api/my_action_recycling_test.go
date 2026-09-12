package api

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionRecycling(t *testing.T) {
	q := 2
	testActionBody(t, actionBodyTestCase{
		name: "recycling",
		path: "/my/hero/action/recycling",
		body: []byte(`{"code":"iron","quantity":2}`),
		call: func(c *Client) error {
			_, err := c.MyActionRecycling("hero", schemas.RecyclingSchema{
				Code:     "iron",
				Quantity: &q,
			})
			return err
		},
	})
}

func TestMyActionRecyclingReturnsErrors(t *testing.T) {
	testActionErrors(t, []actionTestCase{{
		name: "recycling",
		call: func(c *Client) error {
			_, err := c.MyActionRecycling("hero", schemas.RecyclingSchema{})
			return err
		},
	}})
}
