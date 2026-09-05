package views

import (
	"iter"
)

type SliceWrapper[U, E any] interface {
	Wrap(U) E
}

type Slice[U, E any, C SliceWrapper[U, E]] []U

func (s Slice[U, E, C]) Len() int {
	return len(s)
}

func (s Slice[U, E, C]) At(index int) E {
	var u C
	return u.Wrap(s[index])
}

func (s Slice[U, E, C]) All() iter.Seq2[int, E] {
	var u C
	return func(yield func(int, E) bool) {
		for i, v := range s {
			if !yield(i, u.Wrap(v)) {
				break
			}
		}
	}
}

func (s Slice[U, E, C]) Values() iter.Seq[E] {
	var u C
	return func(yield func(E) bool) {
		for _, v := range s {
			if !yield(u.Wrap(v)) {
				break
			}
		}
	}
}

func (s Slice[U, E, C]) Slice(start, end int) Slice[U, E, C] {
	return s[start:end]
}