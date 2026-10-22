// Package group contains the independent multi-character fight simulator.
package group

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type RNG func() float64

type CriticalDecision struct {
	Actor    int
	Critical bool
}

type ActionDecision struct {
	Turn  int
	Actor int
}

type GuardDecision struct {
	Turn, Protector, Target int
}

type Stats struct {
	HP, Initiative, Threat                          int
	AttackFire, AttackEarth, AttackWater, AttackAir int
	Dmg, DmgFire, DmgEarth, DmgWater, DmgAir        int
	ResFire, ResEarth, ResWater, ResAir             int
	CriticalStrike, Haste, Wisdom                   int
	Burn, Lifesteal, Healing, Shell, Guard          int
	Frenzy, EnchantedMirror, Greed, VampiricStrike  int
	Poison, Corrupted                               int
	ProtectiveBubble                                int
	HealingAura                                     int
	SunShield                                       int
}

type Fighter struct {
	Name      string
	Stats     Stats
	Utilities []Utility
}

type Utility struct {
	Code          string
	Restore       int
	SplashRestore int
	Antipoison    int
	Quantity      int
}

type Player struct {
	Fighter       Fighter
	HP            float64
	MaxHP         float64
	Alive         bool
	ShellTurns    int
	ShellUsed     bool
	PlayedTurns   int
	BurnDamage    float64
	BurnApplied   bool
	Poison        int
	FrenzyBonus   int
	FrenzyUntil   int
	GreedStacks   int
	GreedAwakened bool
	VampiricLast  int
	GuardTarget   int
	GuardTurns    int
	GuardUses     int
}

type Monster struct {
	Name           string
	Stats          Stats
	HP             float64
	MaxHP          float64
	Type           string
	Barrier        float64
	BarrierAmount  float64
	Turns          int
	Reconstitution int
	VoidDrain      int
	BerserkerRage  int
	Berserked      bool
	FrenzyPending  bool
	GreedStacks    int
	BubbleElement  int
	BaseRes        [4]int
	BurnDamage     float64
}

type Combat struct {
	Players           []Player
	Monster           Monster
	RNG               RNG
	CriticalSequence  []bool
	CriticalDecisions []CriticalDecision
	TargetDecisions   []int
	ActionDecisions   []ActionDecision
	BubbleDecisions   []int
	GuardDecisions    []GuardDecision
	criticalIndex     int
	decisionUsed      []bool
	targetIndex       int
	guardIndex        int
	bubbleIndex       int
}

type Result struct {
	Win     bool
	Turns   int
	Logs    []string
	FinalHP []float64
}

func NewCombat(fighters []Fighter, monster Monster, rng RNG) Combat {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano())).Float64
	}
	players := make([]Player, len(fighters))
	for i, fighter := range fighters {
		fighter.Utilities = append([]Utility(nil), fighter.Utilities...)
		maxHP := fighter.Stats.HP
		if maxHP < 1 {
			maxHP = 1
		}
		players[i] = Player{
			Fighter:       fighter,
			HP:            float64(maxHP),
			MaxHP:         float64(maxHP),
			Alive:         true,
			GuardTarget:   -1,
			VampiricLast:  -3,
			GreedAwakened: fighter.Stats.Greed <= 0,
		}
	}
	return Combat{Players: players, Monster: monster, RNG: rng}
}

func (combat *Combat) Run(maxTurns int) Result {
	if maxTurns < 1 {
		maxTurns = 1
	}
	logs := []string{}
	if combat.Monster.Barrier > 0 {
		logs = append(logs, fmt.Sprintf("Fight start: The monster raises a barrier of %d HP.", int(combat.Monster.Barrier)))
	}
	logs = append(logs, combat.fightStartLog())
	for i, player := range combat.Players {
		if player.Fighter.Stats.Greed > 0 {
			logs = append(logs, fmt.Sprintf("Turn 1: Greed awakens. Character_%d gains +%d%% damage each time they lose 10%% max HP.", i+1, player.Fighter.Stats.Greed))
		}
	}
	if len(combat.ActionDecisions) > 0 {
		return combat.runReplay(maxTurns, logs)
	}
	turns := 0
	for combat.Monster.HP > 0 && combat.partyAlive() && turns < maxTurns {
		for _, index := range combat.ActorOrder() {
			if turns >= maxTurns || combat.Monster.HP <= 0 {
				break
			}
			if index < 0 {
				for _, player := range combat.Players {
					if !player.Alive || player.HP <= 0 {
						turns++
					}
				}
				if turns >= maxTurns {
					break
				}
				bubbleLog := combat.refreshBubble(turns + 1)
				if bubbleLog != "" {
					logs = append(logs, bubbleLog)
				}
				target := combat.SelectTarget()
				if target < 0 {
					continue
				}
				turns++
				logs = append(logs, combat.executeMonsterAttack(target, turns)...)
				continue
			}
			if !combat.Players[index].Alive || combat.Players[index].HP <= 0 {
				continue
			}
			logs = append(logs, combat.preparePlayerAction(index, turns+1)...)
			if !combat.Players[index].Alive {
				continue
			}
			turns++
			logs = append(logs, combat.executePlayerAttack(index, turns)...)
			logs = append(logs, combat.applyHealingAura(index, turns)...)
		}
		if combat.Monster.HP <= 0 {
			break
		}
	}
	finalHP := make([]float64, len(combat.Players))
	for i, player := range combat.Players {
		finalHP[i] = player.HP
	}
	win := combat.won(turns, maxTurns)
	logs = append(logs, combat.fightResultLog(win))
	return Result{Win: win, Turns: turns, Logs: logs, FinalHP: finalHP}
}

