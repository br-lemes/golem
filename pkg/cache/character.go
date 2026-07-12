package cache

import (
	"encoding/json"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

func CleanCharacter(name string) {
	cache.Where("name = ?", name).Delete(&models.Character{})
}

func GetCharacter(name string) *schemas.CharacterSchema {
	var characterCache models.Character
	if !isNameFresh(&characterCache, name, 10) {
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
	cache.Save(&models.Character{Name: name, Data: string(data)})
}
