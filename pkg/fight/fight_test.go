package fight

import (
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestElementalDamage(t *testing.T) {
	got := elemental(100, 30, 30)
	if got != 91 {
		t.Fatalf("elemental damage = %d, want 91", got)
	}
}

func TestSimulatePlayerActsFirstAndWins(t *testing.T) {
	fighter := Fighter{Stats: Stats{HP: 100, AttackFire: 100, Initiative: 10}}
	monster := schemas.MonsterSchema{Hp: 50, ResFire: 0, Initiative: 1}
	result := Simulate(fighter, monster, SimulationOptions{})
	if !result.Win || !result.PlayerFirst || result.Turns != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSimulateTimesOut(t *testing.T) {
	fighter := Fighter{Stats: Stats{HP: 100, AttackFire: 1}}
	monster := schemas.MonsterSchema{Hp: 1000, ResFire: 100}
	result := Simulate(fighter, monster, SimulationOptions{})
	if !result.TimedOut || result.Win {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSimulatePoisonIsAppliedAndTicks(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "poison", Value: 20}}
	fighter := Fighter{Stats: Stats{HP: 100, AttackFire: 1, Initiative: 1}}
	monster := schemas.MonsterSchema{
		Name:       "Spider",
		Hp:         1000,
		AttackFire: 1,
		Initiative: 10,
		Effects:    effects,
	}
	result := Simulate(fighter, monster, SimulationOptions{})
	if result.HPRemaining != 0 {
		t.Fatalf("poison damage was not applied: %+v", result)
	}
	poisonTicks := 0
	for _, log := range result.Logs {
		if strings.Contains(log, "suffers from poison") {
			poisonTicks++
		}
	}
	if poisonTicks != 5 {
		t.Fatalf("poison tick count = %d, want 5", poisonTicks)
	}
}

func TestConsumeUtilityRestoresHPAndRemovesPoison(t *testing.T) {
	player := combatant{hp: 20, maxHP: 100}
	utilities := []Utility{
		{Code: "minor_health_potion", Restore: 70, Quantity: 1},
		{Code: "small_antidote", Antipoison: 20, Quantity: 1},
	}
	logs := []string{}
	consumeUtility(&utilities, false, &player, 100, &logs, 1)
	if player.hp != 90 {
		t.Fatalf("heal utility failed: hp=%v utilities=%+v", player.hp, utilities)
	}
	if !consumeUtility(&utilities, true, &player, 100, &logs, 2) || utilities[1].Quantity != 0 {
		t.Fatalf("antidote utility failed: utilities=%+v", utilities)
	}
}

func TestSimulateBurnTicksAndDecays(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "burn", Value: 10}}
	fighter := Fighter{Stats: Stats{HP: 1000, Initiative: 1}}
	monster := schemas.MonsterSchema{
		Name:        "Imp",
		Hp:          1000,
		AttackEarth: 100,
		Initiative:  10,
		Effects:     effects,
	}
	result := Simulate(fighter, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "suffers from burn") {
		t.Fatalf("burn was not applied: %+v", result.Logs)
	}
	if result.HPRemaining >= 1000 {
		t.Fatalf("burn did not reduce HP: %+v", result)
	}
}

func TestDecayBurnDamageMatchesAPILogSequence(t *testing.T) {
	want := []int{22, 20, 18, 16, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	for i, value := range want[:len(want)-1] {
		got := int(decayBurnDamage(value))
		if got != want[i+1] {
			t.Fatalf("burn tick %d: got %d, want %d", i, got, want[i+1])
		}
	}
}

func TestPlayerFrenzyActivatesOnCriticalAndAffectsNextAttack(t *testing.T) {
	chance := 100
	fighter := Fighter{
		Stats: Stats{
			HP:             1000,
			AttackEarth:    10,
			CriticalStrike: 0,
			Initiative:     10,
			Frenzy:         12,
		},
	}
	monster := schemas.MonsterSchema{Hp: 1000, Initiative: 1}
	options := SimulationOptions{
		Critical: CriticalOptions{PlayerChance: &chance},
	}
	options.RNG = func() float64 { return 0 }
	result := simulate(fighter, monster, options)
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "Frenzy triggers on critical") || !strings.Contains(logs, "dealt 17 damage") {
		t.Fatalf("player Frenzy did not affect the next attack: %s", logs)
	}
}

func TestMonsterFrenzyActivatesOnCriticalAndAffectsNextAttack(t *testing.T) {
	chance := 100
	effects := &[]schemas.SimpleEffectSchema{{Code: "frenzy", Value: 12}}
	fighter := Fighter{Stats: Stats{HP: 1000, Initiative: 1}}
	monster := schemas.MonsterSchema{
		Name:           "Frenzy monster",
		Hp:             1000,
		AttackEarth:    10,
		CriticalStrike: 0,
		Initiative:     10,
		Effects:        effects,
	}
	options := SimulationOptions{
		Critical: CriticalOptions{MonsterChance: &chance},
	}
	options.RNG = func() float64 { return 0 }
	result := simulate(fighter, monster, options)
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "monster's Frenzy triggers on critical") || !strings.Contains(logs, "monster's Frenzy activates") {
		t.Fatalf("monster Frenzy did not activate: %s", logs)
	}
}

func TestMonsterReconstitutionActivatesEveryTwentyMonsterTurns(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{
		{Code: "reconstitution", Value: 20},
	}
	fighter := Fighter{Stats: Stats{HP: 1000, AttackEarth: 1, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:       "Reconstituting dummy",
		Hp:         100,
		AttackFire: 1,
		Initiative: 1,
		Effects:    effects,
	}
	result := Simulate(fighter, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "uses Reconstitution") {
		t.Fatalf("Reconstitution did not activate: %v", result.Logs)
	}
}

func TestMonsterVoidDrainDamagesPlayerAndHealsMonster(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "void_drain", Value: 10}}
	fighter := Fighter{Stats: Stats{HP: 1000, AttackEarth: 1, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:        "Void dummy",
		Hp:          1000,
		AttackWater: 1,
		Initiative:  1,
		Effects:     effects,
	}
	result := Simulate(fighter, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "uses Void Drain and drains 100 HP") {
		t.Fatalf("Void Drain did not activate: %v", result.Logs)
	}
}

func TestCombatBoostPotionsApplyOnceAtFightStart(t *testing.T) {
	fighter := FromLoadout(40, map[string]string{}, map[string]int{
		"health_boost_potion":   5,
		"enhanced_boost_potion": 3,
		"fire_res_potion":       2,
	})
	if fighter.Stats.HP != 120+39*5+250 {
		t.Fatalf("boost HP = %d, want %d", fighter.Stats.HP, 120+39*5+250)
	}
	if fighter.Stats.DmgFire != 20 || fighter.Stats.DmgEarth != 20 || fighter.Stats.DmgWater != 20 || fighter.Stats.DmgAir != 20 {
		t.Fatalf("boost damage = fire %d earth %d water %d air %d, want 20 each", fighter.Stats.DmgFire, fighter.Stats.DmgEarth, fighter.Stats.DmgWater, fighter.Stats.DmgAir)
	}
	if fighter.Stats.ResFire != 10 {
		t.Fatalf("boost fire resistance = %d, want 10", fighter.Stats.ResFire)
	}
}
