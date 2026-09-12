package api

import "testing"

func TestEventsActivePagination(t *testing.T) {
	testPaginatedEndpoint(t, "/events/active", EventsActiveSize, func(client *Client) (int, error) {
		items, err := client.EventsActive()
		return len(items), err
	})
}