func (combat *Combat) runReplay(maxTurns int, logs []string) Result {
	turns := 0
	for _, decision := range combat.ActionDecisions {
		if decision.Turn > maxTurns {
			break
		}
		if decision.Turn > turns {
			turns = decision.Turn
		}
		if decision.Actor < 0 {
			bubbleLog := combat.refreshBubble(decision.Turn)
			if bubbleLog != "" {
				logs = append(logs, bubbleLog)
			}
			target := combat.SelectTarget()
			if target >= 0 {
				logs = append(logs, combat.executeMonsterAttack(target, decision.Turn)...)
			}
			if combat.Monster.HP <= 0 {
				break
			}
			continue
		}
		if decision.Actor >= len(combat.Players) || !combat.Players[decision.Actor].Alive {
			continue
		}
		logs = append(logs, combat.preparePlayerAction(decision.Actor, decision.Turn)...)
		if combat.Players[decision.Actor].Alive {
			logs = append(logs, combat.executePlayerAttack(decision.Actor, decision.Turn)...)
			logs = append(logs, combat.applyHealingAura(decision.Actor, decision.Turn)...)
		}
		if combat.Monster.HP <= 0 {
			break
		}
	}
	return combat.finishResult(turns, maxTurns, logs)
}

func (combat *Combat) finishResult(turns, maxTurns int, logs []string) Result {
	finalHP := make([]float64, len(combat.Players))
	for i, player := range combat.Players {
		finalHP[i] = player.HP
	}
	win := combat.won(turns, maxTurns)
	logs = append(logs, combat.fightResultLog(win))
	return Result{Win: win, Turns: turns, Logs: logs, FinalHP: finalHP}
}

func (combat *Combat) won(turns, maxTurns int) bool {
	if combat.Monster.HP <= 0 {
		return true
	}
	return combat.Monster.Type == "raid_boss" && turns >= maxTurns && combat.partyAlive()
}

func (combat *Combat) executePlayerAttack(index, turn int) []string {
	return combat.attackPlayer(index, turn)
}

func (combat *Combat) executeMonsterAttack(target, turn int) []string {
	return combat.attackMonster(target, turn)
}

func (combat *Combat) preparePlayerAction(index, turn int) []string {
	logs := []string{}
	combat.expireGuard(index)
	guardLog := combat.activateGuard(index, turn)
	if guardLog != "" {
		logs = append(logs, guardLog)
	}
	combat.Players[index].PlayedTurns++
	poisonLog := combat.applyPoison(index, turn)
	if poisonLog != "" {
		logs = append(logs, poisonLog)
	}
	utilityLog := combat.consumeRestore(index, turn)
	if utilityLog != "" {
		logs = append(logs, utilityLog)
	}
	splashLog := combat.consumeSplashRestore(index, turn)
	if splashLog != "" {
		logs = append(logs, splashLog)
	}
	logs = append(logs, combat.applyPlayerBurn(index, turn)...)
	if !combat.Players[index].Alive {
		return logs
	}
	heal := combat.healPlayer(index)
	if heal > 0 {
		logs = append(logs, fmt.Sprintf("Turn %d: Character_%d heals %d HP from Healing effect. HP: %d/%d", turn, index+1, heal, int(combat.Players[index].HP), int(combat.Players[index].MaxHP)))
	}
	if combat.activateShell(index) {
		logs = append(logs, fmt.Sprintf("Turn %d: Character_%d's shell activates, gaining %d%% resistance to all elements for 3 turns!", turn, index+1, combat.Players[index].Fighter.Stats.Shell))
	}
	return logs
}

func (combat *Combat) fightStartLog() string {
	players := make([]string, len(combat.Players))
	for i, player := range combat.Players {
		players[i] = fmt.Sprintf("Character_%d HP: %d/%d", i+1, int(player.HP), int(player.MaxHP))
	}
	return fmt.Sprintf("Fight start: %s vs %s HP: %d/%d", strings.Join(players, ", "), combat.Monster.Name, int(combat.Monster.HP), int(combat.Monster.MaxHP))
}

func (combat *Combat) fightResultLog(win bool) string {
	players := make([]string, len(combat.Players))
	for i, player := range combat.Players {
		players[i] = fmt.Sprintf("Character_%d HP: %d/%d", i+1, int(player.HP), int(player.MaxHP))
	}
	result := "loss"
	if win {
		result = "win"
	}
	return fmt.Sprintf("Fight result: %s. %s vs %s HP: %d/%d", result, strings.Join(players, ", "), combat.Monster.Name, int(combat.Monster.HP), int(combat.Monster.MaxHP))
}

