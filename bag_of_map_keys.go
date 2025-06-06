package views

import (
	"iter"
	"maps"
)

// BagOfMapKeys implements the Bag and List interfaces for a map's keys.
type BagOfMapKeys[K comparable, V any] struct {
	M map[K]V
}

// Convenience alias for BagOfMapKeys, which can be used as a List.
type ListOfMapKeys[K comparable, V any] = BagOfMapKeys[K, V]

// Len returns the number of entries in the map.
// This is O(1) in complexity.
func (b BagOfMapKeys[K, V]) Len() int {
	return len(b.M)
}

// Values returns a sequence of keys in the map.
func (b BagOfMapKeys[K, V]) Values() iter.Seq[K] {
	return maps.Keys(b.M)
}

// Has checks if the map contains the specified key.
// This is O(1) in complexity.
func (b BagOfMapKeys[K, V]) Has(value K) bool {
	if b.M == nil {
		return false
	}
	_, exists := b.M[value]
	return exists
}
