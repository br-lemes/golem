package utils

import (
	"github.com/br-lemes/golem/pkg/cache"
)

var characters []string

func GetCharacters() []string {
	return characters
}

func init() {
	characters = cache.GetCharacters()
}