func (combat *Combat) consumeRestore(index, turn int) string {
	player := &combat.Players[index]
	if player.HP > player.MaxHP/2 {
		return ""
	}
	for utility := range player.Fighter.Utilities {
		u := &player.Fighter.Utilities[utility]
		if u.Quantity <= 0 || u.Restore <= 0 {
			continue
		}
		u.Quantity--
		player.HP = math.Min(player.MaxHP, player.HP+float64(u.Restore))
		return fmt.Sprintf("Turn %d: Character_%d used %s and restored %d HP. (HP: %d/%d)", turn, index+1, utilityDisplayName(u.Code), u.Restore, int(player.HP), int(player.MaxHP))
	}
	return ""
}

func (combat *Combat) applyHealingAura(index, turn int) []string {
	player := &combat.Players[index]
	if player.PlayedTurns%2 != 0 || player.Fighter.Stats.HealingAura <= 0 {
		return nil
	}
	allies := []string{}
	for allyIndex := range combat.Players {
		if allyIndex == index || !combat.Players[allyIndex].Alive {
			continue
		}
		ally := &combat.Players[allyIndex]
		heal := int(math.Round(ally.MaxHP * float64(player.Fighter.Stats.HealingAura) / 100))
		before := ally.HP
		ally.HP = math.Min(ally.MaxHP, ally.HP+float64(heal))
		healed := int(math.Round(ally.HP - before))
		if healed > 0 {
			allies = append(allies, fmt.Sprintf("Character_%d +%d HP (%d/%d)", allyIndex+1, healed, int(ally.HP), int(ally.MaxHP)))
		}
	}
	if len(allies) == 0 {
		return nil
	}
	return []string{fmt.Sprintf("Turn %d: Character_%d's Healing Aura activates, healing all allies for %d%% of their max HP: %s", turn, index+1, player.Fighter.Stats.HealingAura, strings.Join(allies, ", "))}
}

func (combat *Combat) consumeSplashRestore(index, turn int) string {
	for utility := range combat.Players[index].Fighter.Utilities {
		u := &combat.Players[index].Fighter.Utilities[utility]
		if u.Quantity <= 0 || u.SplashRestore <= 0 {
			continue
		}
		target := -1
		lost := 0.0
		for candidate := range combat.Players {
			ally := &combat.Players[candidate]
			currentLost := ally.MaxHP - ally.HP
			if candidate != index && ally.Alive && ally.HP < ally.MaxHP*0.5 && currentLost > lost {
				target = candidate
				lost = currentLost
			}
		}
		if target < 0 {
			return ""
		}
		u.Quantity--
		ally := &combat.Players[target]
		before := ally.HP
		ally.HP = math.Min(ally.MaxHP, ally.HP+float64(u.SplashRestore))
		healed := int(math.Round(ally.HP - before))
		return fmt.Sprintf("Turn %d: Character_%d used %s and restored %d HP to Character_%d. (Character_%d HP: %d/%d)", turn, index+1, utilityDisplayName(u.Code), healed, target+1, target+1, int(ally.HP), int(ally.MaxHP))
	}
	return ""
}

func utilityDisplayName(code string) string {
	switch code {
	case "health_potion":
		return "Health Potion"
	case "health_splash_potion":
		return "Health Splash Potion"
	default:
		return code
	}
}

func (combat *Combat) applyPoison(index, turn int) string {
	player := &combat.Players[index]
	for utility := range player.Fighter.Utilities {
		u := &player.Fighter.Utilities[utility]
		if player.Poison > 0 && u.Quantity > 0 && u.Antipoison > 0 {
			removed := u.Antipoison
			if removed > player.Poison {
				removed = player.Poison
			}
			player.Poison -= removed
			u.Quantity--
			return fmt.Sprintf("Turn %d: Character_%d used %s and removed %d poison.", turn, index+1, u.Code, removed)
		}
	}
	if player.Poison <= 0 {
		if combat.Monster.Stats.Poison > 0 && combat.Monster.Turns > 0 {
			player.Poison = combat.Monster.Stats.Poison
			return fmt.Sprintf("Turn %d: The monster applies a poison of %d on Character_%d. (Character poison: %d)", turn, player.Poison, index+1, player.Poison)
		}
		return ""
	}
	damage := player.Poison
	player.HP = math.Max(0, player.HP-float64(damage))
	if player.HP == 0 {
		player.Alive = false
	}
	return fmt.Sprintf("Turn %d: Character_%d suffers from poison and loses %d HP. Character_%d HP: %d/%d", turn, index+1, damage, index+1, int(player.HP), int(player.MaxHP))
}

func (combat *Combat) healPlayer(index int) int {
	player := &combat.Players[index]
	if player.PlayedTurns%3 != 0 || player.Fighter.Stats.Healing <= 0 {
		return 0
	}
	amount := int(math.Round(player.MaxHP * float64(player.Fighter.Stats.Healing) / 100))
	before := player.HP
	player.HP = math.Min(player.MaxHP, player.HP+float64(amount))
	return int(math.Round(player.HP - before))
}

