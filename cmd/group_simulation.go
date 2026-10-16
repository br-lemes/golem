package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type groupSimulationFlags struct {
	File       string `flag:"file" desc:"JSON file containing 1 to 3 simulation characters"`
	Iterations int    `flag:"iterations" default:"1" desc:"Number of simulations to run"`
	Logs       bool   `flag:"logs" desc:"Include combat logs"`
}

type groupSimulationOutput struct {
	Losses      int                          `json:"losses"`
	Results     []schemas.CombatResultSchema `json:"results"`
	Winrate     float32                      `json:"winrate"`
	Wins        int                          `json:"wins"`
	Diagnostics map[string]interface{}       `json:"diagnostics"`
}

type groupSimulationComparison struct {
	API         groupSimulationOutput `json:"api"`
	Local       groupSimulationOutput `json:"local"`
	Differences []string              `json:"relevant_differences,omitempty"`
}

var groupSimulationCmd = &cobra.Command{
	Use:   "group-simulation",
	Short: "Run the experimental group simulator",
}

func groupSimulationMonster(args []string) (*schemas.MonsterSchema, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected exactly one monster")
	}
	monster, found := catalog.Monsters.Get(args[0])
	if !found {
		return nil, fmt.Errorf("invalid monster: %s", args[0])
	}
	return monster, nil
}

func readGroupSimulationFlags(cmd *cobra.Command) (groupSimulationFlags, error) {
	flags, err := utils.ReadFlags[groupSimulationFlags](cmd)
	if err != nil {
		return groupSimulationFlags{}, err
	}
	if flags.File == "" {
		return groupSimulationFlags{}, fmt.Errorf("--file is required")
	}
	if flags.Iterations < 1 {
		return groupSimulationFlags{}, fmt.Errorf("iterations must be at least 1")
	}
	return flags, nil
}

func groupSimulationResult(result group.Result, logs bool) schemas.CombatResultSchema {
	characters := make([]map[string]interface{}, len(result.FinalHP))
	for i, hp := range result.FinalHP {
		characters[i] = map[string]interface{}{
			"final_hp":               hp,
			"utility1_slot_quantity": 0.0,
			"utility2_slot_quantity": 0.0,
		}
	}
	outputLogs := []string{}
	if logs {
		outputLogs = result.Logs
	}
	name := "loss"
	if result.Win {
		name = "win"
	}
	return schemas.CombatResultSchema{
		CharacterResults: characters,
		Logs:             outputLogs,
		Result:           name,
		Turns:            result.Turns,
	}
}

func groupSimulationResults(results []group.Result, logs bool) []schemas.CombatResultSchema {
	output := make([]schemas.CombatResultSchema, len(results))
	for i, result := range results {
		output[i] = groupSimulationResult(result, logs)
	}
	return output
}

func groupSimulationReport(results []schemas.CombatResultSchema, includeLogs bool) groupSimulationOutput {
	wins := 0
	for _, result := range results {
		if result.Result == "win" {
			wins++
		}
	}
	winrate := float32(0)
	if len(results) > 0 {
		winrate = float32(wins) * 100 / float32(len(results))
	}
	output := groupSimulationOutput{
		Losses:      len(results) - wins,
		Results:     results,
		Winrate:     winrate,
		Wins:        wins,
		Diagnostics: groupSimulationDiagnostics(results),
	}
	if !includeLogs {
		for i := range output.Results {
			output.Results[i].Logs = []string{}
		}
	}
	return output
}

