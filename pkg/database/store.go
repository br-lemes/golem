package database

import (
	"encoding/json"
	"fmt"
	"sync"
)

type store[T any, K comparable] struct {
	initOnce     func() *storeState[T, K]
	keyExtractor func(*T) K
}

type storeState[T any, K comparable] struct {
	items []*T
	byKey map[K]*T
}

func newStore[T any, K comparable](loader func() []T, keyExtractor func(*T) K) *store[T, K] {
	s := &store[T, K]{keyExtractor: keyExtractor}
	s.initOnce = sync.OnceValue(func() *storeState[T, K] {
		items := loader()
		itemsPtr := make([]*T, len(items))
		byKey := make(map[K]*T, len(items))
		for index := range items {
			itemPtr := &items[index]
			itemsPtr[index] = itemPtr
			byKey[keyExtractor(itemPtr)] = itemPtr
		}
		return &storeState[T, K]{items: itemsPtr, byKey: byKey}
	})
	return s
}

func jsonLoader[T any](data []byte) func() []T {
	return func() []T {
		var items []T
		err := json.Unmarshal(data, &items)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal embedded json: %v", err))
		}
		return items
	}
}

func (s *store[T, K]) All() []*T {
	state := s.initOnce()
	return state.items
}

func (s *store[T, K]) Get(key K) (*T, bool) {
	state := s.initOnce()
	item, ok := state.byKey[key]
	return item, ok
}

func (s *store[T, K]) Filter(predicate func(*T) bool) []*T {
	state := s.initOnce()
	var result []*T
	for _, item := range state.items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func (s *store[T, K]) Find(predicate func(*T) bool) (*T, bool) {
	state := s.initOnce()
	for _, item := range state.items {
		if predicate(item) {
			return item, true
		}
	}
	return nil, false
}