func (combat *Combat) partyAlive() bool {
	for _, player := range combat.Players {
		if player.Alive && player.HP > 0 {
			return true
		}
	}
	return false
}

func (combat *Combat) attackPlayer(index, turn int) []string {
	return combat.attackPlayerHits(index, turn)
}

func (combat *Combat) attackPlayerHits(index, turn int) []string {
	player := &combat.Players[index]
	attacker := player.Fighter.Stats
	attacker.Dmg += combat.activeFrenzy(index, turn)
	hits := attackHits(attacker, combat.Monster.Stats, combat.RNG, &combat.criticalIndex, combat.CriticalSequence, combat.nextCritical(index))
	logs := []string{}
	totalDamage := 0
	shieldUsed := false
	for _, hit := range hits {
		if !shieldUsed && combat.Monster.Stats.SunShield > 0 {
			hit.Damage = int(math.Round(float64(hit.Damage) * (1 - float64(combat.Monster.Stats.SunShield)/100)))
			shieldUsed = true
		}
		absorbed := min(hit.Damage, int(combat.Monster.Barrier))
		combat.Monster.Barrier -= float64(absorbed)
		damage := hit.Damage - absorbed
		totalDamage += damage
		combat.Monster.HP -= float64(damage)
		if combat.Monster.HP < 0 {
			combat.Monster.HP = 0
		}
		logs = append(logs, combat.attackLog(index, turn, hit, absorbed))
		if hit.Critical && player.Fighter.Stats.Lifesteal > 0 {
			heal := int(math.Round(float64(hit.Damage) * float64(player.Fighter.Stats.Lifesteal) / 100))
			before := player.HP
			player.HP = math.Min(player.MaxHP, player.HP+float64(heal))
			healed := int(math.Round(player.HP - before))
			if healed > 0 {
				logs = append(logs, fmt.Sprintf("Turn %d: Character_%d heals %d HP from lifesteal. HP: %d/%d", turn, index+1, healed, int(player.HP), int(player.MaxHP)))
			}
		}
		if absorbed > 0 {
			logs = append(logs, combat.barrierLog(turn, absorbed))
			if combat.Monster.Barrier <= 0 {
				logs = append(logs, fmt.Sprintf("Turn %d: The monster's barrier is destroyed!", turn))
			}
		}
		if combat.Monster.Stats.Corrupted > 0 {
			resistance := combat.reduceMonsterResistance(hit.Element, combat.Monster.Stats.Corrupted)
			logs = append(logs, fmt.Sprintf("Turn %d: The monster's %s resistance is corrupted and decreases by %d%%. New resistance: %d%%", turn, hit.Element, combat.Monster.Stats.Corrupted, resistance))
		}
	}
	for _, hit := range hits {
		if hit.Critical && player.Fighter.Stats.VampiricStrike > 0 && player.PlayedTurns-player.VampiricLast >= 3 {
			ally := combat.lowestWoundedAlly(index)
			if ally >= 0 {
				heal := int(math.Round(float64(totalDamage) * float64(player.Fighter.Stats.VampiricStrike) / 100))
				before := combat.Players[ally].HP
				combat.Players[ally].HP = math.Min(combat.Players[ally].MaxHP, combat.Players[ally].HP+float64(heal))
				healed := int(math.Round(combat.Players[ally].HP - before))
				player.VampiricLast = player.PlayedTurns
				if healed > 0 {
					logs = append(logs, fmt.Sprintf("Turn %d: Character_%d's Vampiric Strike heals Character_%d for %d HP. Character_%d HP: %d/%d", turn, index+1, ally+1, healed, ally+1, int(combat.Players[ally].HP), int(combat.Players[ally].MaxHP)))
				}
			}
			break
		}
	}
	for _, hit := range hits {
		if hit.Critical && player.Fighter.Stats.Frenzy > 0 {
			for ally := range combat.Players {
				combat.Players[ally].FrenzyBonus += player.Fighter.Stats.Frenzy
				combat.Players[ally].FrenzyUntil = turn + len(combat.Players) + 1
			}
			logs = append(logs, fmt.Sprintf("Turn %d: Character_%d's Frenzy triggers on critical. +%d%% damage will apply on each ally's next turn, active until the end of Character_%d's next turn.", turn, index+1, player.Fighter.Stats.Frenzy, index+1))
			break
		}
	}
	if player.Fighter.Stats.Burn > 0 && combat.Monster.BurnDamage == 0 {
		combat.Monster.BurnDamage = burnAmount(player.Fighter.Stats, player.Fighter.Stats.Burn)
		logs = append(logs, fmt.Sprintf("Turn %d: Character_%d applies a burn of %d on the monster. (Monster burn: %d)", turn, index+1, int(math.Floor(combat.Monster.BurnDamage)), int(math.Floor(combat.Monster.BurnDamage))))
	}
	return logs
}

