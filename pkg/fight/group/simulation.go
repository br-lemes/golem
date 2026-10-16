package group

const MaxTurns = 100

type SimulationOptions struct {
	Iterations int
	RNG        RNG
}

type Summary struct {
	Results []Result
	Wins    int
	Losses  int
	Winrate float32
}

type SimulationMetrics struct {
	AverageTurns   float32
	AverageFinalHP float32
}

func SimulateMany(fighters []Fighter, monster Monster, options SimulationOptions) Summary {
	if options.Iterations < 1 {
		options.Iterations = 1
	}
	summary := Summary{Results: make([]Result, 0, options.Iterations)}
	for range options.Iterations {
		combat := NewCombat(fighters, monster, options.RNG)
		result := combat.Run(MaxTurns)
		summary.Results = append(summary.Results, result)
		if result.Win {
			summary.Wins++
		} else {
			summary.Losses++
		}
	}
	summary.Winrate = float32(summary.Wins) * 100 / float32(options.Iterations)
	return summary
}

func Metrics(results []Result) SimulationMetrics {
	metrics := SimulationMetrics{}
	for _, result := range results {
		metrics.AverageTurns += float32(result.Turns)
		for _, hp := range result.FinalHP {
			metrics.AverageFinalHP += float32(hp)
		}
	}
	if len(results) > 0 {
		metrics.AverageTurns /= float32(len(results))
		metrics.AverageFinalHP /= float32(len(results))
	}
	return metrics
}
