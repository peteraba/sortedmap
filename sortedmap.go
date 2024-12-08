package sortedmap

import (
	"errors"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

func insertSorted[K constraints.Ordered](slice []K, value K) []K {
	i := sort.Search(
		len(slice),
		func(i int) bool {
			return slice[i] >= value
		},
	)

	slice = append(slice, value)
	copy(slice[i+1:], slice[i:])
	slice[i] = value

	return slice
}

func deleteSorted[K constraints.Ordered](slice []K, value K) []K {
	i := sort.Search(
		len(slice),
		func(i int) bool {
			return slice[i] >= value
		},
	)

	if i < len(slice) && slice[i] == value {
		copy(slice[i:], slice[i+1:])

		slice = slice[:len(slice)-1]
	}

	return slice
}

type SortedMap[K constraints.Ordered, T any] struct {
	mu         sync.RWMutex
	items      map[K]T
	sortedKeys []K
}

func NewSortedMap[K constraints.Ordered, T any]() *SortedMap[K, T] {
	return &SortedMap[K, T]{
		items:      make(map[K]T),
		sortedKeys: make([]K, 0),
	}
}

func (sm *SortedMap[K, T]) has(key K) bool {
	_, exists := sm.items[key]

	return exists
}

func (sm *SortedMap[K, T]) Set(key K, value T) *SortedMap[K, T] {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.has(key) {
		sm.sortedKeys = insertSorted(sm.sortedKeys, key)
	}

	sm.items[key] = value

	return sm
}

var ErrKeyDoesNotExist = errors.New("key does not exist")

func (sm *SortedMap[K, T]) Get(key K) (T, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	value, exists := sm.items[key]
	if !exists {
		return value, ErrKeyDoesNotExist
	}

	return value, nil
}

func (sm *SortedMap[K, T]) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.items) != len(sm.sortedKeys) {
		panic("sorted keys and items are out of sync")
	}

	return len(sm.items)
}

func (sm *SortedMap[K, T]) Has(key K) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.has(key)
}

func (sm *SortedMap[K, T]) HasAll(key ...K) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, k := range key {
		if _, exists := sm.items[k]; !exists {
			return false
		}
	}

	return true
}

func (sm *SortedMap[K, T]) HasAny(key ...K) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, k := range key {
		if _, exists := sm.items[k]; exists {
			return true
		}
	}

	return false
}

func (sm *SortedMap[K, T]) Delete(key K) *SortedMap[K, T] {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.has(key) {
		return sm
	}

	delete(sm.items, key)

	sm.sortedKeys = deleteSorted(sm.sortedKeys, key)

	return sm
}

func (sm *SortedMap[K, T]) Keys() []K {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.sortedKeys
}

func (sm *SortedMap[K, T]) Items() []T {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	keys := sm.Keys()
	values := make([]T, 0, len(keys))
	for _, key := range keys {
		values = append(values, sm.items[key])
	}

	return values
}