func (combat *Combat) reduceMonsterResistance(element string, amount int) int {
	switch element {
	case "fire":
		combat.Monster.Stats.ResFire -= amount
		return combat.Monster.Stats.ResFire
	case "earth":
		combat.Monster.Stats.ResEarth -= amount
		return combat.Monster.Stats.ResEarth
	case "water":
		combat.Monster.Stats.ResWater -= amount
		return combat.Monster.Stats.ResWater
	case "air":
		combat.Monster.Stats.ResAir -= amount
		return combat.Monster.Stats.ResAir
	}
	return 0
}

func (combat *Combat) refreshBubble(turn int) string {
	if combat.Monster.Stats.ProtectiveBubble <= 0 {
		return ""
	}
	previous := combat.Monster.BubbleElement
	choice := int(combat.RNG() * 4)
	if combat.bubbleIndex < len(combat.BubbleDecisions) {
		choice = combat.BubbleDecisions[combat.bubbleIndex]
		combat.bubbleIndex++
	}
	if previous >= 0 && choice == previous {
		choice = (choice + 1) % 4
	}
	choice %= 4
	combat.Monster.Stats.ResFire = combat.Monster.BaseRes[0]
	combat.Monster.Stats.ResEarth = combat.Monster.BaseRes[1]
	combat.Monster.Stats.ResWater = combat.Monster.BaseRes[2]
	combat.Monster.Stats.ResAir = combat.Monster.BaseRes[3]
	resistances := []*int{
		&combat.Monster.Stats.ResFire,
		&combat.Monster.Stats.ResEarth,
		&combat.Monster.Stats.ResWater,
		&combat.Monster.Stats.ResAir,
	}
	*resistances[choice] += combat.Monster.Stats.ProtectiveBubble
	combat.Monster.BubbleElement = choice
	elements := []string{"fire", "earth", "water", "air"}
	if previous >= 0 {
		return fmt.Sprintf("Turn %d: %s's Protective Bubble shifts from %s to %s, gaining %d%% %s resistance!", turn, combat.Monster.Name, elements[previous], elements[choice], combat.Monster.Stats.ProtectiveBubble, elements[choice])
	}
	return fmt.Sprintf("Turn %d: The monster's Protective Bubble activates, gaining %d%% %s resistance!", turn, combat.Monster.Stats.ProtectiveBubble, elements[choice])
}

func (combat *Combat) lowestWoundedAlly(caster int) int {
	lowest := -1
	ratio := 1.0
	for index := range combat.Players {
		player := &combat.Players[index]
		if index != caster && player.Alive && player.HP < player.MaxHP && player.HP/player.MaxHP < ratio {
			lowest = index
			ratio = player.HP / player.MaxHP
		}
	}
	return lowest
}

func (combat *Combat) activeFrenzy(index, turn int) int {
	player := &combat.Players[index]
	if turn > player.FrenzyUntil {
		player.FrenzyBonus = 0
	}
	return player.FrenzyBonus + player.Fighter.Stats.Greed*player.GreedStacks
}

