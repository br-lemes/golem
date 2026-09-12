package api

import (
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyBankItemsReturnsCachedItems(t *testing.T) {
	cache.CleanBankItems()
	t.Cleanup(cache.CleanBankItems)
	items := []schemas.SimpleItemSchema{{Code: "iron"}}
	cache.SaveBankItems(items)

	client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"request not expected"}}`)))

	got, err := client.MyBankItems()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Code != "iron" {
		t.Fatalf("got cached items %#v", got)
	}
}

func TestMyBankItemsPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/my/bank/items", MyBankItemsSize, func(client *Client) (int, error) {
		cache.CleanBankItems()
		items, err := client.MyBankItems()
		return len(items), err
	})
}
