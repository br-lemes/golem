package fight

import (
	"errors"
	"os"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestReadCharacterFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/character.json"
	err := os.WriteFile(path, []byte(`[{"level": 12}]`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	character, err := ReadCharacterFile(path)
	if err != nil || character.Level != 12 {
		t.Fatalf("character/error = %+v/%v", character, err)
	}
	for _, contents := range []string{"not json", "[]", "[{}, {}]"} {
		err = os.WriteFile(path, []byte(contents), 0600)
		if err != nil {
			t.Fatal(err)
		}
		_, err = ReadCharacterFile(path)
		if err == nil {
			t.Errorf("ReadCharacterFile(%q) unexpectedly succeeded", contents)
		}
	}
	_, err = ReadCharacterFile(dir + "/missing.json")
	if err == nil {
		t.Fatal("missing file unexpectedly succeeded")
	}
}

func TestCharacterSlotsAndUtilities(t *testing.T) {
	weapon, utility := "wooden_stick", "minor_health_potion"
	quantity := 2
	c := schemas.FakeCharacterSchema{
		WeaponSlot:           &weapon,
		Utility1Slot:         &utility,
		Utility1SlotQuantity: &quantity,
	}
	got := CharacterSlots(c)
	if got["weapon"] != weapon || got["utility1"] != utility {
		t.Fatalf("slots = %v", got)
	}
	utilities := CharacterUtilities(c)
	if utilities[utility] != quantity {
		t.Fatalf("utilities = %v", got)
	}
}

func TestResolveCharacterDefaultsAndImplicitUtilityQuantity(t *testing.T) {
	c, err := ResolveCharacter(CharacterOptions{
		Level: 7,
		ExplicitSlots: map[string]string{
			"weapon":   "wooden_stick",
			"utility1": "minor_health_potion",
		},
	})
	if err != nil || c.Level != 7 || c.WeaponSlot == nil || c.Utility1SlotQuantity == nil || *c.Utility1SlotQuantity != 1 {
		t.Fatalf("resolved character = %+v, err=%v", c, err)
	}
	_, err = ResolveCharacter(CharacterOptions{File: "a", Name: "b"})
	if err == nil {
		t.Fatal("expected mutually exclusive options error")
	}
}

func TestResolveCharacterByNameAndPropagatesDependencyError(t *testing.T) {
	want := errors.New("character unavailable")
	d := deps{
		characters: func(name string) (schemas.CharacterSchema, error) {
			if name == "missing" {
				return schemas.CharacterSchema{}, want
			}
			return schemas.CharacterSchema{Level: 9, WeaponSlot: "wooden_stick"}, nil
		},
	}
	character, err := resolveCharacter(d, CharacterOptions{Name: "Ada"})
	if err != nil || character.Level != 9 || character.WeaponSlot == nil {
		t.Fatalf("character/error = %+v/%v", character, err)
	}
	_, err = resolveCharacter(d, CharacterOptions{Name: "missing"})
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}

func TestResolveCharacterReadsFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/character.json"
	err := os.WriteFile(path, []byte(`[{"level": 12}]`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	character, err := ResolveCharacter(CharacterOptions{File: path})
	if err != nil || character.Level != 12 {
		t.Fatalf("character/error = %+v/%v", character, err)
	}
}

func TestResolveCharacterAppliesUtilityQuantities(t *testing.T) {
	quantities := map[string]int{"utility1": 4, "utility2": 5}
	character, err := ResolveCharacter(CharacterOptions{
		ExplicitSlots: map[string]string{
			"utility1": "minor_health_potion",
			"utility2": "small_antidote",
		},
		UtilityQuantities: quantities,
	})
	if err != nil {
		t.Fatalf("ResolveCharacter returned error: %v", err)
	}
	if character.Utility1SlotQuantity == nil || *character.Utility1SlotQuantity != 4 {
		t.Fatalf("utility 1 quantity = %+v, want 4", character.Utility1SlotQuantity)
	}
	if character.Utility2SlotQuantity == nil || *character.Utility2SlotQuantity != 5 {
		t.Fatalf("utility 2 quantity = %+v, want 5", character.Utility2SlotQuantity)
	}
}

func TestFakeCharacterFromCharacterLeavesEmptyUtilitiesUnset(t *testing.T) {
	character := FakeCharacterFromCharacter(schemas.CharacterSchema{Level: 3})
	if character.Utility1SlotQuantity != nil || character.Utility2SlotQuantity != nil {
		t.Fatalf("empty utilities got quantities: %+v", character)
	}
}

func TestFakeCharacterFromCharacterCopiesUtilityQuantities(t *testing.T) {
	character := FakeCharacterFromCharacter(schemas.CharacterSchema{
		Utility1Slot:         "minor_health_potion",
		Utility1SlotQuantity: 2,
		Utility2Slot:         "small_antidote",
		Utility2SlotQuantity: 3,
	})
	if character.Utility1Slot == nil || character.Utility1SlotQuantity == nil || *character.Utility1SlotQuantity != 2 {
		t.Fatalf("utility 1 was not copied: %+v", character)
	}
	if character.Utility2Slot == nil || character.Utility2SlotQuantity == nil || *character.Utility2SlotQuantity != 3 {
		t.Fatalf("utility 2 was not copied: %+v", character)
	}
}

func TestSetSlotAndQuantity(t *testing.T) {
	character := schemas.FakeCharacterSchema{}
	for _, slot := range []string{
		"weapon",
		"rune",
		"shield",
		"helmet",
		"body_armor",
		"leg_armor",
		"boots",
		"ring1",
		"ring2",
		"amulet",
		"artifact1",
		"artifact2",
		"artifact3",
		"utility1",
		"utility2",
	} {
		SetSlot(&character, slot, slot+"-item")
	}
	SetQuantity(&character, "utility1", 2)
	SetQuantity(&character, "utility2", 3)
	SetSlot(&character, "unknown", "ignored")
	SetQuantity(&character, "unknown", 4)
	if len(CharacterSlots(character)) != 15 || len(CharacterUtilities(character)) != 2 {
		t.Fatalf("slots/utilities = %v/%v", CharacterSlots(character), CharacterUtilities(character))
	}
	if *character.Utility1SlotQuantity != 2 || *character.Utility2SlotQuantity != 3 {
		t.Fatalf("quantities = %d/%d", *character.Utility1SlotQuantity, *character.Utility2SlotQuantity)
	}
}

func TestCharacterSlotsAndUtilitiesIgnoreNilPointers(t *testing.T) {
	if len(CharacterSlots(schemas.FakeCharacterSchema{})) != 0 || len(CharacterUtilities(schemas.FakeCharacterSchema{})) != 0 {
		t.Fatal("nil character pointers produced entries")
	}
}