func (combat *Combat) attackMonster(index, turn int) []string {
	combat.Monster.Turns++
	logs := []string{}
	if combat.Monster.Stats.SunShield > 0 {
		logs = append(logs, fmt.Sprintf("Turn %d: Sun Shield activates. The first hit each turn is reduced by %d%%.", turn, combat.Monster.Stats.SunShield))
	}
	if combat.Monster.Reconstitution > 0 && combat.Monster.Turns%20 == 0 {
		combat.Monster.HP = combat.Monster.MaxHP
		logs = append(logs, fmt.Sprintf("Turn %d: The monster uses Reconstitution and restores all HP.", turn))
	}
	if combat.Monster.VoidDrain > 0 && combat.Monster.Turns%4 == 0 && combat.Players[index].HP > 0 {
		drain := math.Min(combat.Players[index].HP, math.Round(combat.Players[index].MaxHP*float64(combat.Monster.VoidDrain)/100))
		combat.Players[index].HP -= drain
		combat.Monster.HP = math.Min(combat.Monster.MaxHP, combat.Monster.HP+drain)
		logs = append(logs, fmt.Sprintf("Turn %d: The monster uses Void Drain and drains %d HP from Character_%d.", turn, int(drain), index+1))
	}
	if combat.Monster.Stats.Healing > 0 && combat.Monster.Turns%3 == 0 {
		heal := int(math.Round(combat.Monster.MaxHP * float64(combat.Monster.Stats.Healing) / 100))
		before := combat.Monster.HP
		combat.Monster.HP = math.Min(combat.Monster.MaxHP, combat.Monster.HP+float64(heal))
		healed := int(math.Round(combat.Monster.HP - before))
		if healed > 0 {
			logs = append(logs, fmt.Sprintf("Turn %d: The monster heals %d HP from Healing effect.(%s HP: %d/%d)", turn, healed, combat.Monster.Name, int(combat.Monster.HP), int(combat.Monster.MaxHP)))
		}
	}
	if combat.Monster.Stats.Burn > 0 && !combat.Players[index].BurnApplied {
		combat.Players[index].BurnDamage = burnAmount(combat.Monster.Stats, combat.Monster.Stats.Burn)
		combat.Players[index].BurnApplied = true
		logs = append(logs, fmt.Sprintf("Turn %d: The monster applies a burn of %d on Character_%d. (Character burn: %d)", turn, int(math.Floor(combat.Players[index].BurnDamage)), index+1, int(math.Floor(combat.Players[index].BurnDamage))))
	}
	if combat.Players[index].Fighter.Stats.Burn > 0 && combat.Monster.BurnDamage > 0 {
		logs = append(logs, combat.applyMonsterBurn(turn)...)
	}
	if combat.Monster.BarrierAmount > 0 && combat.Monster.Turns%5 == 0 {
		combat.Monster.Barrier = combat.Monster.BarrierAmount
		logs = append(logs, fmt.Sprintf("Turn %d: The monster refreshes its barrier to %d HP.", turn, int(combat.Monster.BarrierAmount)))
	}
	player := &combat.Players[index]
	defender := player.Fighter.Stats
	if player.ShellTurns > 0 {
		defender.ResFire += player.Fighter.Stats.Shell
		defender.ResEarth += player.Fighter.Stats.Shell
		defender.ResWater += player.Fighter.Stats.Shell
		defender.ResAir += player.Fighter.Stats.Shell
	}
	wasAlive := player.Alive
	totalDamage := 0
	guarder := combat.guarderFor(index)
	attacker := combat.Monster.Stats
	if combat.Monster.FrenzyPending {
		attacker.Dmg += combat.Monster.Stats.Frenzy
		combat.Monster.FrenzyPending = false
		logs = append(logs, fmt.Sprintf("Turn %d: The monster's Frenzy activates. +%d%% damage applies on its next attack.", turn, combat.Monster.Stats.Frenzy))
	}
	if combat.Monster.Berserked {
		attacker.Dmg += combat.Monster.BerserkerRage
	}
	attacker.Dmg += combat.Monster.Stats.Greed * combat.Monster.GreedStacks
	hits := attackHits(attacker, defender, combat.RNG, &combat.criticalIndex, combat.CriticalSequence, combat.nextCritical(-1))
	for _, hit := range hits {
		damage := hit.Damage
		player.HP -= float64(damage)
		totalDamage += hit.Damage
		if player.HP <= 0 {
			player.HP = 0
		}
		logs = append(logs, combat.monsterLog(index, turn, hit))
		if hit.Critical && combat.Monster.Stats.Lifesteal > 0 {
			heal := int(math.Floor(float64(hit.Damage) * float64(combat.Monster.Stats.Lifesteal) / 100))
			before := combat.Monster.HP
			combat.Monster.HP = math.Min(combat.Monster.MaxHP, combat.Monster.HP+float64(heal))
			healed := int(math.Round(combat.Monster.HP - before))
			if healed > 0 {
				logs = append(logs, fmt.Sprintf("Turn %d: Monster heals %d HP from lifesteal. %s HP: %d/%d", turn, healed, combat.Monster.Name, int(combat.Monster.HP), int(combat.Monster.MaxHP)))
			}
		}
	}
	if guarder >= 0 {
		redirected := int(math.Round(float64(totalDamage) * float64(combat.Players[guarder].Fighter.Stats.Guard) / 100))
		combat.Players[guarder].HP = math.Max(0, combat.Players[guarder].HP-float64(redirected))
		if combat.Players[guarder].HP <= 0 {
			combat.Players[guarder].HP, combat.Players[guarder].Alive = 0, false
		}
		if redirected > 0 {
			logs = append(logs, fmt.Sprintf("Turn %d: Character_%d's Guard redirects %d damage from Character_%d. Character_%d HP: %d/%d", turn, guarder+1, redirected, index+1, guarder+1, int(combat.Players[guarder].HP), int(combat.Players[guarder].MaxHP)))
		}
	}
	if player.HP <= 0 {
		player.Alive = false
	}
	if combat.Monster.Stats.Greed > 0 {
		thresholds := greedThresholds(combat.Monster.MaxHP, combat.Monster.HP)
		for combat.Monster.GreedStacks < thresholds {
			combat.Monster.GreedStacks++
			total := combat.Monster.Stats.Greed * combat.Monster.GreedStacks
			logs = append(logs, fmt.Sprintf("Turn %d: Greed empowers the monster (+%d%% damage, total +%d%%).", turn, combat.Monster.Stats.Greed, total))
		}
	}
	for _, hit := range hits {
		if hit.Critical && combat.Monster.Stats.Frenzy > 0 {
			combat.Monster.FrenzyPending = true
			break
		}
	}
	if combat.Monster.BerserkerRage > 0 && !combat.Monster.Berserked && combat.Monster.HP < combat.Monster.MaxHP*0.25 {
		combat.Monster.Berserked = true
		logs = append(logs, fmt.Sprintf("Turn %d: The monster's Berserker Rage activates, gaining %d%% permanent damage!", turn, combat.Monster.BerserkerRage))
	}
	if player.Fighter.Stats.Greed > 0 {
		if player.GreedAwakened {
			thresholds := greedThresholds(player.MaxHP, player.HP)
			for player.GreedStacks < thresholds {
				player.GreedStacks++
				total := player.Fighter.Stats.Greed * player.GreedStacks
				logs = append(logs, fmt.Sprintf("Turn %d: Greed empowers Character_%d (+%d%% damage, total +%d%%).", turn, index+1, player.Fighter.Stats.Greed, total))
			}
		} else {
			player.GreedAwakened = true
		}
	}
	if player.Fighter.Stats.EnchantedMirror > 0 && combat.Monster.Turns%3 == 1 && totalDamage > 0 {
		reflected := int(math.Round(float64(totalDamage) * float64(player.Fighter.Stats.EnchantedMirror) / 100))
		combat.Monster.HP -= float64(reflected)
		if combat.Monster.HP < 0 {
			combat.Monster.HP = 0
		}
		logs = append(logs, fmt.Sprintf("Turn %d: Character_%d's Enchanted Mirror activates, dealing back %d damage to %s.", turn, index+1, reflected, combat.Monster.Name))
		logs = append(logs, fmt.Sprintf("Turn %d: %s takes %d reflected damage. %s HP: %d/%d", turn, combat.Monster.Name, reflected, combat.Monster.Name, int(combat.Monster.HP), int(combat.Monster.MaxHP)))
	}
	if wasAlive && !player.Alive {
		logs = append(logs, fmt.Sprintf("Turn %d: Character_%d has been defeated!", turn, index+1))
	}
	if player.ShellTurns > 0 {
		player.ShellTurns--
		if player.ShellTurns == 0 {
			logs = append(logs, fmt.Sprintf("Turn %d: Character_%d's shell effect has worn off.", turn, index+1))
		}
	}
	return logs
}

