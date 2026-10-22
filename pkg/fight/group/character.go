package group

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

func FightersFromFile(path string) ([]Fighter, error) {
	characters, err := CharactersFromFile(path)
	if err != nil {
		return nil, err
	}
	result := make([]Fighter, len(characters))
	for i, character := range characters {
		result[i] = FighterFromFakeCharacter(i, character)
	}
	return result, nil
}

func CharactersFromFile(path string) ([]schemas.FakeCharacterSchema, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read group simulation file: %w", err)
	}
	var characters []schemas.FakeCharacterSchema
	err = json.Unmarshal(body, &characters)
	if err != nil {
		return nil, fmt.Errorf("parse group simulation file: %w", err)
	}
	if len(characters) < 1 || len(characters) > 3 {
		return nil, fmt.Errorf("group simulation file must contain 1 to 3 characters")
	}
	return characters, nil
}

func FighterFromFakeCharacter(index int, character schemas.FakeCharacterSchema) Fighter {
	level := character.Level
	if level < 1 {
		level = 1
	}
	name := fmt.Sprintf("Character_%d", index+1)
	items := []string{}
	slots := []*string{
		character.WeaponSlot,
		character.RuneSlot,
		character.ShieldSlot,
		character.HelmetSlot,
		character.BodyArmorSlot,
		character.LegArmorSlot,
		character.BootsSlot,
		character.Ring1Slot,
		character.Ring2Slot,
		character.AmuletSlot,
		character.Artifact1Slot,
		character.Artifact2Slot,
		character.Artifact3Slot,
	}
	for _, slot := range slots {
		if slot != nil {
			items = append(items, *slot)
		}
	}
	fighter := FighterFromLoadout(name, level, items)
	applyUtilityBoosts(&fighter.Stats, character.Utility1Slot)
	applyUtilityBoosts(&fighter.Stats, character.Utility2Slot)
	fighter.Utilities = append(fighter.Utilities, utilitiesFromSlot(character.Utility1Slot, character.Utility1SlotQuantity)...)
	fighter.Utilities = append(fighter.Utilities, utilitiesFromSlot(character.Utility2Slot, character.Utility2SlotQuantity)...)
	return fighter
}

func applyUtilityBoosts(stats *Stats, code *string) {
	if code == nil {
		return
	}
	item, found := catalog.Items().Get(*code)
	if !found || item.Effects == nil {
		return
	}
	for _, effect := range *item.Effects {
		switch effect.Code {
		case "boost_hp":
			stats.HP += effect.Value
		case "boost_dmg_fire":
			stats.DmgFire += effect.Value
		case "boost_dmg_earth":
			stats.DmgEarth += effect.Value
		case "boost_dmg_water":
			stats.DmgWater += effect.Value
		case "boost_dmg_air":
			stats.DmgAir += effect.Value
		case "boost_res_fire":
			stats.ResFire += effect.Value
		case "boost_res_earth":
			stats.ResEarth += effect.Value
		case "boost_res_water":
			stats.ResWater += effect.Value
		case "boost_res_air":
			stats.ResAir += effect.Value
		}
	}
}

func FighterFromLoadout(name string, level int, items []string) Fighter {
	if level < 1 {
		level = 1
	}
	fighter := Fighter{Name: name}
	fighter.Stats.HP = 120 + (level-1)*5
	fighter.Stats.Initiative = 100
	for _, code := range items {
		item, found := catalog.Items().Get(code)
		if !found || item.Effects == nil {
			continue
		}
		for _, effect := range *item.Effects {
			applyEffect(&fighter.Stats, effect.Code, effect.Value)
		}
	}
	return fighter
}

func utilitiesFromSlot(code *string, quantity *int) []Utility {
	if code == nil || quantity == nil || *quantity <= 0 {
		return nil
	}
	item, found := catalog.Items().Get(*code)
	if !found || item.Effects == nil {
		return nil
	}
	utility := Utility{Code: *code, Quantity: *quantity}
	for _, effect := range *item.Effects {
		switch effect.Code {
		case "restore":
			utility.Restore += effect.Value
		case "splash_restore":
			utility.SplashRestore += effect.Value
		case "antipoison":
			utility.Antipoison += effect.Value
		}
	}
	if utility.Restore == 0 && utility.SplashRestore == 0 && utility.Antipoison == 0 {
		return nil
	}
	return []Utility{utility}
}

