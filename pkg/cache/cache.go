package cache

import (
	"os"
	"path/filepath"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var cache *gorm.DB

func APILog(method, path, body, response string, status, cooldown int) {
	cache.Create(&models.APILog{
		Method:   method,
		Path:     path,
		Body:     string(body),
		Response: string(response),
		Status:   status,
		Cooldown: cooldown,
	})
}

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
	cache, err = gorm.Open(sqlite.Open(cacheFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	cache.AutoMigrate(&models.APILog{})
}
