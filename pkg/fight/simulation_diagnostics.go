package fight

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/schemas"
)

type SimulationDiagnostics struct {
	AverageTurns         float32                   `json:"average_turns"`
	AverageFinalHP       float32                   `json:"average_final_hp"`
	AverageFightCooldown float32                   `json:"average_fight_cooldown"`
	XP                   int                       `json:"xp"`
	XPPerCycle           float32                   `json:"xp_per_cycle"`
	Player               SimulationSideDiagnostics `json:"player"`
	Monster              SimulationSideDiagnostics `json:"monster"`
}

type SimulationSideDiagnostics struct {
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

func SimulationDiagnosticsFor(results []schemas.CombatResultSchema) SimulationDiagnostics {
	var diagnostics SimulationDiagnostics
	barrierAbsorption := regexp.MustCompile(`barrier absorbs ([0-9]+) damage`)
	greedLog := regexp.MustCompile(`Greed empowers (Character_1|the monster) \(\+([0-9]+)% damage, total \+([0-9]+)%\)`)
	attackLog := regexp.MustCompile(`^Turn ([0-9]+): (.+?) used `)
	for _, result := range results {
		playerGreedTotal, monsterGreedTotal := 0, 0
		countedCriticalAttack := ""
		for _, log := range result.Logs {
			player := strings.Contains(log, "Character_1")
			side := func(player bool) *SimulationSideDiagnostics {
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
				if attackKey == "" {
					attackKey = strconv.FormatBool(player) + ":" + log
				}
				if attackKey != countedCriticalAttack {
					side(player).CriticalHits++
					countedCriticalAttack = attackKey
				}
			}
			match := attackLog.FindStringSubmatch(log)
			if len(match) == 3 && !strings.Contains(log, "Critical strike") {
				countedCriticalAttack = ""
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

func CriticalSequenceFromLogs(logs []string) []bool {
	attackLog := regexp.MustCompile(`^Turn ([0-9]+): (.+?) used .* attack`)
	sequence := make([]bool, 0)
	lastKey := ""
	for _, log := range logs {
		match := attackLog.FindStringSubmatch(log)
		if match == nil {
			continue
		}
		key := match[1] + "\x00" + match[2]
		if key == lastKey {
			continue
		}
		lastKey = key
		sequence = append(sequence, strings.Contains(log, "Critical strike"))
	}
	return sequence
}
