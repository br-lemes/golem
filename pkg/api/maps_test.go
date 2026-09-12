package api

import "testing"

func TestMapsPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/maps", MapsSize, func(client *Client) (int, error) {
		items, err := client.Maps(MapsOptions{})
		return len(items), err
	})
}
