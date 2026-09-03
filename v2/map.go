package views

import (
	"iter"
	"maps"
)

type Map[K, V any] interface {
	Len() int
	Get(key K) (V, bool)
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	All() iter.Seq2[K, V]
}

// NewMap returns a read-only [Map] view of m.
//
// M may be map[K]V or any type whose underlying type is map[K]V.
//
// A nil m yields a nil [Map], so nil-ness round-trips rather than collapsing
// into an empty view; an empty but non-nil m yields a non-nil view of length
// zero. A nil view is a nil interface, so callers must nil-check the result
// before calling methods on it.
//
// The view aliases m rather than copying it, so it prevents modification
// through the view only; callers that retain m can still change the entries
// the view exposes.
func NewMap[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	if m == nil {
		return nil
	}
	return mapView[K, V]{m: m}
}

// mapView is the [Map] implementation returned by [NewMap].
type mapView[K comparable, V any] struct {
	m map[K]V
}

func (v mapView[K, V]) Len() int {
	return len(v.m)
}

func (v mapView[K, V]) Get(key K) (V, bool) {
	value, ok := v.m[key]
	return value, ok
}

func (v mapView[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(v.m)
}

func (v mapView[K, V]) Values() iter.Seq[V] {
	return maps.Values(v.m)
}

func (v mapView[K, V]) All() iter.Seq2[K, V] {
	return maps.All(v.m)
}

// NewMapVals returns a read-only [Map] view of m whose values are converted by
// conv. Keys are exposed unchanged.
//
// M may be map[K]F or any type whose underlying type is map[K]F.
//
// A nil m yields a nil [Map], so nil-ness round-trips rather than collapsing
// into an empty view; an empty but non-nil m yields a non-nil view of length
// zero. A nil view is a nil interface, so callers must nil-check the result
// before calling methods on it.
//
// conv is applied lazily, once per value read, so a value read twice is
// converted twice and entries that are never read are never converted. Like
// [NewMap], the view aliases m.
func NewMapVals[M ~map[K]F, K comparable, F, V any](m M, conv func(F) V) Map[K, V] {
	if m == nil {
		return nil
	}
	return mapValsView[K, F, V]{m: m, conv: conv}
}

// mapValsView is the [Map] implementation returned by [NewMapVals]. It holds the
// conversion function alongside the backing map and applies it at each read.
type mapValsView[K comparable, F, V any] struct {
	m    map[K]F
	conv func(F) V
}

func (v mapValsView[K, F, V]) Len() int {
	return len(v.m)
}

func (v mapValsView[K, F, V]) Get(key K) (V, bool) {
	from, ok := v.m[key]
	if !ok {
		var zero V
		return zero, false
	}
	return v.conv(from), true
}

func (v mapValsView[K, F, V]) Keys() iter.Seq[K] {
	return maps.Keys(v.m)
}

func (v mapValsView[K, F, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, from := range v.m {
			if !yield(v.conv(from)) {
				return
			}
		}
	}
}

func (v mapValsView[K, F, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, from := range v.m {
			if !yield(key, v.conv(from)) {
				return
			}
		}
	}
}

// NewMapKeyValues returns a read-only [Map] view of m whose entries are
// converted by conv, which sees each key and value together and may derive the
// exposed value from both.
//
// M may be map[FK]FV or any type whose underlying type is map[FK]FV.
//
// A nil m yields a nil [Map], so nil-ness round-trips rather than collapsing
// into an empty view; an empty but non-nil m yields a non-nil view of length
// zero. A nil view is a nil interface, so callers must nil-check the result
// before calling methods on it.
//
// unconv recovers a backing key from an exposed one, which is what makes
// [Map.Get] possible: Get has only the converted key to work with and must
// map it back to look anything up. It reports false for an exposed key with no
// backing counterpart — a converted key space is often larger than the
// original — in which case Get reports false rather than guessing at a zero key.
//
// conv and unconv must agree on keys: unconv(k) must recover the key that conv
// turned into k. The view does not verify this, and it does not require conv to
// be injective on keys — if two backing keys convert to the same exposed key,
// Len still counts backing entries and iteration still yields the duplicate.
//
// conv is applied lazily, once per entry read, so an entry read twice is
// converted twice and entries that are never read are never converted. Like
// [NewMap], the view aliases m.
func NewMapKeyValues[M ~map[FK]FV, FK comparable, K, FV, V any](
	m M,
	conv func(FK, FV) (K, V),
	unconv func(K) (FK, bool),
) Map[K, V] {
	if m == nil {
		return nil
	}
	return mapKeyValuesView[FK, FV, K, V]{m: m, conv: conv, unconv: unconv}
}

// mapKeyValuesView is the [Map] implementation returned by [NewMapKeyValues]. It
// holds both conversion functions alongside the backing map and applies them at
// each read.
type mapKeyValuesView[FK comparable, FV any, K, V any] struct {
	m      map[FK]FV
	conv   func(FK, FV) (K, V)
	unconv func(K) (FK, bool)
}

func (v mapKeyValuesView[FK, FV, K, V]) Len() int {
	return len(v.m)
}

func (v mapKeyValuesView[FK, FV, K, V]) Get(key K) (V, bool) {
	var zero V
	fromKey, ok := v.unconv(key)
	if !ok {
		return zero, false
	}
	fromValue, ok := v.m[fromKey]
	if !ok {
		return zero, false
	}
	_, value := v.conv(fromKey, fromValue)
	return value, true
}

func (v mapKeyValuesView[FK, FV, K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for fromKey, fromValue := range v.m {
			key, _ := v.conv(fromKey, fromValue)
			if !yield(key) {
				return
			}
		}
	}
}

func (v mapKeyValuesView[FK, FV, K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for fromKey, fromValue := range v.m {
			_, value := v.conv(fromKey, fromValue)
			if !yield(value) {
				return
			}
		}
	}
}

func (v mapKeyValuesView[FK, FV, K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for fromKey, fromValue := range v.m {
			if !yield(v.conv(fromKey, fromValue)) {
				return
			}
		}
	}
}
