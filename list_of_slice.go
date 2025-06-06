package views

import (
	"iter"
	"slices"
)

type ListOfSlice[T any] struct {
	S []T
}

func (l ListOfSlice[T]) Len() int {
	return len(l.S)
}

func (l ListOfSlice[T]) Values() iter.Seq[T] {
	return slices.Values(l.S)
}
