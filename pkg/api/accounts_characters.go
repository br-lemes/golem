package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func AccountsCharacters(account string) ([]schemas.CharacterSchema, error) {
	myAccount := cache.GetAccount()
	if myAccount == "" {
		_, err := MyDetails()
		if err != nil {
			return nil, err
		}
		myAccount = cache.GetAccount()
	}
	if account == "" {
		account = myAccount
	}
	if account == myAccount {
		characters := cache.GetAccountCharacters()
		if characters != nil {
			return characters, nil
		}
	}
	resp, err := Get(fmt.Sprintf("/accounts/%s/characters", account), nil)
	if err != nil {
		return nil, err
	}
	var data schemas.CharactersListSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	for _, character := range data.Data {
		cache.SaveCharacter(character.Name, character)
	}
	return data.Data, nil
}
