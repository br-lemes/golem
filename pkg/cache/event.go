package cache

import (
	"encoding/json"
	"time"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

const activeEventsCacheDuration = 30 * time.Second

type activeEventsCache struct {
	FetchedAt time.Time                   `json:"fetched_at"`
	Events    []schemas.ActiveEventSchema `json:"events"`
}

func GetEvent(contentCode string) (*schemas.ActiveEventSchema, bool) {
	now := time.Now().UTC()
	var eventCache models.Cache
	if !findByName(&eventCache, "eventsActive") {
		return nil, true
	}

	var data activeEventsCache
	if json.Unmarshal([]byte(eventCache.Data), &data) != nil {
		return nil, true
	}
	for index := range data.Events {
		event := &data.Events[index]
		if event.Map.Interactions.Content == nil || event.Map.Interactions.Content.Code != contentCode {
			continue
		}
		if !event.Expiration.After(now) {
			return nil, true
		}
		return event, false
	}
	return nil, now.Sub(data.FetchedAt) >= activeEventsCacheDuration
}

func SaveActiveEvents(events []schemas.ActiveEventSchema) {
	if events == nil {
		events = []schemas.ActiveEventSchema{}
	}
	data, err := json.Marshal(activeEventsCache{
		FetchedAt: time.Now().UTC(),
		Events:    events,
	})
	if err != nil {
		return
	}
	cache.Save(&models.Cache{Name: "eventsActive", Data: string(data)})
}

func CleanActiveEvents() {
	cache.Where("name = ?", "eventsActive").Delete(&models.Cache{})
}
