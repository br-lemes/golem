package api

import "testing"

func TestNpcsDetailsPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/npcs/details", NpcsDetailsSize, func(client *Client) (int, error) {
		items, err := client.NpcsDetails(NpcsDetailsOptions{})
		return len(items), err
	})
}