func groupSimulationDiagnostics(results []schemas.CombatResultSchema) map[string]interface{} {
	var turns, hp float32
	for _, result := range results {
		turns += float32(result.Turns)
		for _, character := range result.CharacterResults {
			value, ok := character["final_hp"].(float64)
			if ok {
				hp += float32(value)
			}
		}
	}
	if len(results) > 0 {
		turns /= float32(len(results))
		hp /= float32(len(results))
	}
	player := map[string]interface{}{}
	monster := map[string]interface{}{}
	characters := map[string]map[string]interface{}{}
	for _, result := range results {
		lastCritical := ""
		for _, log := range result.Logs {
			character := groupDiagnosticCharacter(log)
			if character >= 0 && strings.Contains(log, ": Character_") && strings.Contains(log, " used ") && strings.Contains(log, " dealt ") {
				incrementCharacter(characters, character, "attacks")
				incrementCharacterBy(characters, character, "damage_dealt", groupDiagnosticNumber(log, "dealt"))
			}
			if character >= 0 && strings.Contains(log, "against Character_") && strings.Contains(log, " dealt ") {
				incrementCharacter(characters, character, "attacks_received")
				incrementCharacterBy(characters, character, "damage_received", groupDiagnosticNumber(log, "dealt"))
			}
			isPlayer := strings.Contains(log, "Character_") && strings.Contains(log, "used ")
			side := monster
			if isPlayer {
				side = player
			}
			if strings.Contains(log, "Critical strike") {
				key := groupDiagnosticTurn(log)
				if key != lastCritical {
					side["critical_hits"] = intValue(side["critical_hits"]) + 1
					lastCritical = key
					if character >= 0 {
						incrementCharacter(characters, character, "critical_hits")
					}
				}
			}
			if character >= 0 && strings.Contains(log, "has been defeated!") {
				incrementCharacter(characters, character, "defeats")
			}
			if strings.Contains(log, "raises a barrier") || strings.Contains(log, "refreshes its barrier") {
				monster["barrier_activations"] = intValue(monster["barrier_activations"]) + 1
			}
			if strings.Contains(log, "barrier is destroyed") {
				monster["barrier_destructions"] = intValue(monster["barrier_destructions"]) + 1
			}
			if strings.Contains(log, "barrier absorbs") {
				value := groupDiagnosticNumber(log, "absorbs")
				monster["barrier_absorbed_damage"] = intValue(monster["barrier_absorbed_damage"]) + value
			}
			if strings.Contains(log, "Guard redirects") {
				incrementCharacter(characters, character, "guard_redirects")
				incrementCharacterBy(characters, character, "damage_redirected", groupDiagnosticNumber(log, "redirects"))
			}
			if strings.Contains(log, "activates Guard") {
				incrementCharacter(characters, character, "guard_activations")
			}
			if strings.Contains(log, "shell activates") {
				player["shell_activations"] = intValue(player["shell_activations"]) + 1
				incrementCharacter(characters, character, "shell_activations")
			}
			if strings.Contains(log, "shell effect has worn off") {
				player["shell_expirations"] = intValue(player["shell_expirations"]) + 1
				incrementCharacter(characters, character, "shell_expirations")
			}
			if strings.Contains(log, "Greed empowers Character_") {
				player["greed_activations"] = intValue(player["greed_activations"]) + 1
				incrementCharacter(characters, character, "greed_activations")
			}
			if strings.Contains(log, "Greed empowers the monster") {
				monster["greed_activations"] = intValue(monster["greed_activations"]) + 1
			}
			if strings.Contains(log, "applies a burn") {
				if strings.Contains(log, "The monster applies") {
					monster["burn_applications"] = intValue(monster["burn_applications"]) + 1
				} else {
					player["burn_applications"] = intValue(player["burn_applications"]) + 1
					incrementCharacter(characters, character, "burn_applications")
				}
			}
			if strings.Contains(log, "suffers from burn") {
				player["burn_ticks"] = intValue(player["burn_ticks"]) + 1
				incrementCharacter(characters, character, "burn_ticks")
			}
			if strings.Contains(log, "suffers") && strings.Contains(log, "poison") {
				player["poison_ticks"] = intValue(player["poison_ticks"]) + 1
				incrementCharacter(characters, character, "poison_ticks")
			}
			if strings.Contains(log, "applies a poison") {
				if strings.Contains(log, "monster applies") {
					monster["poison_applications"] = intValue(monster["poison_applications"]) + 1
				} else {
					player["poison_applications"] = intValue(player["poison_applications"]) + 1
					incrementCharacter(characters, character, "poison_applications")
				}
			}
			if strings.Contains(log, "from lifesteal") {
				if strings.Contains(log, "Character_") {
					player["lifesteal_activations"] = intValue(player["lifesteal_activations"]) + 1
					incrementCharacter(characters, character, "lifesteal_activations")
				} else {
					monster["lifesteal_activations"] = intValue(monster["lifesteal_activations"]) + 1
				}
			}
			if strings.Contains(log, "Enchanted Mirror activates") {
				player["mirror_activations"] = intValue(player["mirror_activations"]) + 1
			}
			if strings.Contains(log, "Frenzy triggers on critical") {
				player["frenzy_activations"] = intValue(player["frenzy_activations"]) + 1
				incrementCharacter(characters, character, "frenzy_activations")
			}
			if strings.Contains(log, "Frenzy activates") {
				monster["frenzy_activations"] = intValue(monster["frenzy_activations"]) + 1
			}
			if strings.Contains(log, "Vampiric Strike heals") {
				player["lifesteal_activations"] = intValue(player["lifesteal_activations"]) + 1
				incrementCharacter(characters, character, "vampiric_strike_activations")
			}
			if strings.Contains(log, "heals ") && strings.Contains(log, "Healing effect") {
				player["healing_activations"] = intValue(player["healing_activations"]) + 1
				incrementCharacter(characters, character, "healing_activations")
			}
			if strings.Contains(log, "Healing Aura heals") {
				player["healing_aura_activations"] = intValue(player["healing_aura_activations"]) + 1
				incrementCharacter(characters, character, "healing_aura_activations")
			}
			if strings.Contains(log, "resistance is corrupted") {
				player["corrupted_applications"] = intValue(player["corrupted_applications"]) + 1
				incrementCharacter(characters, character, "corrupted_applications")
			}
			if strings.Contains(log, "Berserker Rage activates") {
				monster["berserker_activations"] = intValue(monster["berserker_activations"]) + 1
			}
			if strings.Contains(log, "uses Reconstitution") {
				monster["reconstitution_activations"] = intValue(monster["reconstitution_activations"]) + 1
			}
			if strings.Contains(log, "uses Void Drain") {
				monster["void_drain_activations"] = intValue(monster["void_drain_activations"]) + 1
			}
			if strings.Contains(log, "used antipoison") || strings.Contains(log, "used antidote") {
				player["antidotes"] = intValue(player["antidotes"]) + 1
				incrementCharacter(characters, character, "antidotes")
			}
			if strings.Contains(log, "restored ") {
				player["healing_potions"] = intValue(player["healing_potions"]) + 1
				incrementCharacter(characters, character, "healing_potions")
				incrementCharacterBy(characters, character, "hp_restored", groupDiagnosticNumber(log, "restored"))
			}
			if character >= 0 && strings.Contains(log, "heals ") {
				incrementCharacterBy(characters, character, "hp_restored", groupDiagnosticNumber(log, "heals"))
			}
		}
	}
	return map[string]interface{}{
		"average_turns":          turns,
		"average_final_hp":       hp,
		"average_fight_cooldown": turns * 2,
		"xp":                     0,
		"xp_per_cycle":           0.0,
		"player":                 player,
		"characters":             characters,
		"monster":                monster,
	}
}

