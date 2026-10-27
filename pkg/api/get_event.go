package api

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func GetEvent(contentCode string) (*schemas.ActiveEventSchema, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.GetEvent(contentCode)
}

func (c *Client) GetEvent(contentCode string) (*schemas.ActiveEventSchema, error) {
	event, refresh := cache.GetEvent(contentCode)
	if event != nil {
		return event, nil
	}
	if !refresh {
		return nil, nil
	}

	events, err := c.EventsActive()
	if err != nil {
		return nil, err
	}
	cache.SaveActiveEvents(events)
	event, _ = cache.GetEvent(contentCode)
	return event, nil
}
