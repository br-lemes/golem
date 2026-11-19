package cmd

import (
	"testing"

	"github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestGroupSimulationReportKeepsDiagnosticsWithoutLogs(t *testing.T) {
	result := schemas.CombatResultSchema{
		CharacterResults: []map[string]interface{}{{"final_hp": float64(100)}},
		Logs: []string{
			"Turn 1: Character_1 used earth attack and dealt 20 damage (Critical strike).",
			"Turn 2: Greed empowers Character_1 (+10% damage, total +10%).",
		},
		Result: "win",
		Turns:  2,
	}

	withLogs := groupSimulationReport([]schemas.CombatResultSchema{result}, true, nil, schemas.MonsterSchema{})
	withoutLogs := groupSimulationReport([]schemas.CombatResultSchema{result}, false, nil, schemas.MonsterSchema{})
	if len(withLogs.Results[0].Logs) != 2 {
		t.Fatalf("logs with --logs = %d, want 2", len(withLogs.Results[0].Logs))
	}
	if len(withoutLogs.Results[0].Logs) != 0 {
		t.Fatalf("logs without --logs = %d, want 0", len(withoutLogs.Results[0].Logs))
	}
	if intValue(withLogs.Diagnostics["player"].(map[string]interface{})["greed_activations"]) != 1 {
		t.Fatal("diagnostics with logs missed Greed activation")
	}
	if intValue(withoutLogs.Diagnostics["player"].(map[string]interface{})["greed_activations"]) != 1 {
		t.Fatal("diagnostics without logs missed Greed activation")
	}
}

func TestGroupSimulationDiagnosticsAppliesHasteAndMinimumCooldown(t *testing.T) {
	results := []schemas.CombatResultSchema{{Turns: 10}}
	fighters := []group.Fighter{{Stats: group.Stats{Haste: 50}}}
	diagnostics := groupSimulationDiagnostics(results, fighters, schemas.MonsterSchema{})
	got := diagnostics["average_fight_cooldown"]
	if got != float32(10) {
		t.Fatalf("haste cooldown = %v, want 10", got)
	}
	minimumResults := []schemas.CombatResultSchema{{Turns: 1}}
	fighters = []group.Fighter{{Stats: group.Stats{Haste: 100}}}
	minimum := groupSimulationDiagnostics(minimumResults, fighters, schemas.MonsterSchema{})
	got = minimum["average_fight_cooldown"]
	if got != float32(5) {
		t.Fatalf("minimum cooldown = %v, want 5", got)
	}
}

func TestGroupSimulationDiagnosticsCalculatesGroupXP(t *testing.T) {
	fighters := []group.Fighter{
		{Level: 20},
		{Level: 24, Stats: group.Stats{Wisdom: 100}},
		{Level: 25},
	}
	monster := schemas.MonsterSchema{Level: 20, Hp: 1000}
	results := []schemas.CombatResultSchema{{Result: "win", Turns: 10}}
	diagnostics := groupSimulationDiagnostics(results, fighters, monster)
	if diagnostics["xp"] != 177 {
		t.Fatalf("group XP = %v, want 177", diagnostics["xp"])
	}
	characters := diagnostics["characters"].(map[string]map[string]interface{})
	if characters["Character_2"]["xp"] != 63 {
		t.Fatalf("character XP = %v, want 63", characters["Character_2"]["xp"])
	}
}

func TestGroupSimulationDiagnosticsRecognizesGroupEffects(t *testing.T) {
	result := schemas.CombatResultSchema{
		Logs: []string{
			"Turn 1: Character_1 used health_potion and restored 100 HP.",
			"Turn 2: Character_1 heals 4 HP from Healing effect.",
			"Turn 3: Character_1's Vampiric Strike heals Character_2 for 8 HP.",
			"Turn 4: Character_1's Healing Aura heals Character_2 for 6 HP.",
			"Turn 5: Character_1's resistance is corrupted.",
			"Turn 6: The monster's Berserker Rage activates, gaining 20% permanent damage!",
			"Turn 7: Character_1 used antipoison.",
			"Turn 8: The monster uses Reconstitution and restores all HP.",
			"Turn 9: The monster uses Void Drain and drains 10 HP from Character_1.",
		},
	}
	diagnostics := groupSimulationDiagnostics([]schemas.CombatResultSchema{result}, nil, schemas.MonsterSchema{})
	player := diagnostics["player"].(map[string]interface{})
	if intValue(player["healing_potions"]) != 1 {
		t.Fatal("health potion was not diagnosed")
	}
	if intValue(player["healing_activations"]) != 1 {
		t.Fatal("healing effect was not diagnosed")
	}
	if intValue(player["lifesteal_activations"]) != 1 {
		t.Fatal("Vampiric Strike was not diagnosed")
	}
	if intValue(player["healing_aura_activations"]) != 1 {
		t.Fatal("Healing Aura was not diagnosed")
	}
	if intValue(player["corrupted_applications"]) != 1 {
		t.Fatal("corruption was not diagnosed")
	}
	if intValue(player["antidotes"]) != 1 {
		t.Fatal("antipoison was not diagnosed")
	}
	monster := diagnostics["monster"].(map[string]interface{})
	if intValue(monster["berserker_activations"]) != 1 {
		t.Fatal("Berserker Rage was not diagnosed")
	}
	if intValue(monster["reconstitution_activations"]) != 1 {
		t.Fatal("Reconstitution was not diagnosed")
	}
	if intValue(monster["void_drain_activations"]) != 1 {
		t.Fatal("Void Drain was not diagnosed")
	}
}

