package cache

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestCharacterCache(t *testing.T) {
	initializeTestCache(t)
	SaveAccount("account-name")
	character := schemas.CharacterSchema{}
	character.Name = "hero"
	character.Account = "account-name"
	character.Level = 10
	SaveCharacter(character)
	got := GetCharacter("hero")
	if got == nil || got.Name != "hero" || got.Level != 10 {
		t.Fatalf("GetCharacter() = %#v", got)
	}
	characters := GetCharacters()
	if len(characters) != 1 || characters[0] != "hero" {
		t.Fatalf("GetCharacters() = %#v", characters)
	}
	CleanCharacter("hero")
	if GetCharacter("hero") != nil {
		t.Fatal("character was not cleaned")
	}
}

func TestGetAccountCharacters(t *testing.T) {
	initializeTestCache(t)
	if GetAccountCharacters() != nil {
		t.Fatal("empty character cache returned characters")
	}
	SaveAccount("account-name")
	character := schemas.CharacterSchema{}
	character.Name = "hero"
	character.Account = "account-name"
	character.Level = 10
	SaveCharacter(character)
	characters := GetAccountCharacters()
	if len(characters) != 5 || characters[0].Name != "hero" {
		t.Fatalf("GetAccountCharacters() = %#v", characters)
	}
}

func TestGetAccountCharactersRejectsCorruptData(t *testing.T) {
	initializeTestCache(t)
	cache.Create(&models.Character{Name: "hero", Data: "not-json"})
	if GetAccountCharacters() != nil {
		t.Fatal("corrupt character data was accepted")
	}
}

func TestCleanCharactersRemovesAllCharacters(t *testing.T) {
	initializeTestCache(t)
	SaveAccount("account-name")
	hero := schemas.CharacterSchema{Name: "hero", Account: "account-name"}
	other := schemas.CharacterSchema{Name: "other", Account: "account-name"}
	SaveCharacter(hero)
	SaveCharacter(other)
	CleanCharacters()
	if GetCharacter("hero") != nil || GetCharacter("other") != nil {
		t.Fatal("CleanCharacters left active characters")
	}
}

func TestSaveCharactersSkipsOtherAccounts(t *testing.T) {
	initializeTestCache(t)
	SaveAccount("account-name")
	characters := []schemas.CharacterSchema{
		{Name: "hero", Account: "account-name"},
		{Name: "other", Account: "other-account"},
	}
	SaveCharacters(characters)
	if GetCharacter("hero") == nil || GetCharacter("other") != nil {
		t.Fatal("SaveCharacters cached the wrong account")
	}
}

func TestSaveCharacterSkipsMissingAccount(t *testing.T) {
	initializeTestCache(t)
	SaveCharacter(schemas.CharacterSchema{Name: "hero"})
	if GetCharacter("hero") != nil {
		t.Fatal("character without account was cached")
	}
}

func TestGetCharacterRejectsCorruptData(t *testing.T) {
	initializeTestCache(t)
	cache.Create(&models.Character{Name: "hero", Data: "not-json"})
	if GetCharacter("hero") != nil {
		t.Fatal("corrupt character data was accepted")
	}
}
