package cache

import "github.com/br-lemes/golem/pkg/models"

func ListStocks() []models.Stock {
	var stocks []models.Stock
	if cache.Order("code").Find(&stocks).Error != nil {
		//+gocover:ignore:block SQLite query failure is environmental
		return nil
	}
	return stocks
}

func GetStock(code string) (models.Stock, bool) {
	var stock models.Stock
	result := cache.Where("code = ?", code).Limit(1).Find(&stock)
	return stock, result.Error == nil && result.RowsAffected > 0
}

func SetStock(code string, quantity int, mode models.StockMode) error {
	var stock models.Stock
	result := cache.Where("code = ?", code).Limit(1).Find(&stock)
	if result.Error != nil {
		return result.Error
	}
	stock.Code = code
	stock.Quantity = quantity
	stock.Mode = mode
	stock.Enabled = true
	if mode != models.StockModeExact || quantity != 0 {
		stock.Discard = false
	}
	return cache.Save(&stock).Error
}

func SetStockEnabled(code string, enabled bool) error {
	result := cache.Model(&models.Stock{}).Where("code = ?", code).Update("enabled", enabled)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return nil
	}
	return cache.Create(&models.Stock{Code: code, Enabled: enabled}).Error
}

func SetStockDiscard(code string, discard bool) error {
	return cache.Model(&models.Stock{}).Where("code = ?", code).Update("discard", discard).Error
}

func DiscardStockCodes() []string {
	var codes []string
	result := cache.Model(&models.Stock{}).Where("discard = ? AND enabled = ? AND mode = ? AND quantity = ?", true, true, models.StockModeExact, 0).Order("code").Pluck("code", &codes)
	if result.Error != nil {
		//+gocover:ignore:block SQLite query failure is environmental
		return nil
	}
	return codes
}

func RemoveStock(code string) error {
	return cache.Where("code = ?", code).Delete(&models.Stock{}).Error
}
