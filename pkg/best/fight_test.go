package best

import (
	"testing"

	"github.com/br-lemes/golem/pkg/database"
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

func TestFindWithAvailableDoesNotDependOnInitialEquipment(t *testing.T) {
	monster, ok := database.Monsters.Get("chicken")
	if !ok {
		t.Fatal("chicken not found")
	}

	available := make(map[string]int)
	for _, item := range database.Items.All() {
		if item.Type != "tool" && item.Subtype != "tool" {
			available[item.Code] = 2
		}
	}

	first := schemas.CharacterSchema{Level: 35, BootsSlot: "old_boots"}
	second := schemas.CharacterSchema{
		Level:     35,
		BootsSlot: "hard_leather_boots",
	}
	first.Wisdom = itemMeta(first.BootsSlot, "wisdom")
	first.Prospecting = itemMeta(first.BootsSlot, "prospecting")
	second.Wisdom = itemMeta(second.BootsSlot, "wisdom")
	second.Prospecting = itemMeta(second.BootsSlot, "prospecting")

	firstResult, err := findWithAvailable(first, *monster, available, false, false)
	if err != nil {
		t.Fatal(err)
	}
	secondResult, err := findWithAvailable(second, *monster, available, false, false)
	if err != nil {
		t.Fatal(err)
	}

	sameWinrate := firstResult.Winrate == secondResult.Winrate
	sameCycle := firstResult.CycleCost == secondResult.CycleCost
	sameTurns := firstResult.AverageTurns == secondResult.AverageTurns
	sameHP := firstResult.AverageFinalHP == secondResult.AverageFinalHP
	sameWisdom := firstResult.Wisdom == secondResult.Wisdom
	sameProspecting := firstResult.Prospecting == secondResult.Prospecting
	if !sameWinrate || !sameCycle || !sameTurns || !sameHP || !sameWisdom || !sameProspecting {
		t.Fatalf("initial equipment changed objective result: first=%+v second=%+v", firstResult, secondResult)
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

func TestCatalogLevel50BeamCandidates(t *testing.T) {
	available := make(map[string]int)
	for _, item := range database.Items.All() {
		if item.Type != "tool" && item.Subtype != "tool" {
			available[item.Code] = 2
		}
	}
	character := schemas.CharacterSchema{Level: 50}
	for _, code := range []string{"sandwarden", "red_dragon", "fennec"} {
		monster, ok := database.Monsters.Get(code)
		if !ok {
			t.Fatalf("%s not found", code)
		}
		result, err := findWithAvailable(character, *monster, available, false, false)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s: winrate=%v cycle=%v turns=%v hp=%v", code, result.Winrate, result.CycleCost, result.AverageTurns, result.AverageFinalHP)
	}
}
