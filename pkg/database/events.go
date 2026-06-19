package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

var (
	//go:embed events.json
	events []byte

	eventsList              []schemas.EventSchema
	eventsCache             map[string]schemas.EventSchema
	eventsContentCodesList  []string
	eventsContentCodesCache map[string]bool
)

func GetEvents() []schemas.EventSchema {
	initEventsCache()
	return eventsList
}

func GetEvent(code string) (schemas.EventSchema, bool) {
	initEventsCache()
	event, exists := eventsCache[code]
	return event, exists
}

func GetEventByContent(code string) (schemas.EventSchema, bool) {
	initEventsCache()
	for _, event := range eventsList {
		if event.Content.Code == code {
			return event, true
		}
	}
	return schemas.EventSchema{}, false
}

func GetEventCodes() []string {
	initEventsCache()
	var codes []string
	for code := range eventsCache {
		codes = append(codes, code)
	}
	return codes
}

func GetEventContentCodes() []string {
	initEventsContentCodesCache()
	return eventsContentCodesList
}

var initEventsCache = sync.OnceFunc(func() {
	eventsCache = make(map[string]schemas.EventSchema)
	err := json.Unmarshal(events, &eventsList)
	if err != nil {
		panic("failed to unmarshal events: " + err.Error())
	}
	for _, event := range eventsList {
		eventsCache[event.Code] = event
	}
})

var initEventsContentCodesCache = sync.OnceFunc(func() {
	initEventsCache()

	eventsContentCodesList = make([]string, 0)
	eventsContentCodesCache = make(map[string]bool)
	for _, event := range eventsList {
		code := event.Content.Code
		if !eventsContentCodesCache[code] {
			eventsContentCodesCache[code] = true
			eventsContentCodesList = append(eventsContentCodesList, code)
		}
	}
})
