package api

import "testing"

func TestGrandexchangeHistoryPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/grandexchange/history/iron", GrandexchangeHistorySize, func(client *Client) (int, error) {
		items, err := client.GrandexchangeHistory("iron", GrandexchangeHistoryOptions{})
		return len(items), err
	})
}