func (combat *Combat) guarderFor(target int) int {
	for index := range combat.Players {
		player := &combat.Players[index]
		if player.Alive && player.GuardTarget == target && player.GuardTurns > 0 {
			return index
		}
	}
	return -1
}

func (combat *Combat) expireGuard(index int) {
	player := &combat.Players[index]
	if player.GuardTarget >= 0 {
		player.GuardTarget = -1
		player.GuardTurns = 0
	}
}

func (combat *Combat) activateGuard(index, turn int) string {
	player := &combat.Players[index]
	if player.Fighter.Stats.Guard <= 0 || player.GuardUses >= 3 {
		return ""
	}
	target := -1
	if combat.guardIndex < len(combat.GuardDecisions) && combat.GuardDecisions[combat.guardIndex].Turn == turn && combat.GuardDecisions[combat.guardIndex].Protector == index {
		target = combat.GuardDecisions[combat.guardIndex].Target
		combat.guardIndex++
	}
	if target >= 0 {
		player.GuardTarget = target
		player.GuardTurns = 1
		player.GuardUses++
		return fmt.Sprintf("Turn %d: Character_%d activates Guard to protect Character_%d, redirecting %d%% of damage until next turn. (%d/3 activations used)", turn, index+1, target+1, player.Fighter.Stats.Guard, player.GuardUses)
	}
	if len(combat.GuardDecisions) > 0 {
		return ""
	}
	ratio := 1.0
	for candidate := range combat.Players {
		ally := &combat.Players[candidate]
		if candidate != index && ally.Alive && ally.HP < ally.MaxHP*0.5 && ally.HP/ally.MaxHP < ratio {
			target = candidate
			ratio = ally.HP / ally.MaxHP
		}
	}
	if target < 0 {
		return ""
	}
	player.GuardTarget = target
	player.GuardTurns = 1
	player.GuardUses++
	return fmt.Sprintf("Turn %d: Character_%d activates Guard to protect Character_%d, redirecting %d%% of damage until next turn. (%d/3 activations used)", turn, index+1, target+1, player.Fighter.Stats.Guard, player.GuardUses)
}

func (combat *Combat) applyPlayerBurn(index, turn int) []string {
	player := &combat.Players[index]
	if player.BurnDamage == 0 {
		return nil
	}
	damage := int(math.Floor(player.BurnDamage))
	if damage <= 0 {
		return nil
	}
	player.HP = math.Max(0, player.HP-float64(damage))
	if player.HP == 0 {
		player.Alive = false
	}
	player.BurnDamage = decayBurnDamage(damage)
	return []string{fmt.Sprintf("Turn %d: Character_%d suffers from burn and loses %d HP. Character_%d HP: %d/%d", turn, index+1, damage, index+1, int(player.HP), int(player.MaxHP))}
}

func (combat *Combat) applyMonsterBurn(turn int) []string {
	damage := int(math.Floor(combat.Monster.BurnDamage))
	if damage <= 0 {
		return nil
	}
	combat.Monster.HP = math.Max(0, combat.Monster.HP-float64(damage))
	combat.Monster.BurnDamage = decayBurnDamage(damage)
	return []string{fmt.Sprintf("Turn %d: Monster suffers %d burn damage. %s HP: %d/%d", turn, damage, combat.Monster.Name, int(combat.Monster.HP), int(combat.Monster.MaxHP))}
}

