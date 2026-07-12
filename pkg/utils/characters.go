package utils

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
)

var characters []string

func GetCharacters() []string {
	return characters
}

func init() {
	characters = cache.GetCharacters()
	if len(characters) == 5 {
		return
	}
	chars, err := api.AccountsCharacters("")
	if err != nil {
		panic(err)
	}
	characters = make([]string, len(chars))
	for i, char := range chars {
		characters[i] = char.Name
	}
}
