package cmd

import (
	"fmt"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var groupSimulationCriticalCmd = &cobra.Command{
	Args:              cobra.ExactArgs(1),
	Use:               "critical <monster>",
	Short:             "Compare the group simulator using API criticals",
	ValidArgsFunction: completion.Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := readGroupSimulationFlags(cmd)
		if err != nil {
			return err
		}
		monster, err := groupSimulationMonster(args)
		if err != nil {
			return err
		}
		characters, err := group.CharactersFromFile(flags.File)
		if err != nil {
			return err
		}
		request := schemas.CombatSimulationRequestSchema{
			Characters: characters,
			Monster:    args[0],
			Iterations: 1,
		}
		remote, err := api.SimulationFight(request)
		if err != nil {
			return err
		}
		if len(remote.Results) == 0 {
			return fmt.Errorf("API returned no combat result")
		}
		sequence := groupCriticalSequence(remote.Results[0].Logs)
		targets := groupTargetSequence(remote.Results[0].Logs)
		actions := groupActionSequence(remote.Results[0].Logs)
		bubbles := groupBubbleSequence(remote.Results[0].Logs)
		fighters, err := group.FightersFromFile(flags.File)
		if err != nil {
			return err
		}
		localResults := make([]schemas.CombatResultSchema, 0, flags.Iterations)
		for i := 0; i < flags.Iterations; i++ {
			combat := group.NewCombat(fighters, group.MonsterFromSchema(*monster), nil)
			combat.CriticalDecisions = sequence
			combat.TargetDecisions = targets
			combat.ActionDecisions = actions
			combat.BubbleDecisions = bubbles
			combat.GuardDecisions = groupGuardSequence(remote.Results[0].Logs)
			localResults = append(localResults, groupSimulationResult(combat.Run(100), true))
		}
		diagnostics := groupSimulationDiagnostics(remote.Results)
		if !flags.Logs {
			remote.Results[0].Logs = []string{}
		}
		apiOutput := groupSimulationOutput{
			Losses:      remote.Losses,
			Results:     remote.Results,
			Winrate:     remote.Winrate,
			Wins:        remote.Wins,
			Diagnostics: diagnostics,
		}
		comparison := groupSimulationComparison{
			API:   apiOutput,
			Local: groupSimulationReport(localResults, flags.Logs),
		}
		comparison.Differences = groupDiagnosticDifferences(apiOutput.Diagnostics, comparison.Local.Diagnostics)
		apiDamage := groupDamageSequence(remote.Results[0].Logs)
		localDamage := groupDamageSequence(localResults[0].Logs)
		apiActors := groupDamageActors(remote.Results[0].Logs)
		localActors := groupDamageActors(localResults[0].Logs)
		limit := len(apiDamage)
		if len(localDamage) < limit {
			limit = len(localDamage)
		}
		for index := 0; index < limit; index++ {
			if apiActors[index] != localActors[index] {
				comparison.Differences = append(comparison.Differences, fmt.Sprintf("damage_actor_at_hit_%d_api_%s_local_%s", index+1, apiActors[index], localActors[index]))
				break
			}
			if apiDamage[index] != localDamage[index] {
				comparison.Differences = append(comparison.Differences, fmt.Sprintf("damage_at_hit_%d_api_%d_local_%d", index+1, apiDamage[index], localDamage[index]))
				break
			}
		}
		if len(apiDamage) != len(localDamage) {
			comparison.Differences = append(comparison.Differences, "damage_sequence_length")
		}
		return console.Auto(comparison)
	},
}

func init() {
	groupSimulationCmd.AddCommand(groupSimulationCriticalCmd)
	err := utils.RegisterFlags[groupSimulationFlags](groupSimulationCriticalCmd)
	if err != nil {
		panic(err)
	}
}

