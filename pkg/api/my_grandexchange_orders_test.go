package api

import "testing"

func TestMyGrandexchangeOrdersPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/my/grandexchange/orders", MyGrandexchangeOrdersSize, func(client *Client) (int, error) {
		items, err := client.MyGrandexchangeOrders(MyGrandexchangeOrdersOptions{})
		return len(items), err
	})
}
