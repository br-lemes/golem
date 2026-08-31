package cache

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
	"gorm.io/gorm"
)

func CleanCharacter(name string) {
	cache.Where("name = ?", name).Delete(&models.Character{})
}

func CleanCharacters() {
	cache.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Character{})
}

func GetCharacter(name string) *schemas.CharacterSchema {
	var characterCache models.Character
	if !findByName(&characterCache, name) {
		return nil
	}
	var character schemas.CharacterSchema
	err := json.Unmarshal([]byte(characterCache.Data), &character)
	if err != nil {
		return nil
	}

	return &character
}

func GetCharacters() []string {
	var names []string
	cache.Model(&models.Character{}).Unscoped().Pluck("name", &names)
	return names
}

func SaveCharacter(name string, character schemas.CharacterSchema) {
	SaveCharacters([]schemas.CharacterSchema{character})
}

func SaveCharacters(characters []schemas.CharacterSchema) {
	account := GetAccount()
	for _, character := range characters {
		if character.Account != account {
			continue
		}
		saveCharacter(character.Name, character)
	}
}

func saveCharacter(name string, character schemas.CharacterSchema) {
	if character.Account == "" {
		return
	}
	data, err := json.Marshal(character)
	if err != nil {
		return
	}
	cache.Save(&models.Character{Name: name, Data: string(data)})
}
