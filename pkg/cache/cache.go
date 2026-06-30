package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var cache *gorm.DB

func APILog(method, path, body, response string, status, cooldown int) {
	cache.Create(&models.APILog{
		Method:   method,
		Path:     path,
		Body:     body,
		Response: response,
		Status:   status,
		Cooldown: cooldown,
	})
}

func CleanCharacter(name string) {
	cache.Where("name = ?", name).Delete(&models.Character{})
}

func GetCharacter(name string) *schemas.CharacterSchema {
	var logCache models.APILog
	result := cache.
		Where("created_at > datetime('now', '-5 minutes')").
		Limit(1).Find(&logCache)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	var characterCache models.Character
	result = cache.
		Where("name = ? AND updated_at > datetime('now', '-10 minutes')", name).
		Limit(1).Find(&characterCache)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	var character schemas.CharacterSchema
	err := json.Unmarshal([]byte(characterCache.Data), &character)
	if err != nil {
		return nil
	}

	return &character
}

func SaveCharacter(name string, character schemas.CharacterSchema) {
	data, err := json.Marshal(character)
	if err != nil {
		return
	}
	cache.Save(&models.Character{
		Name: name,
		Data: string(data),
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
	cache, err = gorm.Open(sqlite.Open(cacheFile), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		panic(err)
	}

	cache.AutoMigrate(&models.APILog{}, &models.Character{})
}
