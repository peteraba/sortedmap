package sortedmap

import (
	"errors"
	"sort"
	"sync"
)

func insertSorted(slice []string, value string) []string {
	i := sort.SearchStrings(slice, value)

	slice = append(slice, "")
	copy(slice[i+1:], slice[i:])
	slice[i] = value

	return slice
}

func deleteSorted(slice []string, value string) []string {
	i := sort.SearchStrings(slice, value)

	if i < len(slice) && slice[i] == value {
		copy(slice[i:], slice[i+1:])

		slice = slice[:len(slice)-1]
	}

	return slice
}

type SortedMap[T any] struct {
	mu         sync.RWMutex
	items      map[string]T
	sortedKeys []string
}

func NewSortedMap[T any]() *SortedMap[T] {
	return &SortedMap[T]{
		items:      make(map[string]T),
		sortedKeys: make([]string, 0),
	}
}

func (sm *SortedMap[T]) Set(key string, value T) *SortedMap[T] {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.has(key) {
		sm.sortedKeys = insertSorted(sm.sortedKeys, key)
	}

	sm.items[key] = value

	return sm
}

var ErrKeyDoesNotExist = errors.New("key does not exist")

func (sm *SortedMap[T]) Get(key string) (T, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	value, exists := sm.items[key]
	if !exists {
		return value, ErrKeyDoesNotExist
	}

	return value, nil
}

func (sm *SortedMap[T]) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if len(sm.items) != len(sm.sortedKeys) {
		panic("sorted keys and items are out of sync")
	}

	return len(sm.items)
}

func (sm *SortedMap[T]) has(key string) bool {
	_, exists := sm.items[key]

	return exists
}

func (sm *SortedMap[T]) Has(key string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.has(key)
}

func (sm *SortedMap[T]) HasAll(key ...string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, k := range key {
		if _, exists := sm.items[k]; !exists {
			return false
		}
	}

	return true
}

func (sm *SortedMap[T]) HasAny(key ...string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for _, k := range key {
		if _, exists := sm.items[k]; exists {
			return true
		}
	}

	return false
}

func (sm *SortedMap[T]) Delete(key string) *SortedMap[T] {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.has(key) {
		return sm
	}

	delete(sm.items, key)

	sm.sortedKeys = deleteSorted(sm.sortedKeys, key)

	return sm
}

func (sm *SortedMap[T]) Keys() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.sortedKeys
}

func (sm *SortedMap[T]) Items() []T {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	keys := sm.Keys()
	values := make([]T, 0, len(keys))
	for _, key := range keys {
		values = append(values, sm.items[key])
	}

	return values
}
