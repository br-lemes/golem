package cache

import "github.com/br-lemes/golem/pkg/models"

func APILog(method, path, body, response string, status, cooldown int) {
	cache.Save(&models.Cache{Name: "status"})
	cache.Create(&models.APILog{
		Method:   method,
		Path:     path,
		Body:     body,
		Response: response,
		Status:   status,
		Cooldown: cooldown,
	})
}
