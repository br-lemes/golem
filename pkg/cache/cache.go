package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var cache *gorm.DB

func Initialize(database config.Database) error {
	if database.Driver != "sqlite" {
		return fmt.Errorf("unsupported database driver: %s", database.Driver)
	}
	cacheFile := config.ExpandPath(database.Path)
	err := os.MkdirAll(filepath.Dir(cacheFile), 0755)
	if err != nil {
		return err
	}
	cache, err = gorm.Open(sqlite.Open(cacheFile), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return err
	}

	err = cache.Exec("PRAGMA journal_mode=WAL;").Error
	if err != nil {
		return err
	}
	err = cache.Exec("PRAGMA synchronous=NORMAL;").Error
	if err != nil {
		return err
	}
	err = cache.AutoMigrate(&models.APILog{}, &models.Cache{}, &models.Character{}, &models.FightSimulation{}, &models.UsageCombat{})
	if err != nil {
		return err
	}
	return nil
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
