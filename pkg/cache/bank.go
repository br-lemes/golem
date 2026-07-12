package cache

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
	"gorm.io/gorm"
)

func CleanBank() {
	cache.Where("name = ?", "bank").Delete(&models.Cache{})
}

func GetBank() *schemas.BankSchema {
	var bankCache models.Cache
	if !isNameFresh(&bankCache, "bank", 30) {
		return nil
	}
	var bank schemas.BankSchema
	err := json.Unmarshal([]byte(bankCache.Data), &bank)
	if err != nil {
		return nil
	}
	return &bank
}

func SaveBank(bank schemas.BankSchema) {
	data, err := json.Marshal(bank)
	if err != nil {
		return
	}
	cache.Save(&models.Cache{Name: "bank", Data: string(data)})
}

func UpdateBankGold(gold int) {
	cache.Model(&models.Cache{}).Where("name = ?", "bank").Updates(
		map[string]any{"data": gorm.Expr("json_set(data, '$.gold', ?)", gold)})
}

func CleanBankItems() {
	cache.Where("name = ?", "bankItems").Delete(&models.Cache{})
}

func GetBankItems() *[]schemas.SimpleItemSchema {
	var bankItemsCache models.Cache
	if !isNameFresh(&bankItemsCache, "bankItems", 30) {
		return nil
	}
	var bankItems []schemas.SimpleItemSchema
	err := json.Unmarshal([]byte(bankItemsCache.Data), &bankItems)
	if err != nil {
		return nil
	}
	return &bankItems
}

func SaveBankItems(bankItems []schemas.SimpleItemSchema) {
	data, err := json.Marshal(bankItems)
	if err != nil {
		return
	}
	cache.Save(&models.Cache{Name: "bankItems", Data: string(data)})
}
