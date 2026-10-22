package group

import "testing"

func TestSelectTargetThreatAndLowHP(t *testing.T) {
	combat := NewCombat([]Fighter{
		{Name: "tank", Stats: Stats{HP: 100, Threat: 20}},
		{Name: "damage", Stats: Stats{HP: 100, Threat: 10}},
	}, Monster{}, func() float64 { return 0.5 })
	got := combat.SelectTarget()
	if got != 0 {
		t.Fatalf("target = %d, want highest threat", got)
	}
	combat.Players[0].HP = 40
	combat.RNG = func() float64 { return 0.05 }
	got = combat.SelectTarget()
	if got != 0 {
		t.Fatalf("target = %d, want lowest HP on threat fallback tie", got)
	}
}

func TestTargetTieBreakers(t *testing.T) {
	combat := NewCombat([]Fighter{
		{Name: "healthy", Stats: Stats{HP: 100, Threat: 20}},
		{Name: "wounded", Stats: Stats{HP: 100, Threat: 20}},
	}, Monster{}, func() float64 { return 0.99 })
	combat.Players[0].HP = 90
	combat.Players[1].HP = 40
	got := combat.SelectTarget()
	if got != 1 {
		t.Fatalf("threat tie target = %d, want lower HP character 1", got)
	}

	combat.RNG = func() float64 { return 0.05 }
	combat.Players[0].HP = 40
	got = combat.SelectTarget()
	if got != 0 {
		t.Fatalf("lowest HP tie target = %d, want first tied character 0", got)
	}
	combat.RNG = func() float64 { return 0.99 }
	got = combat.SelectTarget()
	if got != 0 {
		t.Fatalf("lowest HP tie target = %d, want first tied character 0", got)
	}
}

func TestInitiativeTieBreaksByHP(t *testing.T) {
	combat := NewCombat([]Fighter{
		{Name: "low", Stats: Stats{HP: 100, Initiative: 20}},
		{Name: "high", Stats: Stats{HP: 120, Initiative: 20}},
	}, Monster{}, nil)
	order := combat.TurnOrder()
	if order[0] != 1 {
		t.Fatalf("initiative tie first actor = %d, want higher HP character 1", order[0])
	}
}

func TestInitiativeAndHPTiePreservesPartyOrder(t *testing.T) {
	combat := NewCombat([]Fighter{
		{Name: "first", Stats: Stats{HP: 100, Initiative: 20}},
		{Name: "second", Stats: Stats{HP: 100, Initiative: 20}},
	}, Monster{}, nil)
	order := combat.TurnOrder()
	if order[0] != 0 || order[1] != 1 {
		t.Fatalf("tied initiative order = %v, want party order [0 1]", order)
	}
}

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
