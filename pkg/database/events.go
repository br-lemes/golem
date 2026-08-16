package database

import (
	_ "embed"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed events.json
var events []byte

var Events = newStore(jsonLoader[schemas.EventSchema](events), func(event *schemas.EventSchema) string {
	return event.Code
})

var eventContentCodes = sync.OnceValue(func() []string {
	seen := make(map[string]struct{})
	var codes []string

	for _, event := range Events.All() {
		code := event.Content.Code
		if code == "" {
			//+gocover:ignore:block all event contents have codes
			continue
		}
		_, exists := seen[code]
		if exists {
			//+gocover:ignore:block event content codes are unique
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
	}
	return codes
})

func EventContentCodes() []string {
	return eventContentCodes()
}
