package utils

import (
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

type EffectiveStats struct {
	Hp                                              int
	AttackFire, AttackEarth, AttackWater, AttackAir int
	Dmg                                             int
	DmgFire, DmgEarth, DmgWater, DmgAir             int
	ResFire, ResEarth, ResWater, ResAir             int
	CriticalStrike, Haste, Initiative               int
	Wisdom, Prospecting                             int
}

type EffectEntry struct {
	Code  string
	Value int
}

type Fighter struct {
	Stats   EffectiveStats
	Effects []EffectEntry
}

func isStatCode(code string) bool {
	switch code {
	case "hp", "attack_fire", "attack_earth", "attack_water", "attack_air",
		"dmg", "dmg_fire", "dmg_earth", "dmg_water", "dmg_air",
		"res_fire", "res_earth", "res_water", "res_air",
		"critical_strike", "haste", "initiative", "wisdom", "prospecting":
		return true
	}
	return false
}

func isConsumableEffect(code string) bool {
	switch code {
	case "heal", "gold", "teleport", "gems", "inventory_space", "threat":
		return true
	}
	return false
}

func isCombatEffect(code string) bool {
	return !isStatCode(code) && !isConsumableEffect(code)
}

func addToStats(s *EffectiveStats, code string, value int) {
	switch code {
	case "hp":
		s.Hp += value
	case "attack_fire":
		s.AttackFire += value
	case "attack_earth":
		s.AttackEarth += value
	case "attack_water":
		s.AttackWater += value
	case "attack_air":
		s.AttackAir += value
	case "dmg":
		s.Dmg += value
	case "dmg_fire":
		s.DmgFire += value
	case "dmg_earth":
		s.DmgEarth += value
	case "dmg_water":
		s.DmgWater += value
	case "dmg_air":
		s.DmgAir += value
	case "res_fire":
		s.ResFire += value
	case "res_earth":
		s.ResEarth += value
	case "res_water":
		s.ResWater += value
	case "res_air":
		s.ResAir += value
	case "critical_strike":
		s.CriticalStrike += value
	case "haste":
		s.Haste += value
	case "initiative":
		s.Initiative += value
	case "wisdom":
		s.Wisdom += value
	case "prospecting":
		s.Prospecting += value
	}
}

func foldEffects(acc *EffectiveStats, effects *[]schemas.SimpleEffectSchema, combat *[]EffectEntry, sign int) {
	if effects == nil {
		return
	}
	for _, e := range *effects {
		if isStatCode(e.Code) {
			addToStats(acc, e.Code, sign*e.Value)
		} else if isCombatEffect(e.Code) {
			*combat = append(*combat, EffectEntry{Code: e.Code, Value: sign * e.Value})
		}
	}
}

func equippedSlotMap(c schemas.CharacterSchema) map[string]string {
	return map[string]string{
		"amulet": c.AmuletSlot, "artifact1": c.Artifact1Slot, "artifact2": c.Artifact2Slot,
		"artifact3": c.Artifact3Slot, "bag": c.BagSlot, "body_armor": c.BodyArmorSlot,
		"boots": c.BootsSlot, "helmet": c.HelmetSlot, "leg_armor": c.LegArmorSlot,
		"ring1": c.Ring1Slot, "ring2": c.Ring2Slot, "rune": c.RuneSlot, "shield": c.ShieldSlot,
		"utility1": c.Utility1Slot, "utility2": c.Utility2Slot, "weapon": c.WeaponSlot,
	}
}

func currentStats(c schemas.CharacterSchema) EffectiveStats {
	return EffectiveStats{
		Hp:             c.MaxHp,
		AttackFire:     c.AttackFire,
		AttackEarth:    c.AttackEarth,
		AttackWater:    c.AttackWater,
		AttackAir:      c.AttackAir,
		Dmg:            c.Dmg,
		DmgFire:        c.DmgFire,
		DmgEarth:       c.DmgEarth,
		DmgWater:       c.DmgWater,
		DmgAir:         c.DmgAir,
		ResFire:        c.ResFire,
		ResEarth:       c.ResEarth,
		ResWater:       c.ResWater,
		ResAir:         c.ResAir,
		CriticalStrike: c.CriticalStrike,
		Haste:          c.Haste,
		Initiative:     c.Initiative,
		Wisdom:         c.Wisdom,
		Prospecting:    c.Prospecting,
	}
}

func baseStats(c schemas.CharacterSchema) EffectiveStats {
	base := currentStats(c)
	equipped := EffectiveStats{}
	var junk []EffectEntry
	for _, code := range equippedSlotMap(c) {
		if code == "" {
			continue
		}
		it, ok := database.GetItem(code)
		if !ok {
			continue
		}
		foldEffects(&equipped, it.Effects, &junk, 1)
	}
	base.Hp -= equipped.Hp
	base.AttackFire -= equipped.AttackFire
	base.AttackEarth -= equipped.AttackEarth
	base.AttackWater -= equipped.AttackWater
	base.AttackAir -= equipped.AttackAir
	base.Dmg -= equipped.Dmg
	base.DmgFire -= equipped.DmgFire
	base.DmgEarth -= equipped.DmgEarth
	base.DmgWater -= equipped.DmgWater
	base.DmgAir -= equipped.DmgAir
	base.ResFire -= equipped.ResFire
	base.ResEarth -= equipped.ResEarth
	base.ResWater -= equipped.ResWater
	base.ResAir -= equipped.ResAir
	base.CriticalStrike -= equipped.CriticalStrike
	base.Haste -= equipped.Haste
	base.Initiative -= equipped.Initiative
	base.Wisdom -= equipped.Wisdom
	base.Prospecting -= equipped.Prospecting
	return base
}

func applyGear(base EffectiveStats, codes []string) Fighter {
	stats := base
	var effects []EffectEntry
	for _, code := range codes {
		if code == "" {
			continue
		}
		it, ok := database.GetItem(code)
		if !ok {
			continue
		}
		foldEffects(&stats, it.Effects, &effects, 1)
	}
	return Fighter{Stats: stats, Effects: effects}
}

func effMap(entries []EffectEntry) map[string]int {
	m := make(map[string]int)
	for _, e := range entries {
		m[e.Code] += e.Value
	}
	return m
}

func monsterEffectEntries(monster schemas.MonsterSchema) []EffectEntry {
	if monster.Effects == nil {
		return nil
	}
	out := make([]EffectEntry, 0, len(*monster.Effects))
	for _, e := range *monster.Effects {
		out = append(out, EffectEntry{Code: e.Code, Value: e.Value})
	}
	return out
}
