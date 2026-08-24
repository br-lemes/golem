package cmd

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/simulation"
	"github.com/spf13/cobra"
)

// SimulationComparison deliberately compares combat facts, not logs. Logs
// contain wording and formatting differences, while these fields describe
// the behavior users care about when validating the local implementation.
type simulationComparison struct {
	TotalIterations int                    `json:"total_iterations"`
	API             simulationStats        `json:"api"`
	Local           simulationStats        `json:"local"`
	Differences     []simulationDifference `json:"relevant_differences,omitempty"`
	APILogs         []string               `json:"api_logs,omitempty"`
	LocalLogs       []string               `json:"local_logs,omitempty"`
}

type simulationDiagnostics struct {
	AverageTurns         float32                   `json:"average_turns"`
	AverageFinalHP       float32                   `json:"average_final_hp"`
	AverageFightCooldown float32                   `json:"average_fight_cooldown"`
	XP                   int                       `json:"xp"`
	XPPerCycle           float32                   `json:"xp_per_cycle"`
	Player               simulationSideDiagnostics `json:"player"`
	Monster              simulationSideDiagnostics `json:"monster"`
}

type simulationSideDiagnostics struct {
	CriticalHits          int `json:"critical_hits,omitempty"`
	BurnApplications      int `json:"burn_applications,omitempty"`
	BurnTicks             int `json:"burn_ticks,omitempty"`
	PoisonApplications    int `json:"poison_applications,omitempty"`
	PoisonTicks           int `json:"poison_ticks,omitempty"`
	CorruptedApplications int `json:"corrupted_applications,omitempty"`
	HealingPotions        int `json:"healing_potions,omitempty"`
	HealingActivations    int `json:"healing_activations,omitempty"`
	LifestealActivations  int `json:"lifesteal_activations,omitempty"`
	BerserkerActivations  int `json:"berserker_activations,omitempty"`
	GreedActivations      int `json:"greed_activations,omitempty"`
	ShellActivations      int `json:"shell_activations,omitempty"`
	ShellExpirations      int `json:"shell_expirations,omitempty"`
	MirrorActivations     int `json:"mirror_activations,omitempty"`
	FrenzyActivations     int `json:"frenzy_activations,omitempty"`
	BarrierActivations    int `json:"barrier_activations,omitempty"`
	BarrierDestructions   int `json:"barrier_destructions,omitempty"`
	BarrierAbsorbedDamage int `json:"barrier_absorbed_damage,omitempty"`
	Antidotes             int `json:"antidotes,omitempty"`
}

type simulationStats struct {
	Wins        int                   `json:"wins"`
	Losses      int                   `json:"losses"`
	Winrate     float32               `json:"winrate"`
	Diagnostics simulationDiagnostics `json:"diagnostics"`
}

type simulationDifference struct {
	Metric string      `json:"metric"`
	API    interface{} `json:"api"`
	Local  interface{} `json:"local"`
}

