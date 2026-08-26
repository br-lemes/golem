package fight

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func IterationChunks(iterations int) []int {
	chunks := []int{}
	for iterations > 100 {
		chunks = append(chunks, 100)
		iterations -= 100
	}
	return append(chunks, iterations)
}

func ValidateRequest(request schemas.CombatSimulationRequestSchema) error {
	if len(request.Characters) != 1 {
		return fmt.Errorf("characters must contain exactly one entry")
	}
	if request.Monster == "" {
		return fmt.Errorf("monster is required")
	}
	if request.Iterations < 1 {
		return fmt.Errorf("iterations must be at least 1")
	}
	if request.Characters[0].Level < 1 || request.Characters[0].Level > 50 {
		return fmt.Errorf("characters[0].level must be between 1 and 50")
	}
	return nil
}
