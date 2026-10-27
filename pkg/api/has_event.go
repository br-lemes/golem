package api

import "github.com/br-lemes/golem/pkg/cache"

func HasEvent(contentCode string) (bool, error) {
	//+gocover:ignore:block public compatibility wrapper
	return defaultClient.HasEvent(contentCode)
}

func (c *Client) HasEvent(contentCode string) (bool, error) {
	found, refresh := cache.HasEvent(contentCode)
	if found {
		return true, nil
	}
	if !refresh {
		return false, nil
	}

	events, err := c.EventsActive()
	if err != nil {
		return false, err
	}
	cache.SaveActiveEvents(events)
	found, _ = cache.HasEvent(contentCode)
	return found, nil
}
