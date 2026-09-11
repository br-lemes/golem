package fight

import (
	"reflect"
	"strings"
	"testing"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestApplyItemCoversEquipmentEffects(t *testing.T) {
	fields := map[string]string{
		"hp":               "HP",
		"attack_fire":      "AttackFire",
		"attack_earth":     "AttackEarth",
		"attack_water":     "AttackWater",
		"attack_air":       "AttackAir",
		"dmg":              "Dmg",
		"dmg_fire":         "DmgFire",
		"dmg_earth":        "DmgEarth",
		"dmg_water":        "DmgWater",
		"dmg_air":          "DmgAir",
		"res_fire":         "ResFire",
		"res_earth":        "ResEarth",
		"res_water":        "ResWater",
		"res_air":          "ResAir",
		"critical_strike":  "CriticalStrike",
		"initiative":       "Initiative",
		"haste":            "Haste",
		"wisdom":           "Wisdom",
		"burn":             "Burn",
		"lifesteal":        "Lifesteal",
		"healing":          "Healing",
		"shell":            "Shell",
		"enchanted_mirror": "EnchantedMirror",
		"frenzy":           "Frenzy",
		"greed":            "Greed",
	}
	found := make(map[string]bool)
	for _, item := range catalog.Items().All() {
		if item.Effects == nil {
			continue
		}
		for _, effect := range *item.Effects {
			field, ok := fields[effect.Code]
			if !ok || found[effect.Code] || effect.Value == 0 {
				continue
			}
			stats := Stats{}
			applyItem(&stats, item.Code, 1)
			value := reflect.ValueOf(stats).FieldByName(field).Int()
			if value != int64(effect.Value) {
				t.Fatalf("applyItem(%q) %s = %d, want %d", item.Code, effect.Code, value, effect.Value)
			}
			found[effect.Code] = true
		}
	}
	for code := range fields {
		if !found[code] {
			t.Errorf("no catalog item exercises %q", code)
		}
	}
}

func TestApplyItemIgnoresUnknownItemAndSupportsNegativeSign(t *testing.T) {
	stats := Stats{HP: 10}
	applyItem(&stats, "", 1)
	if stats.HP != 10 {
		t.Fatalf("empty item changed stats: %+v", stats)
	}
	applyItem(&stats, "does_not_exist", 1)
	if stats.HP != 10 {
		t.Fatalf("unknown item changed stats: %+v", stats)
	}
	var itemCode string
	var hpValue int
	for _, item := range catalog.Items().All() {
		if item.Effects == nil {
			continue
		}
		for _, effect := range *item.Effects {
			if effect.Code == "hp" && effect.Value != 0 {
				itemCode, hpValue = item.Code, effect.Value
				break
			}
		}
		if itemCode != "" {
			break
		}
	}
	if itemCode == "" {
		t.Fatal("no HP item in catalog")
	}
	applyItem(&stats, itemCode, -1)
	if stats.HP != 10-hpValue {
		t.Fatalf("negative sign was not applied: %+v", stats)
	}
}

func TestApplyCombatBoostAppliesAllSupportedEffects(t *testing.T) {
	codes := []string{
		"boost_hp",
		"boost_dmg_fire",
		"boost_dmg_earth",
		"boost_dmg_water",
		"boost_dmg_air",
		"boost_res_fire",
		"boost_res_earth",
		"boost_res_water",
		"boost_res_air",
	}
	for _, code := range codes {
		stats := Stats{}
		applyCombatBoost(&stats, code, 7)
		got := stats.HP + stats.DmgFire + stats.DmgEarth + stats.DmgWater + stats.DmgAir
		got += stats.ResFire + stats.ResEarth + stats.ResWater + stats.ResAir
		if got != 7 {
			t.Errorf("applyCombatBoost(%q) changed stats by %d, want 7", code, got)
		}
	}
}

func TestApplyCombatBoostIgnoresUnknownEffect(t *testing.T) {
	stats := Stats{HP: 10}
	applyCombatBoost(&stats, "unknown", 7)
	if stats.HP != 10 {
		t.Fatalf("unknown boost changed stats: %+v", stats)
	}
}

func TestApplyUtilitiesAppliesBoostsAndCollectsConsumables(t *testing.T) {
	result := applyUtilities(Fighter{Stats: Stats{HP: 100}}, map[string]int{
		"health_boost_potion": 2,
		"minor_health_potion": 3,
		"small_antidote":      4,
	})
	if result.Stats.HP <= 100 {
		t.Fatalf("HP boost was not applied: %+v", result.Stats)
	}
	if len(result.Utilities) != 2 {
		t.Fatalf("utilities = %+v, want two consumables", result.Utilities)
	}
	for _, utility := range result.Utilities {
		if utility.Quantity <= 0 || (utility.Restore == 0 && utility.Antipoison == 0) {
			t.Errorf("invalid utility = %+v", utility)
		}
	}
}

func TestApplyUtilitiesIgnoresInvalidAndNonPositiveQuantities(t *testing.T) {
	result := applyUtilities(Fighter{Stats: Stats{HP: 100}}, map[string]int{
		"minor_health_potion": 0,
		"small_antidote":      -1,
		"unknown_item":        5,
	})
	if result.Stats.HP != 100 || len(result.Utilities) != 0 {
		t.Fatalf("invalid utilities changed fighter: %+v", result)
	}
}

func TestReduceResistanceReducesSelectedElement(t *testing.T) {
	for _, element := range []string{"fire", "earth", "water", "air"} {
		stats := Stats{ResFire: 50, ResEarth: 50, ResWater: 50, ResAir: 50}
		got := reduceResistance(&stats, element, 12)
		if got != 38 {
			t.Errorf("reduceResistance(%q) = %d, want 38", element, got)
		}
		if stats.ResFire+stats.ResEarth+stats.ResWater+stats.ResAir != 188 {
			t.Errorf("reduceResistance(%q) changed an unexpected resistance: %+v", element, stats)
		}
	}
}

func TestReduceResistanceIgnoresUnknownElement(t *testing.T) {
	stats := Stats{ResFire: 50}
	got := reduceResistance(&stats, "unknown", 12)
	if got != 0 || stats.ResFire != 50 {
		t.Fatalf("unknown element changed resistance: got %d, stats %+v", got, stats)
	}
}

func TestRandomBubbleElement(t *testing.T) {
	got := randomBubbleElement(nil, -1)
	if got != 0 {
		t.Fatalf("nil RNG element = %d, want 0", got)
	}
	got = randomBubbleElement(func() float64 { return 0.50 }, -1)
	if got != 2 {
		t.Fatalf("random element = %d, want 2", got)
	}
	got = randomBubbleElement(func() float64 { return 0.50 }, 2)
	if got != 3 {
		t.Fatalf("repeated random element = %d, want 3", got)
	}
	got = randomBubbleElement(func() float64 { return 1 }, -1)
	if got != 3 {
		t.Fatalf("upper-bound random element = %d, want 3", got)
	}
}

func TestElementalDamage(t *testing.T) {
	got := elemental(100, 30, 30)
	if got != 91 {
		t.Fatalf("elemental damage = %d, want 91", got)
	}
}

func TestFromLoadoutNormalizesInvalidLevel(t *testing.T) {
	fighter := FromLoadout(0, nil, nil)
	if fighter.Stats.HP != 120 || fighter.Stats.Initiative != 100 {
		t.Fatalf("normalized fighter stats = %+v, want level-one stats", fighter.Stats)
	}
}

func TestSimulateManyNormalizesInvalidIterations(t *testing.T) {
	fighter := Fighter{Stats: Stats{HP: 100, AttackFire: 100, Initiative: 10}}
	monster := schemas.MonsterSchema{Hp: 50, Initiative: 1}
	summary := SimulateMany(fighter, monster, SimulationOptions{Iterations: 0})
	if len(summary.Results) != 1 || summary.Wins != 1 || summary.Losses != 0 {
		t.Fatalf("summary = %+v, want one win", summary)
	}
}

func TestSimulateLogsInitialMonsterBarrier(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "barrier", Value: 10}}
	player := Fighter{Stats: Stats{HP: 100, AttackFire: 1, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:    "Shielded dummy",
		Hp:      100,
		Effects: effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "raises a barrier of 10 HP") {
		t.Fatalf("initial barrier was not logged: %v", result.Logs)
	}
}

