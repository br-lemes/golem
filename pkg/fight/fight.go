// Package fight contains the new, standalone combat simulator.
package fight

import (
	"fmt"
	"math"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

const MaxTurns = 100

type Stats struct {
	HP                                              int
	AttackFire, AttackEarth, AttackWater, AttackAir int
	Dmg, DmgFire, DmgEarth, DmgWater, DmgAir        int
	ResFire, ResEarth, ResWater, ResAir             int
	CriticalStrike, Initiative, Haste, Wisdom       int
	Burn                                            int
	Lifesteal                                       int
	Healing                                         int
	Shell                                           int
	EnchantedMirror                                 int
	Frenzy                                          int
	Greed                                           int
}

type Fighter struct {
	Stats     Stats
	Utilities []Utility
	NaturalHP int
}

type Utility struct {
	Code       string
	Restore    int
	Antipoison int
	Quantity   int
}

// FromLoadout creates a fighter without any character/API dependency. The
// The official simulator currently reports 120 HP and 100 initiative for a
// level-one fake character. HP gains 5 points per additional level.
func FromLoadout(level int, slots map[string]string, utilities map[string]int) Fighter {
	if level < 1 {
		level = 1
	}
	stats := Stats{HP: 120 + (level-1)*5, Initiative: 100}
	for _, code := range slots {
		applyItem(&stats, code, 1)
	}
	naturalHP := stats.HP
	result := Fighter{Stats: stats, NaturalHP: naturalHP}
	return applyUtilities(result, utilities)
}

func applyUtilities(result Fighter, utilities map[string]int) Fighter {
	stats := result.Stats
	for code, quantity := range utilities {
		if quantity <= 0 {
			continue
		}
		item, ok := database.Items().Get(code)
		if !ok || item.Effects == nil {
			continue
		}
		utility := Utility{Code: code, Quantity: quantity}
		for _, effect := range *item.Effects {
			switch effect.Code {
			case "boost_hp", "boost_dmg_fire", "boost_dmg_earth", "boost_dmg_water", "boost_dmg_air", "boost_res_fire", "boost_res_earth", "boost_res_water", "boost_res_air":
				// Boost potions apply once when combat starts. Quantity is only
				// relevant to consumables used during combat.
				applyCombatBoost(&stats, effect.Code, effect.Value)
			case "restore":
				utility.Restore += effect.Value
			case "antipoison":
				utility.Antipoison += effect.Value
			}
		}
		if utility.Restore > 0 || utility.Antipoison > 0 {
			result.Utilities = append(result.Utilities, utility)
		}
	}
	result.Stats = stats
	return result
}

func applyCombatBoost(stats *Stats, code string, value int) {
	switch code {
	case "boost_hp":
		stats.HP += value
	case "boost_dmg_fire":
		stats.DmgFire += value
	case "boost_dmg_earth":
		stats.DmgEarth += value
	case "boost_dmg_water":
		stats.DmgWater += value
	case "boost_dmg_air":
		stats.DmgAir += value
	case "boost_res_fire":
		stats.ResFire += value
	case "boost_res_earth":
		stats.ResEarth += value
	case "boost_res_water":
		stats.ResWater += value
	case "boost_res_air":
		stats.ResAir += value
	}
}

func applyItem(stats *Stats, code string, sign int) {
	if code == "" {
		return
	}
	item, ok := database.Items().Get(code)
	if !ok || item.Effects == nil {
		return
	}
	for _, effect := range *item.Effects {
		value := sign * effect.Value
		switch effect.Code {
		case "hp":
			stats.HP += value
		case "attack_fire":
			stats.AttackFire += value
		case "attack_earth":
			stats.AttackEarth += value
		case "attack_water":
			stats.AttackWater += value
		case "attack_air":
			stats.AttackAir += value
		case "dmg":
			stats.Dmg += value
		case "dmg_fire":
			stats.DmgFire += value
		case "dmg_earth":
			stats.DmgEarth += value
		case "dmg_water":
			stats.DmgWater += value
		case "dmg_air":
			stats.DmgAir += value
		case "res_fire":
			stats.ResFire += value
		case "res_earth":
			stats.ResEarth += value
		case "res_water":
			stats.ResWater += value
		case "res_air":
			stats.ResAir += value
		case "critical_strike":
			stats.CriticalStrike += value
		case "initiative":
			stats.Initiative += value
		case "haste":
			stats.Haste += value
		case "wisdom":
			stats.Wisdom += value
		case "burn":
			stats.Burn += value
		case "lifesteal":
			stats.Lifesteal += value
		case "healing":
			stats.Healing += value
		case "shell":
			stats.Shell += value
		case "enchanted_mirror":
			stats.EnchantedMirror += value
		case "frenzy":
			stats.Frenzy += value
		case "greed":
			stats.Greed += value
		}
	}
}

type Result struct {
	Win         bool     `json:"win"`
	Turns       int      `json:"turns"`
	HPRemaining int      `json:"hp_remaining"`
	MaxHP       int      `json:"max_hp"`
	PlayerFirst bool     `json:"player_first"`
	TimedOut    bool     `json:"timed_out"`
	Logs        []string `json:"-"`
}

type Summary struct {
	Results []schemas.CombatResultSchema `json:"results"`
	Wins    int                          `json:"wins"`
	Losses  int                          `json:"losses"`
	Winrate float32                      `json:"winrate"`
}

type RNG func() float64

type CriticalOptions struct {
	PlayerChance  *int
	MonsterChance *int
}

type SimulationOptions struct {
	Iterations int
	RNG        RNG
	Critical   CriticalOptions
	// CriticalSequence, when provided, fixes the result of each complete
	// attack's critical roll. Elements within one attack share that result.
	CriticalSequence []bool
}

func SimulateMany(player Fighter, monster schemas.MonsterSchema, options SimulationOptions) Summary {
	if options.Iterations < 1 {
		options.Iterations = 1
	}
	summary := Summary{
		Results: make([]schemas.CombatResultSchema, 0, options.Iterations),
	}
	for i := 0; i < options.Iterations; i++ {
		r := simulate(player, monster, options)
		result := "loss"
		if r.Win {
			result = "win"
			summary.Wins++
		} else {
			summary.Losses++
		}
		entry := schemas.CombatResultSchema{
			Result: result,
			Turns:  r.Turns,
			Logs:   r.Logs,
			CharacterResults: []map[string]interface{}{{
				"final_hp":               r.HPRemaining,
				"utility1_slot_quantity": 0,
				"utility2_slot_quantity": 0,
			}},
		}
		summary.Results = append(summary.Results, entry)
	}
	summary.Winrate = float32(summary.Wins) * 100 / float32(options.Iterations)
	return summary
}

func Simulate(player Fighter, monster schemas.MonsterSchema, options SimulationOptions) Result {
	return simulate(player, monster, options)
}

func simulate(player Fighter, monster schemas.MonsterSchema, options SimulationOptions) Result {
	criticalIndex := 0
	if options.Critical.PlayerChance != nil {
		player.Stats.CriticalStrike = *options.Critical.PlayerChance
	}
	if options.Critical.MonsterChance != nil {
		monster.CriticalStrike = *options.Critical.MonsterChance
	}
	p := combatant{
		hp:           float64(player.Stats.HP),
		maxHP:        float64(player.Stats.HP),
		stats:        player.Stats,
		damageFactor: 1,
	}
	m := combatant{
		hp:    float64(monster.Hp),
		maxHP: float64(monster.Hp),
		stats: Stats{
			HP:             monster.Hp,
			AttackFire:     monster.AttackFire,
			AttackEarth:    monster.AttackEarth,
			AttackWater:    monster.AttackWater,
			AttackAir:      monster.AttackAir,
			CriticalStrike: monster.CriticalStrike,
			Initiative:     monster.Initiative,
			ResFire:        monster.ResFire,
			ResEarth:       monster.ResEarth,
			ResWater:       monster.ResWater,
			ResAir:         monster.ResAir,
			Lifesteal:      monsterEffect(monster, "lifesteal"),
			Healing:        monsterEffect(monster, "healing"),
		},
		damageFactor: 1,
		barrier:      float64(monsterEffect(monster, "barrier")),
	}
	first := p.stats.Initiative >= m.stats.Initiative
	naturalHP := player.NaturalHP
	if naturalHP <= 0 {
		naturalHP = player.Stats.HP
	}
	result := Result{MaxHP: naturalHP, PlayerFirst: first}
	logs := []string{fmt.Sprintf("Fight start: Character_1 HP: %d/%d vs %s HP: %d/%d", player.Stats.HP, player.Stats.HP, monster.Name, monster.Hp, monster.Hp)}
	if m.barrier > 0 {
		logs = append([]string{fmt.Sprintf("Fight start: The monster raises a barrier of %d HP.", int(m.barrier))}, logs...)
	}
	poison := monsterEffect(monster, "poison")
	poisonApplied := false
	poisonStarted := false
	burn := monsterEffect(monster, "burn")
	burnStarted := false
	berserkerRage := monsterEffect(monster, "berserker_rage")
	berserkerActivated := false
	corrupted := monsterEffect(monster, "corrupted")
	monsterFrenzy := monsterEffect(monster, "frenzy")
	greed := monsterEffect(monster, "greed")
	protectiveBubble := monsterEffect(monster, "protective_bubble")
	reconstitution := monsterEffect(monster, "reconstitution")
	voidDrain := monsterEffect(monster, "void_drain")
	baseMonsterRes := [4]int{
		m.stats.ResFire,
		m.stats.ResEarth,
		m.stats.ResWater,
		m.stats.ResAir,
	}
	bubbleElement := -1
	greedStacks := 0
	if greed > 0 {
		logs = append(logs, fmt.Sprintf("Turn 1: Greed awakens. The monster gains +%d%% damage each time it loses 10%% max HP.", greed))
	}
	// The API applies burn as integer damage, reducing the previous tick by
	// 10% and rounding after each tick.
	burnDamage := float64(monster.AttackFire+monster.AttackEarth+monster.AttackWater+monster.AttackAir) * float64(burn) / 100
	playerAttackTotal := elemental(player.Stats.AttackFire, player.Stats.Dmg+player.Stats.DmgFire, 0)
	playerAttackTotal += elemental(player.Stats.AttackEarth, player.Stats.Dmg+player.Stats.DmgEarth, 0)
	playerAttackTotal += elemental(player.Stats.AttackWater, player.Stats.Dmg+player.Stats.DmgWater, 0)
	playerAttackTotal += elemental(player.Stats.AttackAir, player.Stats.Dmg+player.Stats.DmgAir, 0)
	playerBurnDamage := float64(playerAttackTotal) * float64(player.Stats.Burn) / 100
	playerBurnStarted := false
	playerTurns := 0
	monsterTurns := 0
	shellTurns := 0
	shellActivated := false
	playerFrenzyNext := false
	playerFrenzyActive := false
	monsterFrenzyNext := false
	monsterFrenzyActive := false
	playerGreedStacks := 0
	if player.Stats.Greed > 0 {
		logs = append(logs, fmt.Sprintf("Turn 1: Greed awakens. Character_1 gains +%d%% damage each time they lose 10%% max HP.", player.Stats.Greed))
	}
	utilities := append([]Utility(nil), player.Utilities...)
	for turn := 1; turn <= MaxTurns; turn++ {
		if protectiveBubble > 0 {
			m.stats.ResFire, m.stats.ResEarth, m.stats.ResWater, m.stats.ResAir = baseMonsterRes[0], baseMonsterRes[1], baseMonsterRes[2], baseMonsterRes[3]
			bubbleElement = randomBubbleElement(options.RNG, bubbleElement)
			res := []*int{
				&m.stats.ResFire,
				&m.stats.ResEarth,
				&m.stats.ResWater,
				&m.stats.ResAir,
			}
			*res[bubbleElement] += protectiveBubble
			bubbleElements := []string{"fire", "earth", "water", "air"}
			bubbleMessage := fmt.Sprintf("Turn %d: The monster's protective bubble grants %d%% %s resistance.", turn, protectiveBubble, bubbleElements[bubbleElement])
			logs = append(logs, bubbleMessage)
		}
		attacker, defender := &m, &p
		playerTurn := (turn%2 == 1) == first
		if playerTurn {
			playerTurns++
			if playerFrenzyNext {
				playerFrenzyNext = false
				playerFrenzyActive = true
				greedBonus := playerGreedStacks - 1
				if greedBonus < 0 {
					greedBonus = 0
				}
				p.damageFactor = 1 + float64(player.Stats.Greed*greedBonus+player.Stats.Frenzy)/100
			} else {
				playerFrenzyActive = false
				greedBonus := playerGreedStacks - 1
				if greedBonus < 0 {
					greedBonus = 0
				}
				p.damageFactor = 1 + float64(player.Stats.Greed*greedBonus)/100
			}
			if player.Stats.Shell > 0 && !shellActivated && isBoss(monster) && p.hp < p.maxHP*0.40 {
				shellActivated = true
				shellTurns = 3
				p.stats.ResFire += player.Stats.Shell
				p.stats.ResEarth += player.Stats.Shell
				p.stats.ResWater += player.Stats.Shell
				p.stats.ResAir += player.Stats.Shell
				logs = append(logs, fmt.Sprintf("Turn %d: Character_1's shell activates, gaining %d%% resistance to all elements for 3 turns!", turn, player.Stats.Shell))
			} else if shellActivated && shellTurns > 0 {
				shellTurns--
				if shellTurns == 0 {
					p.stats.ResFire -= player.Stats.Shell
					p.stats.ResEarth -= player.Stats.Shell
					p.stats.ResWater -= player.Stats.Shell
					p.stats.ResAir -= player.Stats.Shell
					logs = append(logs, fmt.Sprintf("Turn %d: Character_1's shell effect has worn off.", turn))
				}
			}
			if player.Stats.Healing > 0 && playerTurns%3 == 0 {
				heal := int(math.Round(float64(player.Stats.HP) * float64(player.Stats.Healing) / 100))
				before := p.hp
				p.hp = math.Min(p.maxHP, p.hp+float64(heal))
				healed := int(math.Round(p.hp - before))
				if healed > 0 {
					logs = append(logs, fmt.Sprintf("Turn %d: Character_1 heals %d HP from Healing effect. HP: %d/%d", turn, healed, int(math.Round(p.hp)), player.Stats.HP))
				}
			}
			if poisonApplied {
				poisonReduction := consumeAntidote(&utilities, &logs, turn)
				if poisonReduction > 0 {
					poison -= poisonReduction
					if poison <= 0 {
						poison = 0
						poisonApplied = false
					}
				}
			}
			consumeUtility(&utilities, false, &p, player.Stats.HP, &logs, turn)
		} else {
			monsterTurns++
			if reconstitution > 0 && monsterTurns%20 == 0 {
				m.hp = m.maxHP
				logs = append(logs, fmt.Sprintf("Turn %d: The monster uses Reconstitution and restores all HP.", turn))
			}
			if voidDrain > 0 && monsterTurns%4 == 0 && p.hp > 0 {
				drain := math.Min(p.hp, math.Round(p.maxHP*float64(voidDrain)/100))
				p.hp -= drain
				m.hp = math.Min(m.maxHP, m.hp+drain)
				logs = append(logs, fmt.Sprintf("Turn %d: The monster uses Void Drain and drains %d HP from Character_1.", turn, int(drain)))
			}
			if monsterFrenzyNext {
				monsterFrenzyNext = false
				monsterFrenzyActive = true
				m.damageFactor *= 1 + float64(monsterFrenzy)/100
				logs = append(logs, fmt.Sprintf("Turn %d: The monster's Frenzy activates. +%d%% damage applies on its next attack.", turn, monsterFrenzy))
			}
			if m.stats.Healing > 0 && monsterTurns%3 == 0 {
				heal := int(math.Round(float64(monster.Hp) * float64(m.stats.Healing) / 100))
				before := m.hp
				m.hp = math.Min(m.maxHP, m.hp+float64(heal))
				healed := int(math.Round(m.hp - before))
				if healed > 0 {
					logs = append(logs, fmt.Sprintf("Turn %d: The monster heals %d HP from Healing effect.(%s HP: %d/%d)", turn, healed, monster.Name, int(math.Round(m.hp)), monster.Hp))
				}
			}
			barrier := monsterEffect(monster, "barrier")
			if barrier > 0 && monsterTurns%5 == 0 {
				m.barrier = float64(barrier)
				logs = append(logs, fmt.Sprintf("Turn %d: The monster refreshes its barrier to %d HP.", turn, barrier))
			}
		}
		if playerTurn && poisonApplied {
			if poisonApplied {
				p.hp -= float64(poison)
			}
			if poisonApplied {
				logs = append(logs, fmt.Sprintf("Turn %d: Character_1 suffers from poison and loses %d HP. Character_1 HP: %d/%d", turn, poison, maxInt(0, int(math.Round(p.hp))), player.Stats.HP))
			}
		}
		if playerTurn && burnStarted && burnDamage > 0 {
			damage := int(math.Floor(burnDamage))
			if damage > 0 {
				p.hp -= float64(damage)
				logs = append(logs, fmt.Sprintf("Turn %d: Character_1 suffers from burn and loses %d HP. Character_1 HP: %d/%d", turn, damage, maxInt(0, int(math.Round(p.hp))), player.Stats.HP))
			}
			if damage <= 1 {
				burnDamage = 0
			} else {
				burnDamage = decayBurnDamage(damage)
			}
		}
		if playerTurn && !playerBurnStarted && player.Stats.Burn > 0 {
			playerBurnStarted = true
			logs = append(logs, fmt.Sprintf("Turn %d: Character_1 applies a burn of %d on the monster. (Monster burn: %d)", turn, int(math.Floor(playerBurnDamage)), int(math.Floor(playerBurnDamage))))
		}
		if !playerTurn && playerBurnStarted && playerBurnDamage > 0 {
			damage := int(math.Floor(playerBurnDamage))
			if damage > 0 {
				m.hp -= float64(damage)
				logs = append(logs, fmt.Sprintf("Turn %d: Monster suffers %d burn damage. %s HP: %d/%d", turn, damage, monster.Name, maxInt(0, int(math.Round(m.hp))), monster.Hp))
			}
			if damage <= 1 {
				playerBurnDamage = 0
			} else {
				playerBurnDamage = decayBurnDamage(damage)
			}
		}
		if !playerTurn && !poisonStarted && poison > 0 {
			poisonStarted = true
			poisonApplied = true
			logs = append(logs, fmt.Sprintf("Turn %d: The monster applies a poison of %d on Character_1. (Character poison: %d)", turn, poison, poison))
		}
		if !playerTurn && !burnStarted && burn > 0 {
			burnStarted = true
			logs = append(logs, fmt.Sprintf("Turn %d: The monster applies a burn of %d on Character_1.", turn, int(math.Floor(burnDamage))))
		}
		if !playerTurn && m.hp <= 0 {
			result.Win = true
			result.Turns = turn
			logs = append(logs, fmt.Sprintf("Turn %d: Monster has been defeated!", turn))
			break
		}
		if !playerTurn && berserkerRage > 0 && !berserkerActivated && m.hp < m.maxHP*0.25 {
			berserkerActivated = true
			m.damageFactor = 1 + float64(berserkerRage)/100
			logs = append(logs, fmt.Sprintf("Turn %d: The monster's Berserker Rage activates, gaining %d%% permanent damage!", turn, berserkerRage))
		}
		if playerTurn {
			attacker, defender = &p, &m
		}
		if defender.hp <= 0 {
			result.Win = playerTurn
			result.Turns = turn - 1
			break
		}
		hits := attackHits(attacker.stats, defender.stats, options.RNG, options.CriticalSequence, &criticalIndex)
		damage := 0
		monsterCriticalDamage := 0
		actor := monster.Name
		target := "Character_1"
		if playerTurn {
			actor, target = "Character_1", monster.Name
		}
		playerCritical := false
		for _, hit := range hits {
			if attacker.damageFactor != 1 {
				hit.damage = int(math.Round(float64(hit.damage) * attacker.damageFactor))
			}
			damage += hit.damage
			absorbed := 0
			through := hit.damage
			if playerTurn && defender.barrier > 0 {
				absorbed = minInt(hit.damage, int(math.Round(defender.barrier)))
				through = hit.damage - absorbed
				defender.barrier -= float64(absorbed)
			}
			defender.hp -= float64(through)
			crit := ""
			if hit.critical {
				if playerTurn {
					playerCritical = true
				}
				crit = " (Critical strike)"
				if !playerTurn {
					monsterCriticalDamage += hit.damage
				}
			}
			defenderHP := maxInt(0, int(math.Round(defender.hp)))
			defenderMaxHP := int(defender.maxHP)
			if playerTurn {
				attackResult := fmt.Sprintf("dealt %d damage%s", hit.damage, crit)
				if absorbed == hit.damage && absorbed > 0 {
					attackResult = fmt.Sprintf("dealt %d damage, fully absorbed by the barrier%s", hit.damage, crit)
				} else if absorbed > 0 {
					attackResult = fmt.Sprintf("dealt %d damage (%d absorbed, %d through)%s", hit.damage, absorbed, through, crit)
				}
				logs = append(logs, fmt.Sprintf("Turn %d: %s used %s attack and %s. %s HP: %d/%d", turn, actor, hit.element, attackResult, target, defenderHP, defenderMaxHP))
				if absorbed > 0 {
					logs = append(logs, fmt.Sprintf("Turn %d: The monster's barrier absorbs %d damage. Barrier HP: %d", turn, absorbed, maxInt(0, int(math.Round(defender.barrier)))))
					if defender.barrier <= 0 {
						logs = append(logs, fmt.Sprintf("Turn %d: The monster's barrier is destroyed!", turn))
					}
				}
				if playerTurn && corrupted > 0 {
					resistance := reduceResistance(&defender.stats, hit.element, corrupted)
					logs = append(logs, fmt.Sprintf("Turn %d: The monster's %s resistance is corrupted and decreases by %d%%. New resistance: %d%%", turn, hit.element, corrupted, resistance))
				}
			} else {
				logs = append(logs, fmt.Sprintf("Turn %d: %s used %s attack against %s and dealt %d damage%s. %s HP: %d/%d", turn, actor, hit.element, target, hit.damage, crit, target, defenderHP, defenderMaxHP))
			}
			if playerTurn && hit.critical && player.Stats.Lifesteal > 0 {
				lifesteal := int(math.Round(float64(hit.damage) * float64(player.Stats.Lifesteal) / 100))
				if lifesteal > 0 {
					before := p.hp
					p.hp = math.Min(p.maxHP, p.hp+float64(lifesteal))
					healed := int(math.Round(p.hp - before))
					logs = append(logs, fmt.Sprintf("Turn %d: Character_1 heals %d HP from lifesteal. HP: %d/%d", turn, healed, int(math.Round(p.hp)), player.Stats.HP))
				}
			}
		}
		if playerTurn && playerCritical && player.Stats.Frenzy > 0 {
			playerFrenzyNext = true
			logs = append(logs, fmt.Sprintf("Turn %d: Character_1's Frenzy triggers on critical. +%d%% damage will apply on each ally's next turn, active until the end of Character_1's next turn.", turn, player.Stats.Frenzy))
		}
		if !playerTurn && player.Stats.Greed > 0 {
			thresholds := int(math.Floor((p.maxHP - p.hp) / (p.maxHP * 0.10)))
			if thresholds > playerGreedStacks {
				playerGreedStacks = thresholds
				logs = append(logs, fmt.Sprintf("Turn %d: Greed empowers Character_1 (+%d%% damage, total +%d%%).", turn, player.Stats.Greed, player.Stats.Greed*playerGreedStacks))
			}
		}
		if !playerTurn && monsterCriticalDamage > 0 && monsterFrenzy > 0 {
			monsterFrenzyNext = true
			logs = append(logs, fmt.Sprintf("Turn %d: The monster's Frenzy triggers on critical. +%d%% damage will apply on its next attack.", turn, monsterFrenzy))
		}
		if playerTurn && playerFrenzyActive {
			p.damageFactor = 1
			playerFrenzyActive = false
		}
		if !playerTurn && monsterCriticalDamage > 0 && attacker.stats.Lifesteal > 0 {
			heal := int(math.Floor(float64(monsterCriticalDamage) * float64(attacker.stats.Lifesteal) / 100))
			before := m.hp
			m.hp = math.Min(m.maxHP, m.hp+float64(heal))
			healed := int(math.Round(m.hp - before))
			if healed > 0 {
				logs = append(logs, fmt.Sprintf("Turn %d: Monster heals %d HP from lifesteal. %s HP: %d/%d", turn, healed, monster.Name, int(math.Round(m.hp)), monster.Hp))
			}
		}
		if !playerTurn && player.Stats.EnchantedMirror > 0 && monsterTurns%3 == 1 && damage > 0 {
			reflected := int(math.Round(float64(damage) * float64(player.Stats.EnchantedMirror) / 100))
			m.hp -= float64(reflected)
			logs = append(logs, fmt.Sprintf("Turn %d: Character_1's Enchanted Mirror activates, dealing back %d damage to %s.", turn, reflected, monster.Name))
			logs = append(logs, fmt.Sprintf("Turn %d: %s takes %d reflected damage. %s HP: %d/%d", turn, monster.Name, reflected, monster.Name, maxInt(0, int(math.Round(m.hp))), monster.Hp))
		}
		if !playerTurn && monsterFrenzyActive {
			m.damageFactor /= 1 + float64(monsterFrenzy)/100
			monsterFrenzyActive = false
		}
		if playerTurn && greed > 0 {
			thresholds := int(math.Floor((m.maxHP - m.hp) / (m.maxHP * 0.10)))
			for greedStacks < thresholds {
				greedStacks++
				m.damageFactor = 1 + float64(greed*greedStacks)/100
				logs = append(logs, fmt.Sprintf("Turn %d: Greed empowers the monster (+%d%% damage, total +%d%%).", turn, greed, greed*greedStacks))
			}
		}
		if defender.hp <= 0 {
			result.Win = playerTurn
			result.Turns = turn
			break
		}
		result.Turns = turn
	}
	if result.Turns == MaxTurns && p.hp > 0 && m.hp > 0 {
		result.TimedOut = true
	}
	result.HPRemaining = minInt(naturalHP, maxInt(0, int(math.Round(p.hp))))
	if result.Win {
		logs = append(logs, fmt.Sprintf("Fight result: win. Character_1 HP: %d/%d vs %s HP: 0/%d", result.HPRemaining, naturalHP, monster.Name, monster.Hp))
	} else if !result.TimedOut {
		logs = append(logs, fmt.Sprintf("Fight result: loss. Character_1 HP: %d/%d vs %s HP: %d/%d", result.HPRemaining, naturalHP, monster.Name, maxInt(0, int(math.Round(m.hp))), monster.Hp))
	}
	result.Logs = logs
	return result
}

// decayBurnDamage mirrors the API's integer burn decay. The API rounds
// larger ticks, while the final low-damage ticks are truncated.
func decayBurnDamage(damage int) float64 {
	if damage <= 5 {
		return math.Floor(float64(damage) * 0.9)
	}
	return math.Round(float64(damage) * 0.9)
}

func monsterEffect(monster schemas.MonsterSchema, code string) int {
	if monster.Effects == nil {
		return 0
	}
	for _, effect := range *monster.Effects {
		if effect.Code == code {
			return effect.Value
		}
	}
	return 0
}

func isBoss(monster schemas.MonsterSchema) bool {
	return monster.Type == "boss" || monster.Type == "raid_boss"
}

func consumeUtility(utilities *[]Utility, antidoteOnly bool, player *combatant, maxHP int, logs *[]string, turn int) bool {
	for i := range *utilities {
		u := &(*utilities)[i]
		if u.Quantity <= 0 || (antidoteOnly && u.Antipoison == 0) || (!antidoteOnly && u.Restore == 0) {
			continue
		}
		if !antidoteOnly && player.hp > float64(maxHP)/2 {
			continue
		}
		u.Quantity--
		if antidoteOnly {
			*logs = append(*logs, fmt.Sprintf("Turn %d: Character_1 used %s and removed poison.", turn, u.Code))
		} else {
			player.hp = math.Min(float64(maxHP), player.hp+float64(u.Restore))
			*logs = append(*logs, fmt.Sprintf("Turn %d: Character_1 used %s and restored %d HP. Character_1 HP: %d/%d", turn, u.Code, u.Restore, int(math.Round(player.hp)), maxHP))
		}
		return antidoteOnly
	}
	return false
}

func consumeAntidote(utilities *[]Utility, logs *[]string, turn int) int {
	for i := range *utilities {
		u := &(*utilities)[i]
		if u.Quantity <= 0 || u.Antipoison == 0 {
			continue
		}
		u.Quantity--
		*logs = append(*logs, fmt.Sprintf("Turn %d: Character_1 used %s and removed %d poison.", turn, u.Code, u.Antipoison))
		return u.Antipoison
	}
	return 0
}

type combatant struct {
	hp, maxHP    float64
	stats        Stats
	damageFactor float64
	barrier      float64
}

type hit struct {
	element  string
	damage   int
	critical bool
}

func attackHits(attacker, defender Stats, rng RNG, sequence []bool, sequenceIndex *int) []hit {
	attacks := []struct {
		attack     int
		bonus      int
		element    string
		resistance int
	}{
		{
			attack:     attacker.AttackFire,
			bonus:      attacker.Dmg + attacker.DmgFire,
			element:    "fire",
			resistance: defender.ResFire,
		},
		{
			attack:     attacker.AttackEarth,
			bonus:      attacker.Dmg + attacker.DmgEarth,
			element:    "earth",
			resistance: defender.ResEarth,
		},
		{
			attack:     attacker.AttackWater,
			bonus:      attacker.Dmg + attacker.DmgWater,
			element:    "water",
			resistance: defender.ResWater,
		},
		{
			attack:     attacker.AttackAir,
			bonus:      attacker.Dmg + attacker.DmgAir,
			element:    "air",
			resistance: defender.ResAir,
		},
	}
	result := make([]hit, 0, len(attacks))
	// Critical chance is rolled once for the whole attack. An attack can
	// produce one hit per element, but the API marks all of those elemental
	// hits as critical when the attack's single critical roll succeeds.
	critical := false
	if sequenceIndex != nil && *sequenceIndex < len(sequence) {
		critical = sequence[*sequenceIndex]
		*sequenceIndex++
	} else {
		critical = rng != nil && rng() < float64(attacker.CriticalStrike)/100
	}
	for _, attack := range attacks {
		if attack.attack == 0 {
			continue
		}
		value := elemental(attack.attack, attack.bonus, attack.resistance)
		if value <= 0 {
			continue
		}
		if critical {
			value = int(math.Round(float64(value) * 1.5))
		}
		result = append(result, hit{
			element:  attack.element,
			damage:   value,
			critical: critical,
		})
	}
	return result
}

func elemental(attack, bonus, resistance int) int {
	if attack == 0 {
		return 0
	}
	raw := math.Round(float64(attack) * (1 + float64(bonus)/100))
	return int(math.Round(raw * math.Max(0, 1-float64(resistance)/100)))
}

func randomBubbleElement(rng RNG, previous int) int {
	if rng == nil {
		return 0
	}
	element := minInt(3, int(rng()*4))
	if element == previous {
		element = (element + 1) % 4
	}
	return element
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func reduceResistance(stats *Stats, element string, amount int) int {
	var resistance *int
	switch element {
	case "fire":
		resistance = &stats.ResFire
	case "earth":
		resistance = &stats.ResEarth
	case "water":
		resistance = &stats.ResWater
	case "air":
		resistance = &stats.ResAir
	default:
		return 0
	}
	*resistance -= amount
	return *resistance
}
