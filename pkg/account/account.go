package account

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

var (
	Bank      = api.MyBank
	BankItems = api.MyBankItems
	Details   = api.MyDetails
)

var Achievements = func() ([]schemas.AccountAchievementSchema, error) {
	return api.AccountsAchievements("", api.AccountsAchievementsOptions{})
}

var Characters = func() ([]schemas.CharacterSchema, error) {
	return api.AccountsCharacters("")
}

var HasAchievement = func(code string) (bool, error) {
	if cache.HasAchievement(code) {
		return true, nil
	}
	achievements, err := api.AccountsAchievements("", api.AccountsAchievementsOptions{
		Completed: true,
	})
	if err != nil {
		return false, err
	}
	found := false
	codes := make([]string, 0, len(achievements))
	for _, achievement := range achievements {
		codes = append(codes, achievement.Code)
		if achievement.Code == code {
			found = true
		}
	}
	cache.SaveAchievements(codes)
	return found, nil
}

var Orders = func() ([]schemas.GEOrderSchema, error) {
	return api.MyGrandexchangeOrders(api.MyGrandexchangeOrdersOptions{})
}