func TestSimulateLogsMonsterGreedAwakening(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "greed", Value: 10}}
	player := Fighter{Stats: Stats{HP: 1000, AttackFire: 20, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:        "Greedy dummy",
		Hp:          100,
		AttackEarth: 1,
		Effects:     effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "Greed awakens. The monster") {
		t.Fatalf("monster Greed was not logged: %v", result.Logs)
	}
}

func TestIsBossRecognizesBossTypes(t *testing.T) {
	for _, monsterType := range []string{"boss", "raid_boss"} {
		if !isBoss(schemas.MonsterSchema{Type: monsterType}) {
			t.Errorf("isBoss(%q) = false, want true", monsterType)
		}
	}
	if isBoss(schemas.MonsterSchema{Type: "normal"}) {
		t.Fatal("normal monster was recognized as boss")
	}
}

func TestSimulateProtectiveBubbleChangesElementEachTurn(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{
		{Code: "protective_bubble", Value: 20},
	}
	player := Fighter{Stats: Stats{HP: 1000, AttackFire: 1, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:        "Bubble dummy",
		Hp:          1000,
		AttackEarth: 1,
		Effects:     effects,
	}
	result := Simulate(player, monster, SimulationOptions{
		RNG: func() float64 { return 0 },
	})
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "protective bubble grants 20% fire resistance") {
		t.Fatalf("bubble fire log missing: %s", logs)
	}
	if !strings.Contains(logs, "protective bubble grants 20% earth resistance") {
		t.Fatalf("bubble logs do not show changing elements: %s", logs)
	}
}

