package views

import (
	"iter"
	"maps"
)

// BagOfMapValues implements the Bag and List interfaces for a map's values.
// It requires that the map's values are comparable, which is necessary to implement the Has method.
type BagOfMapValues[K comparable, V comparable] struct {
	M map[K]V
}

// Len returns the number of entries in the map.
// This is O(1) in complexity.
func (b BagOfMapValues[K, V]) Len() int {
	return len(b.M)
}

// Values returns a sequence of values in the map.
func (b BagOfMapValues[K, V]) Values() iter.Seq[V] {
	return maps.Values(b.M)
}

// Has checks if the map contains the specified value.
// This is O(n) in complexity.
func (b BagOfMapValues[K, V]) Has(value V) bool {
	if b.M == nil {
		return false
	}
	for _, v := range b.M {
		if v == value {
			return true
		}
	}
	return false
}
