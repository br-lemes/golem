package utils

import (
	"math"

	"github.com/br-lemes/golem/pkg/schemas"
)

const maxTurnsTimeout = 100

// SimOptions controls a single Simulate pass.
type SimOptions struct {
	// Evaluates worst-case (0% player crit, 100% monster crit) to reject gear
	// that relies on lucky RNG.
	Pessimistic bool
	Rng         func() float64
}

type FightForecast struct {
	Win         bool
	Turns       int
	HpRemaining int
	HpLostPct   float64
	TimedOut    bool
	PlayerFirst bool
	Margin      float64
	MaxHp       int
}

type combatant struct {
	hp, maxHp      float64
	barrier        float64
	turn           int
	hitEV          float64
	hitNormal      float64
	hitCrit        float64
	dmgMultiplier  float64
	totalAtk       float64
	crit           float64
	init           int
	restore        float64
	healingPct     float64
	reconstitution int
	barrierRestore float64
	berserkerPct   float64
	berserked      bool
	lifestealPct   float64
	poison         float64
	burnCur        float64
	burnActive     bool
}

type elementHits struct {
	fire, earth, water, air float64
}

func (h elementHits) sum() float64 { return h.fire + h.earth + h.water + h.air }

func elementDamage(atk, dmgPct, resPct int) float64 {
	if atk == 0 {
		return 0
	}
	raw := math.Round(float64(atk) * (1 + float64(dmgPct)/100))
	factor := math.Max(0, 1-float64(resPct)/100)
	return math.Round(raw * factor)
}

func playerHits(stats EffectiveStats, boostDmg elementHits, monster schemas.MonsterSchema) elementHits {
	return elementHits{
		fire:  elementDamage(stats.AttackFire, stats.Dmg+stats.DmgFire+int(boostDmg.fire), monster.ResFire),
		earth: elementDamage(stats.AttackEarth, stats.Dmg+stats.DmgEarth+int(boostDmg.earth), monster.ResEarth),
		water: elementDamage(stats.AttackWater, stats.Dmg+stats.DmgWater+int(boostDmg.water), monster.ResWater),
		air:   elementDamage(stats.AttackAir, stats.Dmg+stats.DmgAir+int(boostDmg.air), monster.ResAir),
	}
}

func monsterHits(monster schemas.MonsterSchema, res elementHits) elementHits {
	return elementHits{
		fire:  elementDamage(monster.AttackFire, 0, int(res.fire)),
		earth: elementDamage(monster.AttackEarth, 0, int(res.earth)),
		water: elementDamage(monster.AttackWater, 0, int(res.water)),
		air:   elementDamage(monster.AttackAir, 0, int(res.air)),
	}
}

func opening(player Fighter) (pEff map[string]int, boostDmg, pRes elementHits, pMaxHp float64) {
	pEff = effMap(player.Effects)
	boostDmg = elementHits{
		fire:  float64(pEff["boost_dmg_fire"]),
		earth: float64(pEff["boost_dmg_earth"]),
		water: float64(pEff["boost_dmg_water"]),
		air:   float64(pEff["boost_dmg_air"]),
	}
	boostRes := elementHits{
		fire:  float64(pEff["boost_res_fire"]),
		earth: float64(pEff["boost_res_earth"]),
		water: float64(pEff["boost_res_water"]),
		air:   float64(pEff["boost_res_air"]),
	}
	pMaxHp = float64(player.Stats.Hp + pEff["boost_hp"])
	pRes = elementHits{
		fire:  float64(player.Stats.ResFire) + boostRes.fire,
		earth: float64(player.Stats.ResEarth) + boostRes.earth,
		water: float64(player.Stats.ResWater) + boostRes.water,
		air:   float64(player.Stats.ResAir) + boostRes.air,
	}
	return
}