func TestSimulateShellActivatesAndExpiresForBoss(t *testing.T) {
	player := Fighter{
		Stats: Stats{HP: 100, AttackEarth: 1, Initiative: 10, Shell: 20},
	}
	monster := schemas.MonsterSchema{
		Name:       "Boss dummy",
		Type:       "boss",
		Hp:         1000,
		AttackFire: 10,
		Initiative: 1,
	}
	result := Simulate(player, monster, SimulationOptions{})
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "shell activates") {
		t.Fatalf("shell activation missing: %s", logs)
	}
	if !strings.Contains(logs, "shell effect has worn off") {
		t.Fatalf("shell activation/expiration missing: %s", logs)
	}
}

func TestSimulateHealingActivatesOnThirdPlayerTurn(t *testing.T) {
	player := Fighter{
		Stats: Stats{HP: 100, AttackEarth: 1, Initiative: 10, Healing: 20},
	}
	monster := schemas.MonsterSchema{
		Name:       "Healing dummy",
		Hp:         1000,
		AttackFire: 5,
		Initiative: 1,
	}
	result := Simulate(player, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "Healing effect") {
		t.Fatalf("player Healing activation missing: %v", result.Logs)
	}
}

func TestSimulateAntidoteRemovesAllPoison(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "poison", Value: 20}}
	player := Fighter{
		Stats:     Stats{HP: 1000, AttackEarth: 1, Initiative: 1},
		Utilities: []Utility{{Code: "antidote", Antipoison: 20, Quantity: 1}},
	}
	monster := schemas.MonsterSchema{
		Name:       "Poison dummy",
		Hp:         1000,
		AttackFire: 1,
		Initiative: 10,
		Effects:    effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "removed 20 poison") {
		t.Fatalf("antidote was not used: %s", logs)
	}
}

func TestSimulateMonsterHealingActivatesOnThirdMonsterTurn(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "healing", Value: 20}}
	player := Fighter{Stats: Stats{HP: 1000, AttackEarth: 10, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:       "Healing monster",
		Hp:         100,
		AttackFire: 1,
		Initiative: 1,
		Effects:    effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "monster heals") {
		t.Fatalf("monster Healing activation missing: %v", result.Logs)
	}
}

func TestSimulateRefreshesMonsterBarrierEveryFiveMonsterTurns(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "barrier", Value: 10}}
	player := Fighter{Stats: Stats{HP: 1000, AttackEarth: 1, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:       "Barrier monster",
		Hp:         1000,
		AttackFire: 1,
		Initiative: 1,
		Effects:    effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "refreshes its barrier to 10 HP") {
		t.Fatalf("barrier refresh missing: %v", result.Logs)
	}
}

func TestSimulatePlayerBurnDamagesMonster(t *testing.T) {
	player := Fighter{
		Stats: Stats{HP: 1000, AttackFire: 10, Burn: 20, Initiative: 10},
	}
	monster := schemas.MonsterSchema{
		Name:       "Burn dummy",
		Hp:         1000,
		Initiative: 1,
	}
	result := Simulate(player, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "Monster suffers") {
		t.Fatalf("player burn damage missing: %v", result.Logs)
	}
}

