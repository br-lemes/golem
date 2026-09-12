package api

import (
	"net/http"
	"testing"
)

func TestGrandexchangeOrderBuildsEscapedPath(t *testing.T) {
	tests := []struct {
		name string
		id   string
		path string
	}{
		{
			name: "plain id",
			id:   "123",
			path: "/grandexchange/orders/123",
		},
		{
			name: "escaped id",
			id:   "iron/ore",
			path: "/grandexchange/orders/iron%2Fore",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
				if r.URL.EscapedPath() != test.path {
					t.Errorf("path = %s, want %s", r.URL.EscapedPath(), test.path)
				}
				return []byte(`{"data":{}}`), nil
			}))

			_, err := client.GrandexchangeOrder(test.id)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGrandexchangeOrderReturnsRequestError(t *testing.T) {
	client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"order unavailable"}}`)))

	_, err := client.GrandexchangeOrder("123")
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestGrandexchangeOrderReturnsJSONError(t *testing.T) {
	client := newTestClient(responseTransport(http.StatusOK, []byte("invalid json")))

	_, err := client.GrandexchangeOrder("123")
	if err == nil {
		t.Fatal("expected JSON error")
	}
}

func TestGrandexchangeOrdersPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/grandexchange/orders", GrandexchangeOrdersSize, func(client *Client) (int, error) {
		items, err := client.GrandexchangeOrders(GrandexchangeOrdersOptions{})
		return len(items), err
	})
}
