package views

import (
	"iter"
	"slices"
)

// BagOfSlice implements the Bag interface for a slice with comparable elements.
type BagOfSlice[T comparable] struct {
	S []T
}

// Len returns the number of elements in the slice.
// This is O(1) in complexity.
func (b BagOfSlice[T]) Len() int {
	return len(b.S)
}

// Values returns a sequence of values in the slice.
func (b BagOfSlice[T]) Values() iter.Seq[T] {
	return slices.Values(b.S)
}

// Has checks if the slice contains the specified value.
// This is O(n) in complexity.
func (b BagOfSlice[T]) Has(value T) bool {
	return slices.Contains(b.S, value)
}
