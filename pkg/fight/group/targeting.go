package group

import "sort"

func (combat Combat) ActorOrder() []int {
	order := make([]int, 0, len(combat.Players)+1)
	for i, player := range combat.Players {
		if player.Alive && player.HP > 0 {
			order = append(order, i)
		}
	}
	order = append(order, -1)
	sort.SliceStable(order, func(i, j int) bool {
		left, right := combat.actorStats(order[i]), combat.actorStats(order[j])
		if left.Initiative != right.Initiative {
			return left.Initiative > right.Initiative
		}
		return left.HP > right.HP
	})
	return order
}

func (combat Combat) actorStats(index int) Stats {
	if index < 0 {
		return combat.Monster.Stats
	}
	return combat.Players[index].Fighter.Stats
}

func (combat Combat) TurnOrder() []int {
	order := make([]int, 0, len(combat.Players))
	for i, player := range combat.Players {
		if player.Alive && player.HP > 0 {
			order = append(order, i)
		}
	}
	sort.SliceStable(order, func(i, j int) bool {
		left, right := combat.Players[order[i]], combat.Players[order[j]]
		if left.Fighter.Stats.Initiative != right.Fighter.Stats.Initiative {
			return left.Fighter.Stats.Initiative > right.Fighter.Stats.Initiative
		}
		return left.MaxHP > right.MaxHP
	})
	return order
}

func (combat *Combat) SelectTarget() int {
	living := make([]int, 0, len(combat.Players))
	for i, player := range combat.Players {
		if player.Alive && player.HP > 0 {
			living = append(living, i)
		}
	}
	if len(living) == 0 {
		return -1
	}
	if combat.targetIndex < len(combat.TargetDecisions) {
		target := combat.TargetDecisions[combat.targetIndex]
		combat.targetIndex++
		for _, candidate := range living {
			if candidate == target {
				return target
			}
		}
	}
	candidates := living
	if combat.RNG() < 0.10 {
		lowest := combat.Players[living[0]].HP
		for _, i := range living[1:] {
			if combat.Players[i].HP < lowest {
				lowest = combat.Players[i].HP
			}
		}
		candidates = nil
		for _, i := range living {
			if combat.Players[i].HP == lowest {
				candidates = append(candidates, i)
			}
		}
	} else {
		highest := combat.Players[living[0]].Fighter.Stats.Threat
		for _, i := range living[1:] {
			if combat.Players[i].Fighter.Stats.Threat > highest {
				highest = combat.Players[i].Fighter.Stats.Threat
			}
		}
		candidates = nil
		for _, i := range living {
			if combat.Players[i].Fighter.Stats.Threat == highest {
				candidates = append(candidates, i)
			}
		}
	}
	lowest := combat.Players[candidates[0]].HP
	for _, i := range candidates[1:] {
		if combat.Players[i].HP < lowest {
			lowest = combat.Players[i].HP
		}
	}
	tied := candidates[:0]
	for _, i := range candidates {
		if combat.Players[i].HP == lowest {
			tied = append(tied, i)
		}
	}
	return tied[0]
}
