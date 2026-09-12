package api

import "testing"

func TestItemsPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/items", ItemsSize, func(client *Client) (int, error) {
		items, err := client.Items(ItemsOptions{})
		return len(items), err
	})
}
