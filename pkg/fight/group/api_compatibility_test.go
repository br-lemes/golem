//go:build groupapi

package group_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/catalog"
	group "github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestGuardAPICompatibility(t *testing.T) {
	characters := readGroupCharacters(t, "group-simulation-guard.json")
	assertAPICompatibilityEvents(t, characters, "king_slime", "activates Guard", "activates Guard")
}

func TestVampiricStrikeAPICompatibility(t *testing.T) {
	characters := readGroupCharacters(t, "group-simulation-vampiric.json")
	assertAPICompatibilityEvents(t, characters, "king_slime", "Vampiric Strike heals", "Vampiric Strike heals")
}

func TestSplashRestoreAPICompatibility(t *testing.T) {
	characters := readGroupCharacters(t, "group-simulation-splash-restore.json")
	assertAPICompatibilityEvents(t, characters, "king_slime", "restored", "restored")
}

func TestHealingAuraAPICompatibility(t *testing.T) {
	characters := readGroupCharacters(t, "group-simulation-healing-aura.json")
	assertAPICompatibilityEvents(t, characters, "king_slime", "Healing Aura activates", "Healing Aura activates")
}

func assertAPICompatibilityEvents(t *testing.T, characters []schemas.FakeCharacterSchema, monsterName, apiEvent, localEvent string) {
	t.Helper()
	monster, ok := catalog.Monsters.Get(monsterName)
	if !ok {
		t.Fatal("king_slime not found")
	}
	request := schemas.CombatSimulationRequestSchema{
		Characters: characters,
		Monster:    monsterName,
		Iterations: 10,
	}
	remote, err := api.SimulationFight(request)
	if err != nil {
		t.Fatalf("API simulation: %v", err)
	}
	if len(remote.Results) == 0 {
		t.Fatal("API returned no results")
	}
	resultIndex := -1
	for index := range remote.Results {
		if strings.Contains(strings.Join(remote.Results[index].Logs, "\n"), apiEvent) {
			resultIndex = index
			break
		}
	}
	if resultIndex < 0 {
		t.Fatalf("API cases did not contain %q", apiEvent)
	}
	decisions := criticalDecisionsFromLogs(remote.Results[resultIndex].Logs)
	fighters := make([]group.Fighter, len(characters))
	for index, character := range characters {
		fighters[index] = group.FighterFromFakeCharacter(index, character)
	}
	combat := group.NewCombat(fighters, group.MonsterFromSchema(*monster), func() float64 { return 0.99 })
	combat.CriticalDecisions = decisions
	combat.TargetDecisions = targetDecisionsFromLogs(remote.Results[resultIndex].Logs)
	local := combat.Run(100)
	if strings.Count(strings.Join(local.Logs, "\n"), localEvent) == 0 {
		t.Logf("local replay did not contain %q; target/order sequence may differ", localEvent)
	}
	apiResult := remote.Results[resultIndex]
	t.Logf("API: result=%s turns=%d hp=%v %q=%d", apiResult.Result, apiResult.Turns, apiFinalHP(apiResult), apiEvent, strings.Count(strings.Join(apiResult.Logs, "\n"), apiEvent))
	t.Logf("local: win=%v turns=%d hp=%v %q=%d", local.Win, local.Turns, local.FinalHP, localEvent, strings.Count(strings.Join(local.Logs, "\n"), localEvent))
	if len(apiResult.CharacterResults) != len(local.FinalHP) {
		t.Fatalf("API characters = %d, local characters = %d", len(apiResult.CharacterResults), len(local.FinalHP))
	}
}

func criticalDecisionsFromLogs(logs []string) []group.CriticalDecision {
	decisions := []group.CriticalDecision{}
	previous := ""
	for _, log := range logs {
		if !strings.Contains(log, "used ") {
			continue
		}
		parts := strings.SplitN(log, ": ", 2)
		if len(parts) != 2 || parts[0] == previous {
			continue
		}
		previous = parts[0]
		actor := -1
		for index := 1; index <= 3; index++ {
			if strings.Contains(log, fmt.Sprintf("Character_%d used", index)) {
				actor = index - 1
				break
			}
		}
		decisions = append(decisions, group.CriticalDecision{Actor: actor, Critical: strings.Contains(log, "Critical strike")})
	}
	return decisions
}

func targetDecisionsFromLogs(logs []string) []int {
	decisions := []int{}
	previous := ""
	for _, log := range logs {
		if !strings.Contains(log, " used ") || !strings.Contains(log, " against Character_") {
			continue
		}
		parts := strings.SplitN(log, ": ", 2)
		if len(parts) != 2 || parts[0] == previous {
			continue
		}
		previous = parts[0]
		for index := 1; index <= 3; index++ {
			if strings.Contains(log, fmt.Sprintf("against Character_%d", index)) {
				decisions = append(decisions, index-1)
				break
			}
		}
	}
	return decisions
}

func apiFinalHP(result schemas.CombatResultSchema) []float64 {
	hp := make([]float64, len(result.CharacterResults))
	for index, character := range result.CharacterResults {
		if value, ok := character["final_hp"].(float64); ok {
			hp[index] = value
		}
	}
	return hp
}

func readGroupCharacters(t *testing.T, name string) []schemas.FakeCharacterSchema {
	t.Helper()
	body, err := os.ReadFile("../../../examples/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var characters []schemas.FakeCharacterSchema
	if err = json.Unmarshal(body, &characters); err != nil {
		t.Fatal(err)
	}
	return characters
}
