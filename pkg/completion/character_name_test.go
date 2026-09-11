package completion

import (
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCharacterNameCompletion(t *testing.T) {
	err := cache.Initialize(config.Storage{Cache: t.TempDir() + "/cache.db"})
	if err != nil {
		t.Fatal(err)
	}
	cache.SaveAccount("account-name")
	character := schemas.CharacterSchema{}
	character.Name = "hero"
	character.Account = "account-name"
	cache.SaveCharacter(character)
	complete := CharacterName(1).Build()
	got, _ := complete(nil, nil, "")
	if len(got) != 1 || got[0] != "hero" {
		t.Fatalf("character completion = %#v", got)
	}
}
