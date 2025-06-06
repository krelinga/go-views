package views

import (
	"iter"
	"maps"
)

// DictOfMap implements the Dict interface for a map.
type DictOfMap[K comparable, V any] struct {
	M map[K]V
}

// Len returns the number of entries in the map.
// This is O(1) in complexity.
func (d DictOfMap[K, V]) Len() int {
	return len(d.M)
}

// Values returns a sequence of values in the map.
func (d DictOfMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(d.M)
}

// Keys returns a sequence of keys in the map.
func (d DictOfMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(d.M)
}

// Get retrieves the value associated with the specified key.
// If the key does not exist, it returns the zero value of V and false.
// This is O(1) in complexity.
func (d DictOfMap[K, V]) Get(key K) (V, bool) {
	if d.M == nil {
		var zero V
		return zero, false
	}
	value, exists := d.M[key]
	return value, exists
}

// Has checks if the map contains the specified key.
func (d DictOfMap[K, V]) Has(key K) bool {
	if d.M == nil {
		return false
	}
	_, exists := d.M[key]
	return exists
}

// All returns the sequence of all key-value pairs in the map.
func (d DictOfMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(d.M)
}
