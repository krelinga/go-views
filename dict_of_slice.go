package views

import (
	"iter"
	"slices"
)

// DictOfSlice implements the Dict interface for a slice, where keys are integers corresponding to indices.
type DictOfSlice[V any] struct {
	S []V
}

// Len returns the number of elements in the slice.
// This is O(1) in complexity.
func (d DictOfSlice[V]) Len() int {
	return len(d.S)
}

// Values returns the sequence of values in the slice.
func (d DictOfSlice[V]) Values() iter.Seq[V] {
	return slices.Values(d.S)
}

// Keys returns the sequence of keys (indices) in the slice.
func (d DictOfSlice[V]) Keys() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range d.S {
			if !yield(i) {
				return
			}
		}
	}
}

// Get retrieves the value associated with the specified key (index).
// If the key does not exist (index is out of range), it returns the zero value of V and false.
// This is O(1) in complexity.
func (d DictOfSlice[V]) Get(key int) (V, bool) {
	if !d.Has(key) {
		var zero V
		return zero, false
	}
	return d.S[key], true
}

// Has checks if the slice contains the specified key (index is valid).
// This is O(1) in complexity.
func (d DictOfSlice[V]) Has(key int) bool {
	if key < 0 || key >= len(d.S) {
		return false
	}
	return true
}

// All returns the sequence of all key-value pairs in the slice.
// The keys are the indices of the slice.
func (d DictOfSlice[V]) All() iter.Seq2[int, V] {
	return slices.All(d.S)
}