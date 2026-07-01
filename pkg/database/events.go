package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/tidwall/gjson"
)

var (
	//go:embed events.json
	events []byte

	eventsList             []schemas.EventSchema
	eventsCache            map[string]schemas.EventSchema
	eventCodesList         []string
	eventsContentCodesList []string
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
	initEventCodes()
	return eventCodesList
}

func GetEventContentCodes() []string {
	initEventContentCodes()
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

var initEventCodes = sync.OnceFunc(func() {
	gjson.GetBytes(events, "#.code").ForEach(func(_, value gjson.Result) bool {
		eventCodesList = append(eventCodesList, value.String())
		return true
	})
})

var initEventContentCodes = sync.OnceFunc(func() {
	uniqueContentCodes := make(map[string]struct{})
	gjson.GetBytes(events, "#.content.code").ForEach(func(_, value gjson.Result) bool {
		s := value.String()
		if s != "" {
			uniqueContentCodes[s] = struct{}{}
		}
		return true
	})
	for code := range uniqueContentCodes {
		eventsContentCodesList = append(eventsContentCodesList, code)
	}
})
