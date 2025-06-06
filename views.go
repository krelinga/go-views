// Package views provides read-only views of collections.
//
// This is useful in situations where you want to expose a collection as part of an API while
// maintaining encapsulation and preventing modification of the underlying data structure.
//
// This package also provides wrappers around Go's built-in slices and maps to satsfiy these
// read-only interfaces.
//
// Callers are, of course, free to implement these interfaces themselves if they need more control.
package views

import (
	"iter"
)

// List is a read-only view of a collection of values.
type List[T any] interface {
	// Len returns the number of elements in the list.
	Len() int

	// Values returns a sequence of the values in the list.
	// Callers should not assume that the order of the values is consistent.
	Values() iter.Seq[T]
}

// Bag is a read-only view of a collection of values that allows checking for the presence of values.
type Bag[T comparable] interface {
	// All Bags are also Lists.
	List[T]

	// Has checks if the bag contains the specified value.
	Has(T) bool
}

// Dict is a read-only view of a collection of key-value pairs.
type Dict[K comparable, V any] interface {
	// All Dicts are also Lists, where the values are the values of the map.
	List[V]

	// Keys returns a sequence of the keys in the Dict.
	Keys() iter.Seq[K]

	// Get retrieves the value associated with the specified key.
	// If the key does not exist, it returns the zero value of V and false.
	Get(K) (V, bool)

	// Has checks if the Dict contains the specified key.
	Has(K) bool

	// All returns the sequence of key-value pairs in the Dict.
	All() iter.Seq2[K, V]
}
