package api

import "testing"

func TestEffectsPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/effects", EffectsSize, func(client *Client) (int, error) {
		items, err := client.Effects(EffectsOptions{})
		return len(items), err
	})
}
