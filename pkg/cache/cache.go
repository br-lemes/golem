package cache

import (
	"os"
	"path/filepath"
	"time"

	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var cache *gorm.DB
var simulationDB *gorm.DB

func Initialize(storage config.Storage) error {
	cacheFile := config.ExpandPath(storage.Cache)
	simulationFile := config.ExpandPath(storage.Simulation)
	if storage.Simulation == "" {
		simulationFile = cacheFile
	}
	var err error
	cache, err = openDatabase(cacheFile)
	if err != nil {
		//+gocover:ignore:block database setup failure is environmental
		return err
	}

	err = cache.AutoMigrate(&models.Cache{}, &models.Character{}, &models.OutputFilter{})
	if err != nil {
		//+gocover:ignore:block schema migration failure is environmental
		return err
	}
	simulationDB, err = openDatabase(simulationFile)
	if err != nil {
		//+gocover:ignore:block database setup failure is environmental
		return err
	}
	err = simulationDB.AutoMigrate(&models.FightSimulation{}, &models.UsageCombat{})
	if err != nil {
		//+gocover:ignore:block schema migration failure is environmental
		return err
	}
	return nil
}

func openDatabase(path string) (*gorm.DB, error) {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		//+gocover:ignore:block filesystem setup failure is environmental
		return nil, err
	}
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		//+gocover:ignore:block SQLite setup failure is environmental
		return nil, err
	}
	err = database.Exec("PRAGMA journal_mode=WAL;").Error
	if err != nil {
		//+gocover:ignore:block SQLite pragma failure is environmental
		return nil, err
	}
	err = database.Exec("PRAGMA synchronous=NORMAL;").Error
	if err != nil {
		//+gocover:ignore:block SQLite pragma failure is environmental
		return nil, err
	}
	return database, nil
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