func TestSimulateBarrierPartiallyAbsorbsDamageAndIsDestroyed(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "barrier", Value: 10}}
	player := Fighter{Stats: Stats{HP: 1000, AttackFire: 15, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:    "Barrier dummy",
		Hp:      100,
		Effects: effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	logs := strings.Join(result.Logs, "\n")
	if !strings.Contains(logs, "absorbed") || !strings.Contains(logs, "barrier is destroyed") {
		t.Fatalf("barrier damage logs missing: %s", logs)
	}
}

func TestSimulateCorruptsMonsterResistance(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "corrupted", Value: 5}}
	player := Fighter{Stats: Stats{HP: 1000, AttackFire: 10, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:    "Corrupted dummy",
		Hp:      100,
		ResFire: 20,
		Effects: effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "resistance is corrupted") {
		t.Fatalf("corruption log missing: %v", result.Logs)
	}
}

func TestSimulatePlayerLifestealActivatesOnCritical(t *testing.T) {
	chance := 100
	player := Fighter{
		Stats: Stats{HP: 100, AttackFire: 10, Initiative: 10, Lifesteal: 50},
	}
	monster := schemas.MonsterSchema{
		Name:       "Lifesteal dummy",
		Hp:         1000,
		Initiative: 1,
	}
	result := Simulate(player, monster, SimulationOptions{
		Critical: CriticalOptions{PlayerChance: &chance},
		RNG:      func() float64 { return 0 },
	})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "from lifesteal") {
		t.Fatalf("player lifesteal log missing: %v", result.Logs)
	}
}

func TestSimulateMonsterLifestealActivatesOnCritical(t *testing.T) {
	chance := 100
	effects := &[]schemas.SimpleEffectSchema{{Code: "lifesteal", Value: 50}}
	player := Fighter{Stats: Stats{HP: 1000, AttackEarth: 1, Initiative: 1}}
	monster := schemas.MonsterSchema{
		Name:       "Monster lifesteal",
		Hp:         1000,
		AttackFire: 10,
		Initiative: 10,
		Effects:    effects,
	}
	result := Simulate(player, monster, SimulationOptions{
		Critical: CriticalOptions{MonsterChance: &chance},
		RNG:      func() float64 { return 0 },
	})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "Monster heals") {
		t.Fatalf("monster lifesteal log missing: %v", result.Logs)
	}
}

func TestSimulateEnchantedMirrorReflectsMonsterDamage(t *testing.T) {
	player := Fighter{
		Stats: Stats{HP: 1000, Initiative: 1, EnchantedMirror: 50},
	}
	monster := schemas.MonsterSchema{
		Name:       "Mirror dummy",
		Hp:         1000,
		AttackFire: 10,
		Initiative: 10,
	}
	result := Simulate(player, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "Enchanted Mirror activates") {
		t.Fatalf("mirror log missing: %v", result.Logs)
	}
}

func TestSimulateBurnDefeatsMonsterDuringMonsterTurn(t *testing.T) {
	player := Fighter{
		Stats: Stats{HP: 1000, AttackFire: 100, Burn: 10, Initiative: 10},
	}
	monster := schemas.MonsterSchema{
		Name:       "Burn victim",
		Hp:         105,
		Initiative: 1,
	}
	result := Simulate(player, monster, SimulationOptions{})
	logs := strings.Join(result.Logs, "\n")
	if !result.Win || !strings.Contains(logs, "Monster has been defeated") {
		t.Fatalf("burn did not defeat monster: %+v", result)
	}
}

