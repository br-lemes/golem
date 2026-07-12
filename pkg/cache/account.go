package cache

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

func CleanAccount() {
	cache.Where("name = ?", "account").Delete(&models.Cache{})
}

func GetAccount() string {
	var accountCache models.Cache
	cache.Where("name = ?", "account").Find(&accountCache)
	return accountCache.Data
}

func SaveAccount(account string) {
	cache.Save(&models.Cache{Name: "account", Data: account})
}

func GetAccountCharacters() []schemas.CharacterSchema {
	var charactersCache []models.Character
	if !isTableFresh(&charactersCache, 10) || len(charactersCache) != 5 {
		return nil
	}
	characters := make([]schemas.CharacterSchema, 5)
	for i, characterCache := range charactersCache {
		err := json.Unmarshal([]byte(characterCache.Data), &characters[i])
		if err != nil {
			return nil
		}
	}
	return characters
}
