package damage

import (
	"math"

	"github.com/br-lemes/golem/pkg/schemas"
)

func damage(attack, gear bonus, monster schemas.MonsterSchema) int {
	total := 0
	total += elementDamage(attack[bonusFire], gear[bonusDmg]+gear[bonusFire], monster.ResFire)
	total += elementDamage(attack[bonusWater], gear[bonusDmg]+gear[bonusWater], monster.ResWater)
	total += elementDamage(attack[bonusEarth], gear[bonusDmg]+gear[bonusEarth], monster.ResEarth)
	total += elementDamage(attack[bonusAir], gear[bonusDmg]+gear[bonusAir], monster.ResAir)
	return total
}

func elementDamage(attack, bonus, resistance int) int {
	if attack <= 0 {
		return 0
	}
	raw := math.Round(float64(attack) * (1 + float64(bonus)/100))
	return int(math.Round(raw * resistanceFactor(resistance)))
}

func resistanceFactor(resistance int) float64 {
	return math.Max(0, 1-float64(resistance)/100)
}