var simulationCompareCmd = &cobra.Command{
	Use:   "compare <monster>",
	Short: "Compare the official and local simulators",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := readSimulationScalarFlags(cmd)
		if err != nil {
			return err
		}
		err = validateSimulationScalarFlags(flags)
		if err != nil {
			return err
		}
		request, err := simulationRequest(cmd, args[0], flags)
		if err != nil {
			return err
		}
		var apiResults []schemas.CombatResultSchema
		var lastRequest time.Time
		for _, iterations := range simulation.IterationChunks(flags.Iterations) {
			request.Iterations = iterations
			if !lastRequest.IsZero() {
				wait := time.Second - time.Since(lastRequest)
				if wait > 0 {
					time.Sleep(wait)
				}
			}
			lastRequest = time.Now()
			remote, requestErr := api.SimulationFight(request)
			if requestErr != nil {
				return requestErr
			}
			apiResults = append(apiResults, remote.Results...)
		}
		monster, ok := database.Monsters.Get(args[0])
		if !ok {
			return fmt.Errorf("invalid monster: %s", args[0])
		}
		c := request.Characters[0]
		localIterations := flags.Iterations
		localOptions := simulationCriticalOptions(flags)
		localOptions.Iterations = localIterations
		localOptions.RNG = rand.New(rand.NewSource(time.Now().UnixNano())).Float64
		local := fight.SimulateMany(fight.FromLoadout(c.Level, simulation.CharacterSlots(c), simulation.CharacterUtilities(c)), *monster, localOptions)
		comparison := simulationComparison{
			TotalIterations: localIterations,
			API: simulationStats{
				Wins:   countSimulationWins(apiResults),
				Losses: len(apiResults) - countSimulationWins(apiResults),
			},
			Local: simulationStats{
				Wins:    local.Wins,
				Losses:  local.Losses,
				Winrate: local.Winrate,
			},
		}
		comparison.API.Winrate = float32(comparison.API.Wins) * 100 / float32(len(apiResults))
		apiAverageTurns, apiAverageFinalHP := simulationAverages(apiResults)
		localAverageTurns, localAverageFinalHP := simulationAverages(local.Results)
		apiFighter := fight.FromLoadout(c.Level, simulation.CharacterSlots(c), simulation.CharacterUtilities(c))
		apiMetrics := fight.Metrics(apiFighter, c.Level, *monster, apiResults)
		localMetrics := fight.Metrics(fight.FromLoadout(c.Level, simulation.CharacterSlots(c), simulation.CharacterUtilities(c)), c.Level, *monster, local.Results)
		comparison.API.Diagnostics = simulationDiagnosticsFor(apiResults)
		comparison.API.Diagnostics.AverageTurns, comparison.API.Diagnostics.AverageFinalHP = apiAverageTurns, apiAverageFinalHP
		comparison.API.Diagnostics.AverageFightCooldown, comparison.API.Diagnostics.XP, comparison.API.Diagnostics.XPPerCycle = apiMetrics.AverageFightCooldown, apiMetrics.XP, apiMetrics.XPPerCycle
		comparison.Local.Diagnostics = simulationDiagnosticsFor(local.Results)
		comparison.Local.Diagnostics.AverageTurns, comparison.Local.Diagnostics.AverageFinalHP = localAverageTurns, localAverageFinalHP
		comparison.Local.Diagnostics.AverageFightCooldown, comparison.Local.Diagnostics.XP, comparison.Local.Diagnostics.XPPerCycle = localMetrics.AverageFightCooldown, localMetrics.XP, localMetrics.XPPerCycle
		if flags.Logs && len(apiResults) > 0 && len(local.Results) > 0 {
			comparison.APILogs = apiResults[0].Logs
			comparison.LocalLogs = local.Results[0].Logs
		}
		if comparison.API.Wins != comparison.Local.Wins {
			comparison.Differences = append(comparison.Differences, simulationDifference{
				Metric: "wins",
				API:    comparison.API.Wins,
				Local:  comparison.Local.Wins,
			})
		}
		if comparison.API.Losses != comparison.Local.Losses {
			comparison.Differences = append(comparison.Differences, simulationDifference{
				Metric: "losses",
				API:    comparison.API.Losses,
				Local:  comparison.Local.Losses,
			})
		}
		return console.Auto(comparison)
	},
}

