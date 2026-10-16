package group

import (
	"strings"
	"testing"
)

func TestPlayerGreedAwakensBeforeEvaluatingDamageThresholds(t *testing.T) {
	fighter := Fighter{
		Stats: Stats{HP: 100, AttackEarth: 1, Greed: 10, Initiative: 1},
	}
	monster := Monster{
		Name:  "Flameche",
		Stats: Stats{HP: 1000, AttackFire: 41, Initiative: 10, Burn: 20},
		HP:    1000,
		MaxHP: 1000,
	}

	combat := NewCombat([]Fighter{fighter}, monster, func() float64 { return 0.99 })
	result := combat.Run(20)
	joined := strings.Join(result.Logs, "\n")
	if strings.Contains(joined, "Turn 2: Greed empowers Character_1") {
		t.Fatalf("Greed was applied during its awakening turn: %s", joined)
	}
	if !strings.Contains(joined, "Greed empowers Character_1") {
		t.Fatalf("Greed never activated: %s", joined)
	}
	if !strings.Contains(joined, "Turn 3: Greed empowers Character_1") {
		t.Fatalf("Greed was not evaluated after the first monster turn: %s", joined)
	}
	if !strings.Contains(joined, "total +90%") {
		t.Fatalf("Greed did not include burn damage in the threshold calculation: %s", joined)
	}
}

func TestMonsterGreedUsesBoundedHPThresholds(t *testing.T) {
	fighter := Fighter{Stats: Stats{HP: 100, AttackEarth: 41, Initiative: 1}}
	monster := Monster{
		Name:  "Flameche",
		Stats: Stats{HP: 100, AttackFire: 1, Initiative: 10, Greed: 10},
		HP:    100,
		MaxHP: 100,
	}

	combat := NewCombat([]Fighter{fighter}, monster, func() float64 { return 0.99 })
	result := combat.Run(10)
	joined := strings.Join(result.Logs, "\n")
	if !strings.Contains(joined, "Greed empowers the monster") {
		t.Fatalf("monster Greed never activated: %s", joined)
	}
	if !strings.Contains(joined, "total +50%") {
		t.Fatalf("monster Greed used an incorrect HP threshold: %s", joined)
	}
}

func TestGreedThresholdsUseIntegerTenPercentBlocks(t *testing.T) {
	got := greedThresholds(455, 275)
	if got != 4 {
		t.Fatalf("greed thresholds = %d, want 4", got)
	}
	got = greedThresholds(455, 274)
	if got != 4 {
		t.Fatalf("greed thresholds = %d, want 4", got)
	}
	got = greedThresholds(455, 273)
	if got != 4 {
		t.Fatalf("greed thresholds = %d, want 4", got)
	}
}
