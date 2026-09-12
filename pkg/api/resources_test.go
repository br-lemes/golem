package api

import "testing"

func TestResourcesPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/resources", ResourcesSize, func(client *Client) (int, error) {
		items, err := client.Resources(ResourcesOptions{})
		return len(items), err
	})
}
