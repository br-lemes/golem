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

	err = cache.Exec("PRAGMA journal_mode=WAL;").Error
	if err != nil {
		panic(err)
	}
	err = cache.Exec("PRAGMA synchronous=NORMAL;").Error
	if err != nil {
		panic(err)
	}
	err = cache.AutoMigrate(&models.APILog{}, &models.Cache{}, &models.Character{}, &models.FightSimulation{}, &models.UsageCombat{})
	if err != nil {
		panic(err)
	}
}

func findByName(dest any, name string) bool {
	result := cache.Where("name = ?", name).Limit(1).Find(dest)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}

func findTable(dest any) bool {
	result := cache.Find(dest)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}
	return true
}
