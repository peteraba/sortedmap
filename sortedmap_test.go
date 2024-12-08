package sortedmap_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/peteraba/sortedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortedMap_SetHasAndGet(t *testing.T) {
	key1 := "key1"
	value1 := "value1"

	sm := sortedmap.New[string, string]().
		Set(key1, value1)
	require.True(t, sm.Has(key1))

	actualValue, err := sm.Get("key1")
	require.NoError(t, err)

	assert.Equal(t, value1, actualValue)
}

func TestSortedMap_HasGetNonExistentKey(t *testing.T) {
	key1 := "key1"

	sm := sortedmap.New[string, string]()

	_, err := sm.Get(key1)
	require.Error(t, err)
	require.False(t, sm.Has(key1))

	assert.ErrorIs(t, sortedmap.ErrKeyDoesNotExist, err)
}

func TestSortedMap_HasAll(t *testing.T) {
	key1, key2, key3 := "key1", "key2", "key3"
	value1, value2 := "value1", "value2"

	sm := sortedmap.New[string, string]().
		Set(key1, value1).
		Set(key2, value2)

	assert.True(t, sm.HasAll(key1, key2))
	assert.False(t, sm.HasAll(key1, key2, key3))
}

func TestSortedMap_HasAny(t *testing.T) {
	key1, key2, key3 := "key1", "key2", "key3"
	value1 := "value1"

	sm := sortedmap.New[string, string]().
		Set(key1, value1)

	assert.True(t, sm.HasAny(key1, key2, key3))
	assert.False(t, sm.HasAny(key2, key3))
}

func TestSortedMap_Delete(t *testing.T) {
	key1, key2, key3 := "key1", "key2", "key3"
	value1, value2, value3 := "value1", "value2", "value3"

	sm := sortedmap.New[string, string]().
		Set(key2, value2).
		Set(key1, value1).
		Set(key3, value3)

	assert.True(t, sm.Has(key1))
	assert.True(t, sm.Has(key2))
	assert.True(t, sm.Has(key3))

	sm.Delete(key1, key2)

	assert.False(t, sm.Has(key1))
	assert.False(t, sm.Has(key2))
	assert.True(t, sm.Has(key3))
}

func TestSortedMap_Keys(t *testing.T) {
	key1, key2, key3 := "key1", "key2", "key3"
	value1, value2, value3 := "value1", "value2", "value3"
	expectedKeys := []string{key1, key2, key3}

	sm := sortedmap.New[string, string]().
		Set(key2, value2).
		Set(key1, value1).
		Set(key3, value3)

	actual := sm.Keys()
	assert.Equal(t, expectedKeys, actual)
}

func TestSortedMap_Items(t *testing.T) {
	key1, key2, key3 := "key1", "key2", "key3"
	value1, value2, value3 := "value1", "value2", "value3"

	sm := sortedmap.New[string, string]().
		Set(key2, value2).
		Set(key1, value1).
		Set(key3, value3)

	expectedKeys := []string{key1, key2, key3}
	expectedValues := []string{value1, value2, value3}

	actualKeys := make([]string, 0, 3)
	actualValues := make([]string, 0, 3)
	for key, value := range sm.Items() {
		actualKeys = append(actualKeys, key)
		actualValues = append(actualValues, value)
	}

	assert.Equal(t, expectedKeys, actualKeys)
	assert.Equal(t, expectedValues, actualValues)
}

func TestSortedMap_Complex(t *testing.T) {
	key1, key2, key2b, key3 := "key1", "key2", "key2", "key3"
	value1, value2, value2b, value3 := 1, 2, -2, 3

	sm := sortedmap.New[string, int]().
		Set(key1, value1).
		Set(key2, value2).
		Set(key2b, value2b).
		Set(key3, value3)

	assert.Equal(t, 3, sm.Len())

	actualValue, err := sm.Get(key2)
	require.NoError(t, err)

	sm.Delete(key2)

	assert.Equal(t, 2, sm.Len())

	assert.Equal(t, value2b, actualValue)
}

func TestSortedMap_ParallelSet(t *testing.T) {
	key1, key2, key2b, key3 := "key1", "key2", "key2", "key3"
	value1, value2, value2b, value3 := 1, 2, -2, 3

	sm := sortedmap.New[string, int]()

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Set(key1, value1).Set(key2, value2)
	}()

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Set(key2b, value2b).Set(key3, value3)
	}()

	time.Sleep(time.Millisecond)

	assert.Equal(t, 3, sm.Len())
	assert.Equal(t, []string{key1, key2, key3}, sm.Keys())
}

func TestSortedMap_ParallelSetFloatFloat(t *testing.T) {
	key1, key2, key3 := 1.23, 2.34, 3.45
	value1, value2, value2b, value3 := 8.34, 9.45, -9.45, 10.56

	sm := sortedmap.New[float64, float64]()

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Set(key1, value1)

		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Set(key2, value2)
	}()

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Set(key3, value3)

		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Set(key2, value2b)
	}()

	time.Sleep(time.Millisecond)

	assert.Equal(t, 3, sm.Len())
	assert.Equal(t, []float64{key1, key2, key3}, sm.Keys())
	assert.Equal(t, value1, sm.MustGet(key1))
	assert.Equal(t, value3, sm.MustGet(key3))
}

func TestSortedMap_ParallelDelete(t *testing.T) {
	key1, key2, key2b, key3 := "key1", "key2", "key2", "key3"
	value1, value2, value2b, value3 := 1, 2, -2, 3

	sm := sortedmap.New[string, int]().
		Set(key1, value1).
		Set(key2, value2).
		Set(key2b, value2b).
		Set(key3, value3)

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Delete(key1)
	}()

	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)

		sm.Delete(key3)
	}()

	time.Sleep(time.Millisecond)

	assert.Equal(t, 1, sm.Len())
	assert.Equal(t, []string{key2}, sm.Keys())
}

func TestSortedMap_NewWithCapacity(t *testing.T) {
	sm := sortedmap.NewWithCapacity[string, string](200)

	assert.Equal(t, 0, sm.Len())
	assert.Equal(t, 0, len(sm.Keys()))
	assert.Equal(t, 0, len(sm.Values()))
}

func TestSortedMap_NewFrom(t *testing.T) {
	key1, value1 := "key1", "value1"

	sm := sortedmap.NewFrom(key1, value1)

	assert.Equal(t, 1, sm.Len())
	assert.Equal(t, []string{key1}, sm.Keys())
	assert.Equal(t, []string{value1}, sm.Values())
}

type A struct {
	B []string
	C []int
}

func (a *A) AddB(newB string) *A {
	a.B = append(a.B, newB)

	return a
}

func TestSortedMap_NewFromComplex(t *testing.T) {
	key1 := "key1"
	value1 := &A{
		B: []string{"b1", "b2"},
		C: []int{1, 2},
	}

	sm := sortedmap.NewFrom(key1, value1)

	sm.Set(key1, sm.MustGet(key1).AddB("b3"))

	assert.Equal(t, 1, sm.Len())

	actualValue := sm.MustGet(key1)
	assert.Equal(t, []string{"b1", "b2", "b3"}, actualValue.B)
	assert.Equal(t, []int{1, 2}, actualValue.C)
}
