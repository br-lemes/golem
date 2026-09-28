package api

import "github.com/br-lemes/golem/pkg/cache"

func HasAchievement(code string) (bool, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.HasAchievement(code)
}

func (c *Client) HasAchievement(code string) (bool, error) {
	account, err := c.accountName()
	if err != nil {
		return false, err
	}
	if cache.HasAchievement(code) {
		return true, nil
	}
	achievements, err := c.AccountsAchievements(account, AccountsAchievementsOptions{
		Completed: true,
	})
	if err != nil {
		return false, err
	}
	codes := make([]string, 0, len(achievements))
	for _, achievement := range achievements {
		codes = append(codes, achievement.Code)
	}
	cache.SaveAchievements(codes)
	return cache.HasAchievement(code), nil
}
