package fight

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/schemas"
)

func CompareCritical(request schemas.CombatSimulationRequestSchema, monster schemas.MonsterSchema, options SimulationOptions, includeLogs bool) (SimulationComparison, error) {
	iterations := request.Iterations
	request.Iterations = 1
	remote, err := api.SimulationFight(request)
	if err != nil {
		return SimulationComparison{}, err
	}
	if len(remote.Results) == 0 {
		return SimulationComparison{}, fmt.Errorf("API returned no combat result")
	}
	character := request.Characters[0]
	fighter := FromLoadout(character.Level, CharacterSlots(character), CharacterUtilities(character))
	options.Iterations = iterations
	options.CriticalSequence = CriticalSequenceFromLogs(remote.Results[0].Logs)
	if options.RNG == nil {
		options.RNG = defaultRNG()
	}
	local := SimulateMany(fighter, monster, options)
	comparison := SimulationComparison{
		API: SimulationStats{
			Wins:   CountSimulationWins(remote.Results),
			Losses: len(remote.Results) - CountSimulationWins(remote.Results),
		},
		Local: SimulationStats{
			Wins:    local.Wins,
			Losses:  local.Losses,
			Winrate: local.Winrate,
		},
	}
	comparison.API.Winrate = float32(comparison.API.Wins) * 100 / float32(len(remote.Results))
	apiReport := Report(fighter, character.Level, monster, remote.Results)
	localReport := Report(fighter, character.Level, monster, local.Results)
	comparison.API.Diagnostics = apiReport.Diagnostics
	comparison.Local.Diagnostics = localReport.Diagnostics
	if includeLogs && len(local.Results) > 0 {
		comparison.API.Logs = remote.Results[0].Logs
		comparison.Local.Logs = local.Results[0].Logs
	}
	return comparison, nil
}
