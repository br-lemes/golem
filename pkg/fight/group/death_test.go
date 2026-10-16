package group

import (
	"strings"
	"testing"
)

func TestDefeatedCharacterLeavesOrderAndTargetCandidates(t *testing.T) {
	fighters := []Fighter{
		{Stats: Stats{HP: 100, Threat: 100, Initiative: 100}},
		{Stats: Stats{HP: 100, Threat: 1, Initiative: 90}},
	}
	combat := NewCombat(fighters, Monster{Stats: Stats{Initiative: 200}}, func() float64 { return 0.99 })
	combat.Players[0].HP = 0
	combat.Players[0].Alive = false

	order := combat.ActorOrder()
	for _, actor := range order {
		if actor == 0 {
			t.Fatal("defeated character remained in actor order")
		}
	}
	got := combat.SelectTarget()
	if got != 1 {
		t.Fatalf("selected defeated character: got %d, want 1", got)
	}
}

func TestDefeatedCharacterKeepsItsScheduledTurn(t *testing.T) {
	fighters := []Fighter{
		{
			Name:  "survivor",
			Stats: Stats{HP: 100, AttackEarth: 1, Initiative: 200},
		},
		{
			Name:  "victim",
			Stats: Stats{HP: 10, Initiative: 150},
		},
	}
	monster := Monster{
		Name:  "wolf",
		Stats: Stats{HP: 1000, AttackFire: 100, Initiative: 100},
		HP:    1000,
		MaxHP: 1000,
	}
	combat := NewCombat(fighters, monster, func() float64 { return 0.99 })
	result := combat.Run(4)
	joined := strings.Join(result.Logs, "\n")
	if !strings.Contains(joined, "Turn 6: wolf used fire attack") {
		t.Fatalf("monster did not preserve the defeated slot: %s", joined)
	}
	if strings.Contains(joined, "Turn 5: wolf used fire attack") {
		t.Fatalf("monster acted during the defeated character's slot: %s", joined)
	}
}
