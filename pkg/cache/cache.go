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

	cache.AutoMigrate(&models.APILog{}, &models.Cache{}, &models.Character{})
}

func isFresh(dest any, name string, minutes int) bool {
	var logCache models.APILog
	result := cache.Where("created_at >= ?", limit(5)).Limit(1).Find(&logCache)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}

	result = cache.Where("name = ? AND updated_at >= ?", name, limit(minutes)).
		Limit(1).Find(dest)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func limit(minutes int) time.Time {
	return time.Now().UTC().Add(time.Duration(-minutes) * time.Minute)
}