func burnAmount(stats Stats, percent int) float64 {
	total := stats.AttackFire + stats.AttackEarth + stats.AttackWater + stats.AttackAir
	return float64(total) * float64(percent) / 100
}

func greedThresholds(maxHP, hp float64) int {
	block := math.Floor(maxHP * 0.10)
	if block < 1 {
		block = 1
	}
	lost := math.Max(0, maxHP-math.Max(0, hp))
	return int(math.Floor(lost / block))
}

func decayBurnDamage(damage int) float64 {
	if damage <= 5 {
		return math.Floor(float64(damage) * 0.9)
	}
	return math.Round(float64(damage) * 0.9)
}

func (combat *Combat) activateShell(index int) bool {
	player := &combat.Players[index]
	if player.ShellUsed || player.Fighter.Stats.Shell <= 0 {
		return false
	}
	if (combat.Monster.Type != "boss" && combat.Monster.Type != "raid_boss") || player.HP >= player.MaxHP*0.40 {
		return false
	}
	player.ShellUsed = true
	player.ShellTurns = 3
	return true
}

type hit struct {
	Element  string
	Damage   int
	Critical bool
}

func (combat *Combat) nextCritical(actor int) *bool {
	if len(combat.CriticalDecisions) == 0 {
		return nil
	}
	if len(combat.decisionUsed) != len(combat.CriticalDecisions) {
		combat.decisionUsed = make([]bool, len(combat.CriticalDecisions))
	}
	for index, decision := range combat.CriticalDecisions {
		if !combat.decisionUsed[index] && decision.Actor == actor {
			combat.decisionUsed[index] = true
			value := decision.Critical
			return &value
		}
	}
	return nil
}

func attackHits(attacker, defender Stats, rng RNG, index *int, sequence []bool, forced ...*bool) []hit {
	attacks := []struct {
		attack, bonus, resistance int
		element                   string
	}{
		{
			attack:     attacker.AttackFire,
			bonus:      attacker.Dmg + attacker.DmgFire,
			resistance: defender.ResFire,
			element:    "fire",
		},
		{
			attack:     attacker.AttackEarth,
			bonus:      attacker.Dmg + attacker.DmgEarth,
			resistance: defender.ResEarth,
			element:    "earth",
		},
		{
			attack:     attacker.AttackWater,
			bonus:      attacker.Dmg + attacker.DmgWater,
			resistance: defender.ResWater,
			element:    "water",
		},
		{
			attack:     attacker.AttackAir,
			bonus:      attacker.Dmg + attacker.DmgAir,
			resistance: defender.ResAir,
			element:    "air",
		},
	}
	critical := false
	if len(forced) > 0 && forced[0] != nil {
		critical = *forced[0]
	} else if *index < len(sequence) {
		critical = sequence[*index]
		*index++
	} else {
		critical = rng() < float64(attacker.CriticalStrike)/100
	}
	hits := []hit{}
	for _, attack := range attacks {
		if attack.attack == 0 {
			continue
		}
		damage := elemental(attack.attack, attack.bonus, attack.resistance)
		if critical {
			damage = int(math.Round(float64(damage) * 1.5))
		}
		hits = append(hits, hit{attack.element, damage, critical})
	}
	return hits
}

func elemental(attack, bonus, resistance int) int {
	if attack <= 0 {
		return 0
	}
	raw := math.Round(float64(attack) * (1 + float64(bonus)/100))
	value := int(math.Round(raw * math.Max(0, 1-float64(resistance)/100)))
	if value < 0 {
		return 0
	}
	return value
}
func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}
func (combat *Combat) attackLog(index, turn int, hit hit, absorbed int) string {
	suffix := ""
	if hit.Critical {
		suffix = " (Critical strike)"
	}
	detail := ""
	if absorbed == hit.Damage {
		detail = ", fully absorbed by the barrier"
	} else if absorbed > 0 {
		detail = fmt.Sprintf(" (%d absorbed, %d through)", absorbed, hit.Damage-absorbed)
	}
	return fmt.Sprintf("Turn %d: Character_%d used %s attack and dealt %d damage%s%s. %s HP: %d/%d", turn, index+1, hit.Element, hit.Damage, suffix, detail, combat.Monster.Name, max(0, int(combat.Monster.HP)), int(combat.Monster.MaxHP))
}

func (combat *Combat) barrierLog(turn, absorbed int) string {
	return fmt.Sprintf("Turn %d: The monster's barrier absorbs %d damage. Barrier HP: %d", turn, absorbed, max(0, int(combat.Monster.Barrier)))
}
func (combat *Combat) monsterLog(index, turn int, hit hit) string {
	suffix := ""
	if hit.Critical {
		suffix = " (Critical strike)"
	}
	player := combat.Players[index]
	return fmt.Sprintf("Turn %d: %s used %s attack against Character_%d and dealt %d damage%s. Character_%d HP: %d/%d", turn, combat.Monster.Name, hit.Element, index+1, hit.Damage, suffix, index+1, max(0, int(player.HP)), int(player.MaxHP))
}