func Simulate(player Fighter, monster schemas.MonsterSchema, opts SimOptions) FightForecast {
	pEff, boostDmg, pRes, pMaxHp := opening(player)
	mEff := effMap(monsterEffectEntries(monster))

	pCrit := math.Max(0, math.Min(1, float64(player.Stats.CriticalStrike)/100))
	mCrit := math.Max(0, math.Min(1, float64(monster.CriticalStrike)/100))
	if opts.Pessimistic {
		pCrit = 0
		mCrit = 1
	}

	pTotalAtk := float64(player.Stats.AttackFire + player.Stats.AttackEarth + player.Stats.AttackWater + player.Stats.AttackAir)
	mTotalAtk := float64(monster.AttackFire + monster.AttackEarth + monster.AttackWater + monster.AttackAir)

	pHits := playerHits(player.Stats, boostDmg, monster)
	mHits := monsterHits(monster, pRes)

	P := &combatant{
		hp:             pMaxHp,
		maxHp:          pMaxHp,
		barrier:        float64(pEff["barrier"]),
		hitEV:          pHits.sum() * (1 + 0.5*pCrit),
		hitNormal:      pHits.sum(),
		hitCrit:        math.Round(pHits.sum() * 1.5),
		dmgMultiplier:  1,
		totalAtk:       pTotalAtk,
		crit:           pCrit,
		init:           player.Stats.Initiative,
		restore:        float64(pEff["restore"]),
		healingPct:     float64(pEff["healing"]),
		reconstitution: pEff["reconstitution"],
		barrierRestore: float64(pEff["barrier"]),
		berserkerPct:   float64(pEff["berserker_rage"]),
		lifestealPct:   float64(pEff["lifesteal"]),
		poison:         float64(mEff["poison"]),
		burnCur:        float64(mEff["burn"]) / 100 * mTotalAtk,
		burnActive:     mEff["burn"] > 0,
	}
	M := &combatant{
		hp:             float64(monster.Hp),
		maxHp:          float64(monster.Hp),
		barrier:        float64(mEff["barrier"]),
		hitEV:          mHits.sum() * (1 + 0.5*mCrit),
		hitNormal:      mHits.sum(),
		hitCrit:        math.Round(mHits.sum() * 1.5),
		dmgMultiplier:  1,
		totalAtk:       mTotalAtk,
		crit:           mCrit,
		init:           monster.Initiative,
		restore:        float64(mEff["restore"]),
		healingPct:     float64(mEff["healing"]),
		reconstitution: mEff["reconstitution"],
		barrierRestore: float64(mEff["barrier"]),
		berserkerPct:   float64(mEff["berserker_rage"]),
		lifestealPct:   float64(mEff["lifesteal"]),
		poison:         float64(pEff["poison"]),
		burnCur:        float64(pEff["burn"]) / 100 * pTotalAtk,
		burnActive:     pEff["burn"] > 0,
	}

	playerFirst := P.init >= M.init

	checkBerserk := func(c *combatant) {
		if c.berserkerPct > 0 && !c.berserked && c.hp < 0.25*c.maxHp {
			c.dmgMultiplier *= 1 + c.berserkerPct/100
			c.berserked = true
		}
	}

	startOfTurn := func(c *combatant) bool {
		c.turn++
		if c.reconstitution > 0 && c.turn%c.reconstitution == 0 {
			c.hp = c.maxHp
		}
		if c.healingPct > 0 && c.turn%3 == 0 {
			c.hp = math.Min(c.maxHp, c.hp+c.maxHp*c.healingPct/100)
		}
		if c.barrierRestore > 0 && c.turn > 1 && (c.turn-1)%5 == 0 {
			c.barrier += c.barrierRestore
		}
		if c.restore > 0 && c.hp < 0.5*c.maxHp {
			c.hp = math.Min(c.maxHp, c.hp+c.restore)
		}
		if c.poison > 0 {
			c.hp -= c.poison
		}
		if c.burnActive {
			c.hp -= c.burnCur
			c.burnCur *= 0.9
			if c.burnCur < 1 {
				c.burnActive = false
			}
		}
		checkBerserk(c)
		return c.hp > 0
	}

	attack := func(a, d *combatant) bool {
		isCrit := false
		dmg := a.hitEV
		if opts.Rng != nil {
			isCrit = opts.Rng() < a.crit
			dmg = a.hitNormal
			if isCrit {
				dmg = a.hitCrit
			}
		}
		dmg *= a.dmgMultiplier
		if d.barrier > 0 {
			absorbed := math.Min(d.barrier, dmg)
			d.barrier -= absorbed
			dmg -= absorbed
		}
		d.hp -= dmg
		if a.lifestealPct > 0 {
			share := a.crit
			if opts.Rng != nil {
				share = 0
				if isCrit {
					share = 1
				}
			}
			a.hp = math.Min(a.maxHp, a.hp+share*(a.lifestealPct/100)*a.totalAtk)
		}
		checkBerserk(d)
		return d.hp <= 0
	}

	win := false
	timedOut := false
	turn := 0
	for turn = 1; turn <= maxTurnsTimeout; turn++ {
		actorIsPlayer := (turn%2 == 1) == playerFirst
		actor, defender := M, P
		if actorIsPlayer {
			actor, defender = P, M
		}
		if !startOfTurn(actor) {
			win = !actorIsPlayer
			break
		}
		if attack(actor, defender) {
			win = actorIsPlayer
			break
		}
	}
	if turn > maxTurnsTimeout {
		win = false
		timedOut = true
		turn = maxTurnsTimeout
	}

	clampedHp := int(math.Max(0, math.Round(P.hp)))
	lost := int(pMaxHp) - clampedHp
	hpLostPct := 0.0
	if pMaxHp > 0 {
		hpLostPct = math.Max(0, math.Min(1, float64(lost)/pMaxHp))
	}
	margin := 0.0
	if win && pMaxHp > 0 {
		margin = float64(clampedHp) / pMaxHp
	}
	return FightForecast{
		Win:         win,
		Turns:       turn,
		HpRemaining: clampedHp,
		HpLostPct:   hpLostPct,
		TimedOut:    timedOut,
		PlayerFirst: playerFirst,
		Margin:      margin,
		MaxHp:       int(pMaxHp),
	}
}