func groupDiagnosticCharacter(log string) int {
	marker := "Character_"
	index := strings.Index(log, marker)
	if index < 0 {
		return -1
	}
	value := strings.TrimLeft(log[index+len(marker):], "0123456789")
	number := log[index+len(marker) : len(log)-len(value)]
	parsed, err := strconv.Atoi(number)
	if err != nil || parsed < 1 {
		return -1
	}
	return parsed - 1
}

func incrementCharacter(characters map[string]map[string]interface{}, index int, key string) {
	incrementCharacterBy(characters, index, key, 1)
}

func incrementCharacterBy(characters map[string]map[string]interface{}, index int, key string, amount int) {
	if index < 0 {
		return
	}
	name := fmt.Sprintf("Character_%d", index+1)
	if characters[name] == nil {
		characters[name] = map[string]interface{}{}
	}
	characters[name][key] = intValue(characters[name][key]) + amount
}

func intValue(value interface{}) int {
	result, _ := value.(int)
	return result
}

func groupDiagnosticTurn(log string) string {
	parts := strings.SplitN(log, ": ", 2)
	if len(parts) != 2 {
		return log
	}
	return parts[0] + ":" + strings.SplitN(parts[1], " used ", 2)[0]
}

func groupDiagnosticNumber(log, marker string) int {
	parts := strings.Split(log, marker+" ")
	if len(parts) < 2 {
		return 0
	}
	value := strings.Fields(parts[1])
	if len(value) == 0 {
		return 0
	}
	result, _ := strconv.Atoi(value[0])
	return result
}

func init() {
	rootCmd.AddCommand(groupSimulationCmd)
}
