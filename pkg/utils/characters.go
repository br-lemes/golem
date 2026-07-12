package utils

import "github.com/br-lemes/golem/pkg/api"

var characters []string

func GetCharacters() []string {
	return characters
}

func init() {
	chars, err := api.AccountsCharacters("")
	if err != nil {
		panic(err)
	}
	characters = make([]string, len(chars))
	for i, char := range chars {
		characters[i] = char.Name
	}
}