func simulationDiagnosticsFor(results []schemas.CombatResultSchema) simulationDiagnostics {
	var diagnostics simulationDiagnostics
	barrierAbsorption := regexp.MustCompile(`barrier absorbs ([0-9]+) damage`)
	greedLog := regexp.MustCompile(`Greed empowers (Character_1|the monster) \(\+([0-9]+)% damage, total \+([0-9]+)%\)`)
	attackLog := regexp.MustCompile(`^Turn ([0-9]+): (.+?) used `)
	for _, result := range results {
		playerGreedTotal, monsterGreedTotal := 0, 0
		countedCriticalAttack := ""
		for _, log := range result.Logs {
			player := strings.Contains(log, "Character_1")
			side := func(player bool) *simulationSideDiagnostics {
				if player {
					return &diagnostics.Player
				}
				return &diagnostics.Monster
			}
			if strings.Contains(log, "Critical strike") {
				player := strings.Contains(log, "Character_1 used")
				attackKey := ""
				match := attackLog.FindStringSubmatch(log)
				if len(match) == 3 {
					attackKey = match[1] + ":" + match[2]
				}
				// The API emits one log line per elemental hit, but critical
				// chance is rolled once for the attack. Count the first
				// critical line in the consecutive elemental group, not every
				// element.
				if attackKey == "" {
					attackKey = strconv.FormatBool(player) + ":" + log
				}
				if attackKey != countedCriticalAttack {
					side(player).CriticalHits++
					countedCriticalAttack = attackKey
				}
			}
			match := attackLog.FindStringSubmatch(log)
			if len(match) == 3 {
				if !strings.Contains(log, "Critical strike") {
					countedCriticalAttack = ""
				}
			}
			if strings.Contains(log, "applies a burn") {
				side(!strings.Contains(log, "The monster applies")).BurnApplications++
			}
			if strings.Contains(log, "suffers from burn") {
				side(player).BurnTicks++
			}
			if strings.Contains(log, "applies a poison") {
				side(!strings.Contains(log, "The monster applies")).PoisonApplications++
			}
			if strings.Contains(log, "suffers from poison") {
				side(player).PoisonTicks++
			}
			if strings.Contains(log, "resistance is corrupted") {
				diagnostics.Player.CorruptedApplications++
			}
			if strings.Contains(log, "Health Potion") || strings.Contains(log, "health_potion") {
				diagnostics.Player.HealingPotions++
			}
			if strings.Contains(log, "Healing effect") {
				side(player).HealingActivations++
			}
			if strings.Contains(log, "from lifesteal") {
				side(strings.Contains(log, "Character_1 heals")).LifestealActivations++
			}
			if strings.Contains(log, "Berserker Rage activates") {
				diagnostics.Monster.BerserkerActivations++
			}
			match = greedLog.FindStringSubmatch(log)
			if len(match) == 4 {
				increment, _ := strconv.Atoi(match[2])
				total, _ := strconv.Atoi(match[3])
				previous := &monsterGreedTotal
				if match[1] == "Character_1" {
					previous = &playerGreedTotal
				}
				if increment > 0 && total > *previous {
					activations := (total - *previous) / increment
					side(match[1] == "Character_1").GreedActivations += activations
					*previous = total
				}
			}
			if strings.Contains(log, "shell activates") {
				diagnostics.Player.ShellActivations++
			}
			if strings.Contains(log, "shell effect has worn off") {
				diagnostics.Player.ShellExpirations++
			}
			if strings.Contains(log, "Enchanted Mirror activates") {
				diagnostics.Player.MirrorActivations++
			}
			if strings.Contains(log, "Frenzy triggers on critical") {
				side(strings.Contains(log, "Character_1")).FrenzyActivations++
			}
			if strings.Contains(log, "raises a barrier") || strings.Contains(log, "refreshes its barrier") {
				diagnostics.Monster.BarrierActivations++
			}
			if strings.Contains(log, "barrier is destroyed") {
				diagnostics.Monster.BarrierDestructions++
			}
			match = barrierAbsorption.FindStringSubmatch(log)
			if len(match) == 2 {
				absorbed, err := strconv.Atoi(match[1])
				if err == nil {
					diagnostics.Monster.BarrierAbsorbedDamage += absorbed
				}
			}
			if strings.Contains(log, "Antidote") || strings.Contains(log, "antidote") {
				diagnostics.Player.Antidotes++
			}
		}
	}
	return diagnostics
}

func countSimulationWins(results []schemas.CombatResultSchema) int {
	n := 0
	for _, r := range results {
		if r.Result == "win" {
			n++
		}
	}
	return n
}

func finalHP(r schemas.CombatResultSchema) float32 {
	if len(r.CharacterResults) == 0 {
		return 0
	}
	switch value := r.CharacterResults[0]["final_hp"].(type) {
	case float64:
		return float32(value)
	case int:
		return float32(value)
	default:
		return 0
	}
}

func simulationAverages(results []schemas.CombatResultSchema) (float32, float32) {
	var turns, hp float32
	for _, result := range results {
		turns += float32(result.Turns)
		hp += finalHP(result)
	}
	if len(results) == 0 {
		return 0, 0
	}
	return turns / float32(len(results)), hp / float32(len(results))
}

func init() {
	simulationCmd.AddCommand(simulationCompareCmd)
}
