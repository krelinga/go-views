package views

import (
	"iter"
	"maps"
)

type ListOfMapValues[K comparable, V any] struct {
	M map[K]V
}

func (l ListOfMapValues[K, V]) Len() int {
	return len(l.M)
}

func (l ListOfMapValues[K, V]) Values() iter.Seq[V] {
	return maps.Values(l.M)
}
