package api

import "testing"

func TestEventsPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/events", EventsSize, func(client *Client) (int, error) {
		items, err := client.Events(EventsOptions{})
		return len(items), err
	})
}
