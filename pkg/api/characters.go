package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Characters(name string) (schemas.CharacterSchema, error) {
	if name == "." {
		name = os.Getenv("GOLEM_NAME")
		if name == "" {
			return schemas.CharacterSchema{}, fmt.Errorf("GOLEM_NAME is not set")
		}
	}
	character := cache.GetCharacter(name)
	if character != nil {
		return *character, nil
	}
	resp, err := Get(fmt.Sprintf("/characters/%s", name), nil)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	var data schemas.CharacterResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterSchema{}, err
	}
	cache.SaveCharacter(data.Data)
	return data.Data, nil
}
