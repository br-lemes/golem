package damage

import (
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestEquipmentBonusesRepeatsRings(t *testing.T) {
	pools := []bonusPool{{
		rule: slotRule{kind: "ring", count: 2},
		items: []bonusItem{
			{bonus: bonus{bonusCritical: 7}},
			{bonus: bonus{bonusCritical: 7}},
		},
	}}
	got := equipmentBonuses(pools)
	if got[bonusCritical] != 14 {
		t.Fatalf("ring critical = %d, want 14", got[bonusCritical])
	}
}

func TestEquipmentBonusesUsesDistinctArtifacts(t *testing.T) {
	pools := []bonusPool{{
		rule: slotRule{kind: "artifact", count: 3, distinct: true},
		items: []bonusItem{
			{bonus: bonus{bonusCritical: 10}},
			{bonus: bonus{bonusCritical: 10}},
			{bonus: bonus{bonusCritical: 3}},
		},
	}}
	got := equipmentBonuses(pools)
	if got[bonusCritical] != 23 {
		t.Fatalf("artifact critical = %d, want 23", got[bonusCritical])
	}
}

func TestMaxDamageRepeatsRings(t *testing.T) {
	attack := bonus{bonusAir: 10}
	pools := []bonusPool{{
		rule: slotRule{kind: "ring", count: 2},
		items: []bonusItem{
			{code: "ring", bonus: bonus{bonusAir: 50}},
			{code: "ring", bonus: bonus{bonusAir: 50}},
		},
	}}
	got := maxDamage(attack, pools, schemas.MonsterSchema{})
	if got != 20 {
		t.Fatalf("damage = %d, want 20", got)
	}
}

func TestPoolCopiesRespectsAvailability(t *testing.T) {
	ring := slotRule{kind: "ring", count: 2}
	one := poolCopies(ring, map[string]int{"ring": 1}, "ring")
	if one != 1 {
		t.Fatalf("ring copies = %d, want 1", one)
	}
	many := poolCopies(ring, map[string]int{"ring": 5}, "ring")
	if many != 2 {
		t.Fatalf("ring copies = %d, want 2", many)
	}
	unlimited := poolCopies(ring, nil, "ring")
	if unlimited != 2 {
		t.Fatalf("ring copies = %d, want 2", unlimited)
	}
	artifact := slotRule{kind: "artifact", count: 3, distinct: true}
	distinct := poolCopies(artifact, map[string]int{"artifact": 5}, "artifact")
	if distinct != 1 {
		t.Fatalf("artifact copies = %d, want 1", distinct)
	}
}

func TestEquipmentPoolsFiltersByAvailability(t *testing.T) {
	options := Options{Available: map[string]int{"topaz_ring": 1}}
	pools := equipmentPools(schemas.CharacterSchema{Level: 50}, options)
	for _, pool := range pools {
		if pool.rule.kind != "ring" {
			if len(pool.items) != 0 {
				t.Fatalf("pool %s has %d items, want 0", pool.rule.kind, len(pool.items))
			}
			continue
		}
		if len(pool.items) != 1 || pool.items[0].code != "topaz_ring" {
			t.Fatalf("ring pool = %v, want one topaz_ring", pool.items)
		}
	}
}

func TestEquipmentPoolsExcludesWithout(t *testing.T) {
	available := map[string]int{"topaz_ring": 1}
	options := Options{
		Available: available,
		Without:   map[string]bool{"topaz_ring": true},
	}
	pools := equipmentPools(schemas.CharacterSchema{Level: 50}, options)
	for _, pool := range pools {
		if pool.rule.kind == "ring" && len(pool.items) != 0 {
			t.Fatalf("ring pool = %v, want empty", pool.items)
		}
	}
}

func TestMaxDamageUsesDistinctArtifacts(t *testing.T) {
	attack := bonus{bonusAir: 10}
	pools := []bonusPool{{
		rule: slotRule{kind: "artifact", count: 3, distinct: true},
		items: []bonusItem{
			{code: "a", bonus: bonus{bonusAir: 50}},
			{code: "b", bonus: bonus{bonusAir: 50}},
			{code: "c", bonus: bonus{bonusAir: 50}},
			{code: "d", bonus: bonus{bonusAir: 50}},
		},
	}}
	got := maxDamage(attack, pools, schemas.MonsterSchema{})
	if got != 25 {
		t.Fatalf("damage = %d, want 25", got)
	}
}
