package database

import (
	"sync"
)

type view[T any, K comparable] struct {
	initOnce     func() *storeState[T, K]
	keyExtractor func(*T) K
}

func (s *store[T, K]) View(predicate func(*T) bool) *view[T, K] {
	v := &view[T, K]{keyExtractor: s.keyExtractor}
	v.initOnce = sync.OnceValue(func() *storeState[T, K] {
		items := s.All()
		viewItems := make([]*T, 0, len(items))
		byKey := make(map[K]*T)
		for _, item := range items {
			if !predicate(item) {
				continue
			}
			key := s.keyExtractor(item)
			viewItems = append(viewItems, item)
			byKey[key] = item
		}
		return &storeState[T, K]{items: viewItems, byKey: byKey}
	})
	return v
}

func (v *view[T, K]) All() []*T {
	state := v.initOnce()
	return state.items
}

func (v *view[T, K]) Get(key K) (*T, bool) {
	state := v.initOnce()
	item, ok := state.byKey[key]
	return item, ok
}

func (v *view[T, K]) Filter(predicate func(*T) bool) []*T {
	state := v.initOnce()
	var result []*T
	for _, item := range state.items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func (v *view[T, K]) Find(predicate func(*T) bool) (*T, bool) {
	state := v.initOnce()
	for _, item := range state.items {
		if predicate(item) {
			return item, true
		}
	}
	return nil, false
}

func (v *view[T, K]) Keys() []K {
	state := v.initOnce()
	keys := make([]K, len(state.items))
	for i, item := range state.items {
		keys[i] = v.keyExtractor(item)
	}
	return keys
}
