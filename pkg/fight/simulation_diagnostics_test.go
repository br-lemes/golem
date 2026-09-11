package fight

import (
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSimulationDiagnosticsForCountsCombatEvents(t *testing.T) {
	logs := []string{
		"Turn 1: Character_1 used fire attack and dealt 10 damage (Critical strike)",
		"Turn 1: Character_1 used earth attack and dealt 10 damage (Critical strike)",
		"Turn 1: The monster used fire attack against Character_1 and dealt 5 damage (Critical strike)",
		"Turn 1: Character_1 applies a burn of 4 on the monster.",
		"Turn 1: The monster applies a burn of 3 on Character_1.",
		"Turn 2: Character_1 suffers from burn and loses 3 HP.",
		"Turn 2: Character_1 suffers from poison and loses 2 HP.",
		"Turn 2: The monster applies a poison of 2 on Character_1.",
		"Turn 2: The monster's fire resistance is corrupted and decreases by 5%.",
		"Turn 2: Character_1 used health_potion and restored 20 HP.",
		"Turn 3: Character_1 heals 10 HP from Healing effect.",
		"Turn 3: Character_1 heals 2 HP from lifesteal.",
		"Turn 3: Monster heals 2 HP from lifesteal.",
		"Turn 3: The monster's Berserker Rage activates, gaining 20% permanent damage!",
		"Turn 3: Greed empowers Character_1 (+10% damage, total +10%).",
		"Turn 3: Greed empowers the monster (+5% damage, total +5%).",
		"Turn 3: Character_1's shell activates, gaining 10% resistance to all elements for 3 turns!",
		"Turn 6: Character_1's shell effect has worn off.",
		"Turn 4: Character_1's Enchanted Mirror activates, dealing back 2 damage to Dummy.",
		"Turn 4: Character_1's Frenzy triggers on critical.",
		"Turn 4: The monster raises a barrier of 10 HP.",
		"Turn 5: The monster refreshes its barrier to 10 HP.",
		"Turn 5: The monster's barrier absorbs 7 damage.",
		"Turn 5: The monster's barrier is destroyed!",
		"Turn 5: Character_1 used small_antidote and removed 2 poison.",
	}
	got := SimulationDiagnosticsFor([]schemas.CombatResultSchema{{Logs: logs}})
	if got.Player.CriticalHits != 1 || got.Monster.CriticalHits != 1 {
		t.Fatalf("critical hits = %+v/%+v", got.Player, got.Monster)
	}
	if got.Player.BurnApplications != 1 || got.Monster.BurnApplications != 1 || got.Player.BurnTicks != 1 {
		t.Fatalf("burn diagnostics = %+v/%+v", got.Player, got.Monster)
	}
	if got.Player.PoisonApplications != 0 || got.Monster.PoisonApplications != 1 || got.Player.PoisonTicks != 1 {
		t.Fatalf("poison diagnostics = %+v/%+v", got.Player, got.Monster)
	}
	if got.Player.CorruptedApplications != 1 || got.Player.HealingPotions != 1 || got.Player.HealingActivations != 1 {
		t.Fatalf("player effect diagnostics = %+v", got.Player)
	}
	if got.Player.LifestealActivations != 1 || got.Monster.LifestealActivations != 1 {
		t.Fatalf("lifesteal diagnostics = %+v/%+v", got.Player, got.Monster)
	}
	if got.Monster.BerserkerActivations != 1 || got.Player.GreedActivations != 1 || got.Monster.GreedActivations != 1 {
		t.Fatalf("power diagnostics = %+v/%+v", got.Player, got.Monster)
	}
	if got.Player.ShellActivations != 1 || got.Player.ShellExpirations != 1 || got.Player.MirrorActivations != 1 || got.Player.FrenzyActivations != 1 {
		t.Fatalf("player activation diagnostics = %+v", got.Player)
	}
	if got.Monster.BarrierActivations != 2 || got.Monster.BarrierDestructions != 1 || got.Monster.BarrierAbsorbedDamage != 7 || got.Player.Antidotes != 1 {
		t.Fatalf("barrier/antidote diagnostics = %+v/%+v", got.Player, got.Monster)
	}
}

func TestCriticalSequenceFromLogsDeduplicatesElementalHits(t *testing.T) {
	logs := []string{
		"Turn 1: Character_1 used fire attack and dealt 10 damage (Critical strike)",
		"Turn 1: Character_1 used earth attack and dealt 10 damage (Critical strike)",
		"Turn 2: Dummy used fire attack against Character_1 and dealt 5 damage",
		"Turn 3: Character_1 used water attack and dealt 4 damage",
	}
	got := CriticalSequenceFromLogs(logs)
	want := []bool{true, false, false}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("critical sequence = %v, want %v", got, want)
	}
}

func TestSimulationDiagnosticsCountsMalformedCriticalLog(t *testing.T) {
	result := schemas.CombatResultSchema{
		Logs: []string{
			"Character_1 used an unrecognized move (Critical strike)",
		},
	}
	got := SimulationDiagnosticsFor([]schemas.CombatResultSchema{result})
	if got.Player.CriticalHits != 1 {
		t.Fatalf("player critical hits = %d, want 1", got.Player.CriticalHits)
	}
}

func TestCriticalSequenceFromLogsIgnoresNonAttackLogs(t *testing.T) {
	logs := []string{
		"Fight start: Character_1 HP: 100/100",
		"Turn 1: Character_1 applies a burn of 2 on the monster.",
		"Turn 1: Character_1 used fire attack and dealt 10 damage",
	}
	got := CriticalSequenceFromLogs(logs)
	if !reflect.DeepEqual(got, []bool{false}) {
		t.Fatalf("critical sequence = %v, want [false]", got)
	}
}