func TestGroupActionSequenceDeduplicatesElementalHits(t *testing.T) {
	logs := []string{
		"Turn 1: Character_2 used earth attack and dealt 20 damage.",
		"Turn 1: Character_2 used fire attack and dealt 20 damage.",
		"Turn 2: King Slime used water attack against Character_1 and dealt 14 damage.",
		"Turn 2: King Slime used air attack against Character_1 and dealt 14 damage.",
	}
	sequence := groupActionSequence(logs)
	if len(sequence) != 2 {
		t.Fatalf("action sequence length = %d, want 2", len(sequence))
	}
	if sequence[0].Turn != 1 || sequence[0].Actor != 1 {
		t.Fatalf("first action = %+v, want turn 1 actor 1", sequence[0])
	}
	if sequence[1].Turn != 2 || sequence[1].Actor != -1 {
		t.Fatalf("second action = %+v, want turn 2 monster", sequence[1])
	}
}

func TestGroupSimulationDiagnosticsByCharacter(t *testing.T) {
	result := schemas.CombatResultSchema{
		Logs: []string{
			"Turn 3: Character_2 used earth attack and dealt 10 damage (Critical strike).",
			"Turn 4: Character_1's Vampiric Strike heals Character_2 for 8 HP.",
			"Turn 5: Character_1 used Health Potion and restored 100 HP. (HP: 100/200)",
			"Turn 6: King Slime used fire attack against Character_2 and dealt 13 damage. Character_2 HP: 87/100",
			"Turn 6: King Slime used earth attack against Character_2 and dealt 13 damage. Character_2 HP: 74/100",
			"Turn 7: Character_2 has been defeated!",
		},
	}
	diagnostics := groupSimulationDiagnostics([]schemas.CombatResultSchema{result}, nil, schemas.MonsterSchema{})
	characters := diagnostics["characters"].(map[string]map[string]interface{})
	if intValue(characters["Character_2"]["critical_hits"]) != 1 {
		t.Fatal("critical hit was assigned to the wrong character")
	}
	if intValue(characters["Character_1"]["vampiric_strike_activations"]) != 1 {
		t.Fatal("Vampiric Strike was assigned to the wrong character")
	}
	if intValue(characters["Character_1"]["healing_potions"]) != 1 {
		t.Fatal("potion use was not assigned to Character_1")
	}
	if intValue(characters["Character_2"]["damage_dealt"]) != 10 {
		t.Fatal("damage dealt was parsed incorrectly")
	}
	if intValue(characters["Character_2"]["damage_received"]) != 26 {
		t.Fatal("damage received was parsed incorrectly")
	}
	if intValue(characters["Character_1"]["damage_received"]) != 0 {
		t.Fatal("damage was assigned to the wrong character")
	}
	if intValue(characters["Character_2"]["defeats"]) != 1 {
		t.Fatal("defeat was not assigned to Character_2")
	}
}

func TestGroupDiagnosticDifferencesDetectsMissingMetrics(t *testing.T) {
	api := map[string]interface{}{
		"characters": map[string]map[string]interface{}{
			"Character_1": {"critical_hits": 2},
		},
	}
	local := map[string]interface{}{
		"characters": map[string]map[string]interface{}{
			"Character_1": {"critical_hits": 1, "damage_dealt": 10},
		},
	}
	differences := groupDiagnosticDifferences(api, local)
	if len(differences) != 2 {
		t.Fatalf("differences = %v, want two entries", differences)
	}
}

func TestGroupDamageSequenceReadsEachHit(t *testing.T) {
	logs := []string{
		"Turn 1: King Slime used fire attack against Character_1 and dealt 13 damage.",
		"Turn 1: King Slime used earth attack against Character_1 and dealt 20 damage (Critical strike).",
	}
	damages := groupDamageSequence(logs)
	if len(damages) != 2 || damages[0] != 13 || damages[1] != 20 {
		t.Fatalf("damage sequence = %v, want [13 20]", damages)
	}
}
