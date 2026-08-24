package simulation

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

type CharacterOptions struct {
	File              string
	Name              string
	Level             int
	ExplicitSlots     map[string]string
	UtilityQuantities map[string]int
}

func ReadCharacterFile(path string) (schemas.FakeCharacterSchema, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return schemas.FakeCharacterSchema{}, fmt.Errorf("read simulation file: %w", err)
	}
	var characters []schemas.FakeCharacterSchema
	err = json.Unmarshal(body, &characters)
	if err != nil {
		return schemas.FakeCharacterSchema{}, fmt.Errorf("parse simulation file: %w", err)
	}
	if len(characters) != 1 {
		return schemas.FakeCharacterSchema{}, fmt.Errorf("characters must contain exactly one entry")
	}
	return characters[0], nil
}

func ResolveCharacter(options CharacterOptions) (schemas.FakeCharacterSchema, error) {
	if options.File != "" && options.Name != "" {
		return schemas.FakeCharacterSchema{}, fmt.Errorf("--name and --file are mutually exclusive")
	}
	var character schemas.FakeCharacterSchema
	var err error
	switch {
	case options.File != "":
		character, err = ReadCharacterFile(options.File)
	case options.Name != "":
		var source schemas.CharacterSchema
		source, err = api.Characters(options.Name)
		if err == nil {
			character = FakeCharacterFromCharacter(source)
		}
	default:
		character.Level = options.Level
		character.WeaponSlot = stringPointer("wooden_stick")
	}
	if err != nil {
		return schemas.FakeCharacterSchema{}, err
	}
	if options.Level > 0 && options.File == "" && options.Name == "" {
		character.Level = options.Level
	}
	for slot, value := range options.ExplicitSlots {
		SetSlot(&character, slot, value)
	}
	for slot, quantity := range options.UtilityQuantities {
		SetQuantity(&character, slot, quantity)
	}
	if character.Utility1Slot != nil && character.Utility1SlotQuantity == nil {
		SetQuantity(&character, "utility1", 1)
	}
	if character.Utility2Slot != nil && character.Utility2SlotQuantity == nil {
		SetQuantity(&character, "utility2", 1)
	}
	return character, nil
}

func FakeCharacterFromCharacter(c schemas.CharacterSchema) schemas.FakeCharacterSchema {
	result := schemas.FakeCharacterSchema{
		Level:         c.Level,
		WeaponSlot:    stringPointer(c.WeaponSlot),
		RuneSlot:      stringPointer(c.RuneSlot),
		ShieldSlot:    stringPointer(c.ShieldSlot),
		HelmetSlot:    stringPointer(c.HelmetSlot),
		BodyArmorSlot: stringPointer(c.BodyArmorSlot),
		LegArmorSlot:  stringPointer(c.LegArmorSlot),
		BootsSlot:     stringPointer(c.BootsSlot),
		Ring1Slot:     stringPointer(c.Ring1Slot),
		Ring2Slot:     stringPointer(c.Ring2Slot),
		AmuletSlot:    stringPointer(c.AmuletSlot),
		Artifact1Slot: stringPointer(c.Artifact1Slot),
		Artifact2Slot: stringPointer(c.Artifact2Slot),
		Artifact3Slot: stringPointer(c.Artifact3Slot),
		Utility1Slot:  stringPointer(c.Utility1Slot),
		Utility2Slot:  stringPointer(c.Utility2Slot),
	}
	if c.Utility1Slot != "" {
		result.Utility1SlotQuantity = intPointer(c.Utility1SlotQuantity)
	}
	if c.Utility2Slot != "" {
		result.Utility2SlotQuantity = intPointer(c.Utility2SlotQuantity)
	}
	return result
}

func CharacterSlots(c schemas.FakeCharacterSchema) map[string]string {
	slots := map[string]string{}
	for slot, value := range map[string]*string{
		"weapon":     c.WeaponSlot,
		"rune":       c.RuneSlot,
		"shield":     c.ShieldSlot,
		"helmet":     c.HelmetSlot,
		"body_armor": c.BodyArmorSlot,
		"leg_armor":  c.LegArmorSlot,
		"boots":      c.BootsSlot,
		"ring1":      c.Ring1Slot,
		"ring2":      c.Ring2Slot,
		"amulet":     c.AmuletSlot,
		"artifact1":  c.Artifact1Slot,
		"artifact2":  c.Artifact2Slot,
		"artifact3":  c.Artifact3Slot,
		"utility1":   c.Utility1Slot,
		"utility2":   c.Utility2Slot,
	} {
		if value != nil {
			slots[slot] = *value
		}
	}
	return slots
}

func CharacterUtilities(c schemas.FakeCharacterSchema) map[string]int {
	utilities := map[string]int{}
	if c.Utility1Slot != nil && c.Utility1SlotQuantity != nil {
		utilities[*c.Utility1Slot] = *c.Utility1SlotQuantity
	}
	if c.Utility2Slot != nil && c.Utility2SlotQuantity != nil {
		utilities[*c.Utility2Slot] = *c.Utility2SlotQuantity
	}
	return utilities
}

func SetSlot(c *schemas.FakeCharacterSchema, slot, value string) {
	p := stringPointer(value)
	switch slot {
	case "weapon":
		c.WeaponSlot = p
	case "rune":
		c.RuneSlot = p
	case "shield":
		c.ShieldSlot = p
	case "helmet":
		c.HelmetSlot = p
	case "body_armor":
		c.BodyArmorSlot = p
	case "leg_armor":
		c.LegArmorSlot = p
	case "boots":
		c.BootsSlot = p
	case "ring1":
		c.Ring1Slot = p
	case "ring2":
		c.Ring2Slot = p
	case "amulet":
		c.AmuletSlot = p
	case "artifact1":
		c.Artifact1Slot = p
	case "artifact2":
		c.Artifact2Slot = p
	case "artifact3":
		c.Artifact3Slot = p
	case "utility1":
		c.Utility1Slot = p
	case "utility2":
		c.Utility2Slot = p
	}
}

func SetQuantity(c *schemas.FakeCharacterSchema, slot string, quantity int) {
	p := intPointer(quantity)
	if slot == "utility1" {
		c.Utility1SlotQuantity = p
	}
	if slot == "utility2" {
		c.Utility2SlotQuantity = p
	}
}

func stringPointer(value string) *string { return &value }
func intPointer(value int) *int          { return &value }
