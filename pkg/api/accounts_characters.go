package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func AccountsCharacters(account string) ([]schemas.CharacterSchema, error) {
	resp, err := Get(fmt.Sprintf("/accounts/%s/characters", account), nil)
	if err != nil {
		return nil, err
	}
	var data schemas.CharactersListSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}