func TestSimulateBerserkerRageActivatesBelowQuarterHealth(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{
		{Code: "berserker_rage", Value: 50},
	}
	player := Fighter{Stats: Stats{HP: 1000, AttackFire: 80, Initiative: 10}}
	monster := schemas.MonsterSchema{
		Name:       "Berserker",
		Hp:         100,
		Initiative: 1,
		Effects:    effects,
	}
	result := Simulate(player, monster, SimulationOptions{})
	if !strings.Contains(strings.Join(result.Logs, "\n"), "Berserker Rage activates") {
		t.Fatalf("Berserker Rage activation missing: %v", result.Logs)
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

func TestConsumeUtilitySkipsHealingAboveHalfHealth(t *testing.T) {
	player := combatant{hp: 51, maxHP: 100}
	utilities := []Utility{{Code: "health_potion", Restore: 40, Quantity: 1}}
	logs := []string{}
	used := consumeUtility(&utilities, false, &player, 100, &logs, 1)
	if used || player.hp != 51 || utilities[0].Quantity != 1 || len(logs) != 0 {
		t.Fatalf("healing potion was consumed above half health: used=%v player=%+v utilities=%+v logs=%v", used, player, utilities, logs)
	}
}

func TestConsumeAntidoteReturnsConfiguredReduction(t *testing.T) {
	utilities := []Utility{{Code: "antidote", Antipoison: 50, Quantity: 1}}
	logs := []string{}
	got := consumeAntidote(&utilities, &logs, 2)
	if got != 50 {
		t.Fatalf("antidote reduction = %d, want 50", got)
	}
	if !strings.Contains(logs[0], "removed 50 poison") {
		t.Fatalf("antidote log = %q", logs[0])
	}
}

func TestConsumeAntidoteSkipsUnavailableUtilities(t *testing.T) {
	utilities := []Utility{
		{Code: "empty_antidote", Antipoison: 50, Quantity: 0},
		{Code: "health_potion", Restore: 20, Quantity: 1},
		{Code: "antidote", Antipoison: 30, Quantity: 1},
	}
	logs := []string{}
	got := consumeAntidote(&utilities, &logs, 3)
	if got != 30 {
		t.Fatalf("antidote reduction = %d, want 30", got)
	}
	if utilities[2].Quantity != 0 || len(logs) != 1 {
		t.Fatalf("utilities/logs after antidote = %+v, %v", utilities, logs)
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

func TestPlayerGreedAwakensBeforeEvaluatingDamageThresholds(t *testing.T) {
	effects := &[]schemas.SimpleEffectSchema{{Code: "burn", Value: 20}}
	fighter := Fighter{
		Stats: Stats{HP: 100, AttackEarth: 1, Greed: 10, Initiative: 1},
	}
	monster := schemas.MonsterSchema{
		Name:       "Flameche",
		Hp:         1000,
		AttackFire: 20,
		Initiative: 10,
		Effects:    effects,
	}

	zeroPlayerCritical, zeroMonsterCritical := 0, 0
	result := Simulate(fighter, monster, SimulationOptions{
		Critical: CriticalOptions{
			PlayerChance:  &zeroPlayerCritical,
			MonsterChance: &zeroMonsterCritical,
		},
	})
	joined := strings.Join(result.Logs, "\n")
	if strings.Contains(joined, "Turn 1: Greed empowers Character_1") {
		t.Fatalf("Greed was applied during its awakening turn: %s", joined)
	}
	want := "Turn 3: Greed empowers Character_1 (+10% damage, total +40%)."
	if !strings.Contains(joined, want) {
		t.Fatalf("Greed activation log missing; got logs: %s", joined)
	}
}

func TestGreedFlamecheReproduction(t *testing.T) {
	fighter := FromLoadout(50, map[string]string{
		"rune":       "powerful_rune",
		"shield":     "fire_shield",
		"helmet":     "obsidian_helmet",
		"body_armor": "medic_armor",
		"leg_armor":  "enchanter_pants",
		"boots":      "adamantite_boots",
		"ring1":      "mithril_ring",
		"ring2":      "mithril_ring",
		"amulet":     "heart_amulet",
		"artifact1":  "life_crystal",
		"artifact2":  "life_crystal",
		"artifact3":  "sandwhisper_codex",
	}, nil)
	effects := &[]schemas.SimpleEffectSchema{{Code: "burn", Value: 20}}
	monster := schemas.MonsterSchema{
		Name:       "Flameche",
		Hp:         2000,
		AttackFire: 1250,
		Initiative: 2000,
		ResAir:     -50,
		ResFire:    0,
		ResWater:   -50,
		Effects:    effects,
	}
	result := Simulate(fighter, monster, SimulationOptions{
		CriticalSequence: []bool{false, false, false, true, false},
	})
	joined := strings.Join(result.Logs, "\n")
	for _, want := range []string{
		"Fight start: Character_1 HP: 2565/2565",
		"Flameche used fire attack against Character_1 and dealt 563 damage",
		"Character_1 suffers from burn and loses 250 HP. Character_1 HP: 1752/2565",
		"Turn 3: Greed empowers Character_1 (+15% damage, total +75%).",
		"Turn 5: Greed empowers Character_1 (+15% damage, total +120%).",
		"Turn 7: Greed empowers Character_1 (+15% damage, total +150%).",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("reproduction log missing %q; got logs: %s", want, joined)
		}
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