func applyEffect(stats *Stats, code string, value int) {
	switch code {
	case "hp":
		stats.HP += value
	case "attack_fire":
		stats.AttackFire += value
	case "attack_earth":
		stats.AttackEarth += value
	case "attack_water":
		stats.AttackWater += value
	case "attack_air":
		stats.AttackAir += value
	case "dmg":
		stats.Dmg += value
	case "dmg_fire":
		stats.DmgFire += value
	case "dmg_earth":
		stats.DmgEarth += value
	case "dmg_water":
		stats.DmgWater += value
	case "dmg_air":
		stats.DmgAir += value
	case "res_fire":
		stats.ResFire += value
	case "res_earth":
		stats.ResEarth += value
	case "res_water":
		stats.ResWater += value
	case "res_air":
		stats.ResAir += value
	case "critical_strike":
		stats.CriticalStrike += value
	case "initiative":
		stats.Initiative += value
	case "threat":
		stats.Threat += value
	case "shell":
		stats.Shell += value
	case "healing":
		stats.Healing += value
	case "healing_aura":
		stats.HealingAura += value
	case "lifesteal":
		stats.Lifesteal += value
	case "burn":
		stats.Burn += value
	case "frenzy":
		stats.Frenzy += value
	case "greed":
		stats.Greed += value
	case "vampiric_strike":
		stats.VampiricStrike += value
	case "guard":
		stats.Guard += value
	case "enchanted_mirror":
		stats.EnchantedMirror += value
	}
}

func FighterFromCharacter(character schemas.CharacterSchema) Fighter {
	return Fighter{
		Name: character.Name,
		Stats: Stats{
			HP:             character.MaxHp,
			Initiative:     character.Initiative,
			Threat:         character.Threat,
			CriticalStrike: character.CriticalStrike,
			AttackFire:     character.AttackFire,
			AttackEarth:    character.AttackEarth,
			AttackWater:    character.AttackWater,
			AttackAir:      character.AttackAir,
			Dmg:            character.Dmg,
			DmgFire:        character.DmgFire,
			DmgEarth:       character.DmgEarth,
			DmgWater:       character.DmgWater,
			DmgAir:         character.DmgAir,
			ResFire:        character.ResFire,
			ResEarth:       character.ResEarth,
			ResWater:       character.ResWater,
			ResAir:         character.ResAir,
			Haste:          character.Haste,
			Wisdom:         character.Wisdom,
		},
	}
}

func MonsterFromSchema(monster schemas.MonsterSchema) Monster {
	result := Monster{
		Name:          monster.Name,
		Type:          monster.Type,
		BubbleElement: -1,
		BaseRes: [4]int{
			monster.ResFire,
			monster.ResEarth,
			monster.ResWater,
			monster.ResAir,
		},
		HP:    float64(monster.Hp),
		MaxHP: float64(monster.Hp),
		Stats: Stats{
			HP:             monster.Hp,
			Initiative:     monster.Initiative,
			AttackFire:     monster.AttackFire,
			AttackEarth:    monster.AttackEarth,
			AttackWater:    monster.AttackWater,
			AttackAir:      monster.AttackAir,
			CriticalStrike: monster.CriticalStrike,
			ResFire:        monster.ResFire,
			ResEarth:       monster.ResEarth,
			ResWater:       monster.ResWater,
			ResAir:         monster.ResAir,
		},
	}
	if monster.Effects != nil {
		for _, effect := range *monster.Effects {
			if effect.Code == "barrier" {
				result.Barrier = float64(effect.Value)
				result.BarrierAmount = float64(effect.Value)
			}
			if effect.Code == "lifesteal" {
				result.Stats.Lifesteal = effect.Value
			}
			if effect.Code == "healing" {
				result.Stats.Healing = effect.Value
			}
			if effect.Code == "poison" {
				result.Stats.Poison = effect.Value
			}
			if effect.Code == "burn" {
				result.Stats.Burn = effect.Value
			}
			if effect.Code == "corrupted" {
				result.Stats.Corrupted = effect.Value
			}
			if effect.Code == "reconstitution" {
				result.Reconstitution = effect.Value
			}
			if effect.Code == "void_drain" {
				result.VoidDrain = effect.Value
			}
			if effect.Code == "frenzy" {
				result.Stats.Frenzy = effect.Value
			}
			if effect.Code == "greed" {
				result.Stats.Greed = effect.Value
			}
			if effect.Code == "berserker_rage" {
				result.BerserkerRage = effect.Value
			}
			if effect.Code == "protective_bubble" {
				result.Stats.ProtectiveBubble = effect.Value
			}
			if effect.Code == "sun_shield" {
				result.Stats.SunShield = effect.Value
			}
		}
	}
	return result
}
