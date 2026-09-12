package api

import "testing"

func TestAchievementsPagination(t *testing.T) {
	testPaginatedEndpoint(t, "/achievements", AchievementsSize, func(client *Client) (int, error) {
		items, err := client.Achievements(AchievementsOptions{})
		return len(items), err
	})
}