func groupGuardSequence(logs []string) []group.GuardDecision {
	decisions := []group.GuardDecision{}
	for _, log := range logs {
		if !strings.Contains(log, "activates Guard to protect Character_") {
			continue
		}
		turn := 0
		protector, target := 0, 0
		format := "Turn %d: Character_%d activates Guard to protect"
		format += " Character_%d"
		_, err := fmt.Sscanf(log, format, &turn, &protector, &target)
		if err == nil {
			decision := group.GuardDecision{
				Turn:      turn,
				Protector: protector - 1,
				Target:    target - 1,
			}
			decisions = append(decisions, decision)
		}
	}
	return decisions
}

func groupDiagnosticDifferences(api, local map[string]interface{}) []string {
	apiCharacters, _ := api["characters"].(map[string]map[string]interface{})
	localCharacters, _ := local["characters"].(map[string]map[string]interface{})
	differences := []string{}
	for name, apiValues := range apiCharacters {
		localValues := localCharacters[name]
		for key, apiValue := range apiValues {
			if intValue(apiValue) != intValue(localValues[key]) {
				differences = append(differences, fmt.Sprintf("diagnostics.%s.%s", name, key))
			}
		}
		for key := range localValues {
			_, ok := apiValues[key]
			if !ok {
				differences = append(differences, fmt.Sprintf("diagnostics.%s.%s", name, key))
			}
		}
	}
	for name := range localCharacters {
		if apiCharacters[name] == nil {
			differences = append(differences, fmt.Sprintf("diagnostics.%s", name))
		}
	}
	return differences
}

func groupBubbleSequence(logs []string) []int {
	elements := []string{"fire", "earth", "water", "air"}
	sequence := []int{}
	for _, log := range logs {
		if !strings.Contains(log, "Protective Bubble") {
			continue
		}
		for index, element := range elements {
			if strings.Contains(log, element+" resistance") {
				sequence = append(sequence, index)
				break
			}
		}
	}
	return sequence
}

func groupActionSequence(logs []string) []group.ActionDecision {
	decisions := []group.ActionDecision{}
	previous := ""
	for _, log := range logs {
		if !strings.Contains(log, " used ") {
			continue
		}
		parts := strings.SplitN(log, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0] + ":" + strings.SplitN(parts[1], " used ", 2)[0]
		if key == previous {
			continue
		}
		previous = key
		turnParts := strings.Split(parts[0], " ")
		if len(turnParts) != 2 {
			continue
		}
		turn := 0
		_, err := fmt.Sscanf(turnParts[1], "%d", &turn)
		if err != nil {
			continue
		}
		actor := -1
		if strings.HasPrefix(parts[1], "Character_") {
			_, err = fmt.Sscanf(parts[1], "Character_%d", &actor)
			if err != nil {
				continue
			}
			actor--
		}
		decision := group.ActionDecision{Turn: turn, Actor: actor}
		decisions = append(decisions, decision)
	}
	return decisions
}

func groupTargetSequence(logs []string) []int {
	targets := []int{}
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
				targets = append(targets, index-1)
				break
			}
		}
	}
	return targets
}

func groupCriticalSequence(logs []string) []group.CriticalDecision {
	sequence := []group.CriticalDecision{}
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
		sequence = append(sequence, group.CriticalDecision{
			Actor:    actor,
			Critical: strings.Contains(log, "Critical strike"),
		})
	}
	return sequence
}

func groupDamageSequence(logs []string) []int {
	sequence := []int{}
	for _, log := range logs {
		if !strings.Contains(log, " used ") || !strings.Contains(log, " dealt ") {
			continue
		}
		sequence = append(sequence, groupDiagnosticNumber(log, "dealt"))
	}
	return sequence
}

func groupDamageActors(logs []string) []string {
	sequence := []string{}
	for _, log := range logs {
		if !strings.Contains(log, " used ") || !strings.Contains(log, " dealt ") {
			continue
		}
		parts := strings.SplitN(log, ": ", 2)
		if len(parts) == 2 {
			sequence = append(sequence, strings.SplitN(parts[1], " used ", 2)[0])
		}
	}
	return sequence
}
