package views

import (
	"iter"
	"slices"
)

type Slice[T any] interface {
	Len() int
	At(index int) T
	All() iter.Seq2[int, T]
	Values() iter.Seq[T]
}

// NewSlice returns a read-only [Slice] view of s.
//
// S may be []T or any type whose underlying type is []T.
//
// A nil s yields a nil [Slice], so nil-ness round-trips rather than collapsing
// into an empty view; an empty but non-nil s yields a non-nil view of length
// zero. A nil view is a nil interface, so callers must nil-check the result
// before calling methods on it.
//
// The view aliases s rather than copying it, so it prevents modification
// through the view only; callers that retain s can still change the elements
// the view exposes.
func NewSlice[S ~[]T, T any](s S) Slice[T] {
	if s == nil {
		return nil
	}
	return sliceView[T]{s: s}
}

// sliceView is the [Slice] implementation returned by [NewSlice].
type sliceView[T any] struct {
	s []T
}

func (v sliceView[T]) Len() int {
	return len(v.s)
}

func (v sliceView[T]) At(index int) T {
	return v.s[index]
}

func (v sliceView[T]) All() iter.Seq2[int, T] {
	return slices.All(v.s)
}

func (v sliceView[T]) Values() iter.Seq[T] {
	return slices.Values(v.s)
}

// NewSliceVals returns a read-only [Slice] view of s whose elements are
// converted by conv.
//
// S may be []F or any type whose underlying type is []F.
//
// A nil s yields a nil [Slice], so nil-ness round-trips rather than collapsing
// into an empty view; an empty but non-nil s yields a non-nil view of length
// zero. A nil view is a nil interface, so callers must nil-check the result
// before calling methods on it.
//
// conv is applied lazily, once per element read, so a value read twice is
// converted twice and elements that are never read are never converted. Like
// [NewSlice], the view aliases s.
func NewSliceVals[S ~[]F, F, T any](s S, conv func(F) T) Slice[T] {
	if s == nil {
		return nil
	}
	return sliceValsView[F, T]{s: s, conv: conv}
}

// sliceValsView is the [Slice] implementation returned by [NewSliceVals]. It
// holds the conversion function alongside the backing slice and applies it at
// each read.
type sliceValsView[F, T any] struct {
	s    []F
	conv func(F) T
}

func (v sliceValsView[F, T]) Len() int {
	return len(v.s)
}

func (v sliceValsView[F, T]) At(index int) T {
	return v.conv(v.s[index])
}

func (v sliceValsView[F, T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, e := range v.s {
			if !yield(i, v.conv(e)) {
				return
			}
		}
	}
}

func (v sliceValsView[F, T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range v.s {
			if !yield(v.conv(e)) {
				return
			}
		}
	}
}
