# Sorted Map

A sorted map implementation. Probably just one of the millions out there.  It uses a map, an `RWMutex` and a sorted slice underneath. It also uses binary search for insertion and deletion of the keys to make them performant.

## Examples

```go
key1, key2, key3 := "key1", "key2", "key3"
value1, value2, value2b, value3 := 1, 2, -2, 3

sm := sortedmap.New[string, int]().
    Set(key1, value1).
    Set(key2, value2).
    Set(key2, value2b).
    Set(key3, value3)

sm.Get(key2) // -2
sm.Len() // 3
sm.Has(key2) // true
sm.HasAll(key1, key2) // true
sm.HasAll(key1, key2, "nope") // false
sm.HasAny(key1, key2) // true
sm.HasAny(key1, key2, "nope") // true
sm.Delete(key1)
```

### Creating a sorted map with capacity

Creating a sorted map with capacity will mimic creating a map with a capacity, making it faster to insert until the given capacity, but using more memory at creation.

```go
key1, key2, key3 := "key1", "key2", "key3"
value1, value2, value2b, value3 := 1, 2, -2, 3

sm := sortedmap.NewWithCapacity[string, int](200).
    Set(key1, value1).
    Set(key2, value2).
    Set(key2, value2b).
    Set(key3, value3)

sm.Get(key2) // -2
sm.Len() // 3
```

### Creating a sorted map with an initial item

```go
key1, key2, key3 := "key1", "key2", "key3"
value1, value2, value2b, value3 := 1, 2, -2, 3

sm := sortedmap.NewFrom(key1, value1)

sm.Get(key1) // 1
sm.Len() // 1
```