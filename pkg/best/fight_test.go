package best

import (
	"testing"

	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSimulationItemScoreUsesWeaponElement(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{
		{Code: "dmg", Value: 10},
		{Code: "dmg_fire", Value: 20},
		{Code: "dmg_water", Value: 30},
	}
	item := schemas.ItemSchema{Effects: effects}

	fire := simulationItemScore(item, schemas.MonsterSchema{}, "fire", "damage")
	water := simulationItemScore(item, schemas.MonsterSchema{}, "water", "damage")
	if fire == water {
		t.Fatalf("element-specific damage should depend on weapon element: fire=%v water=%v", fire, water)
	}
}

func TestEffectivePlayerDamageAppliesResistanceFraction(t *testing.T) {
	player := fight.Fighter{Stats: fight.Stats{AttackFire: 100}}
	monster := schemas.MonsterSchema{ResFire: 25}
	got := effectivePlayerDamage(player, monster)
	if got != 75 {
		t.Fatalf("effective damage = %v, want 75", got)
	}
}

func TestAlignSimulationArtifactsPreservesEquivalentOrder(t *testing.T) {
	character := schemas.CharacterSchema{
		Artifact1Slot: "novice_guide",
		Artifact2Slot: "lost_world_map",
	}
	slots := map[string]string{
		"artifact1": "lost_world_map",
		"artifact2": "novice_guide",
		"artifact3": "lich_race_medal",
	}

	got := alignSimulationArtifacts(character, slots)
	if got["artifact1"] != "novice_guide" || got["artifact2"] != "lost_world_map" || got["artifact3"] != "lich_race_medal" {
		t.Fatalf("artifact order = %#v", got)
	}
}

func TestSimulationEquipmentChangesOmitsEmptyResults(t *testing.T) {
	current := map[string]string{"boots": "", "weapon": "new_weapon"}
	original := map[string]string{"boots": "old_boots", "weapon": "old_weapon"}

	got := simulationEquipmentChanges(current, original)
	_, bootsChanged := got["boots"]
	if bootsChanged {
		t.Fatal("empty equipment change should be omitted")
	}
	if got["weapon"] != "new_weapon" {
		t.Fatalf("weapon change = %q", got["weapon"])
	}
}

func TestSimulationUtilityChangesIncludesRemoval(t *testing.T) {
	character := schemas.CharacterSchema{Utility1Slot: "health_potion"}
	current := map[string]string{"utility1": "", "utility2": ""}

	got := simulationUtilityChanges(current, character)
	value, utilityChanged := got["utility1"]
	if !utilityChanged || value != "" {
		t.Fatalf("utility removal = %#v", got)
	}
	_, emptyUtilityChanged := got["utility2"]
	if emptyUtilityChanged {
		t.Fatal("unchanged empty utility should be omitted")
	}
}

func TestLoadoutHasQuantityRejectsDuplicateUtilities(t *testing.T) {
	slots := map[string]string{
		"utility1": "health_potion",
		"utility2": "health_potion",
	}
	available := map[string]int{"health_potion": 2}
	if loadoutHasQuantity(slots, available) {
		t.Fatal("duplicate utility should not be allowed")
	}
}

func TestBetterSimulationScorePrefersLowerCycleCost(t *testing.T) {
	faster := Result{Winrate: 100, CycleCost: 84}
	slower := Result{Winrate: 100, CycleCost: 94}
	if !betterSimulationScore(faster, slower) {
		t.Fatal("lower cycle cost should be preferred")
	}
	if betterSimulationScore(slower, faster) {
		t.Fatal("higher cycle cost should not be preferred")
	}
}
