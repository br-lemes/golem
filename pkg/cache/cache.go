package cache

import (
	"os"
	"path/filepath"
	"time"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var cache *gorm.DB

func init() {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	cacheDir := filepath.Join(userCacheDir, "golem")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		panic(err)
	}
	cacheFile := filepath.Join(cacheDir, "season8.db")
	cache, err = gorm.Open(sqlite.Open(cacheFile), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		panic(err)
	}

	cache.Exec("PRAGMA journal_mode=WAL;")
	cache.Exec("PRAGMA synchronous=NORMAL;")
	cache.AutoMigrate(&models.APILog{}, &models.Cache{}, &models.Character{})
}

func isNameFresh(dest any, name string, minutes int) bool {
	if !isStatusFresh() {
		return false
	}

	result := cache.Where("name = ? AND updated_at >= ?", name, limit(minutes)).
		Limit(1).Find(dest)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func isTableFresh(dest any, minutes int) bool {
	if !isStatusFresh() {
		return false
	}

	result := cache.Where("updated_at >= ?", limit(minutes)).Find(dest)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func isStatusFresh() bool {
	var statusCache models.Cache
	result := cache.Where("name = ? AND updated_at >= ?", "status", limit(5)).
		Limit(1).Find(&statusCache)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func limit(minutes int) time.Time {
	return time.Now().UTC().Add(time.Duration(-minutes) * time.Minute)
}
