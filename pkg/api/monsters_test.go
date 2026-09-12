package api

import "testing"

func TestMonstersPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/monsters", MonstersSize, func(client *Client) (int, error) {
		items, err := client.Monsters(MonstersOptions{})
		return len(items), err
	})
}
