package views_test

import (
	"reflect"
	"slices"
	"strconv"
	"testing"

	views "github.com/krelinga/go-views/v2"
)

// namedStrings and namedInts exercise the ~[]T constraint: a view must accept
// any type whose underlying type is a slice, not just []T itself.
type namedStrings []string

type namedInts []int

// TestSlice covers the behavior every Slice implementation shares. Declaring the
// table field as views.Slice[string] makes each entry a compile-time
// conformance check, and lets both constructors share one table by having them
// produce the same element type.
//
// Unlike the v1 views, iteration order is part of the contract here: Slice is
// index-addressed, so results are compared in order rather than sorted.
func TestSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice views.Slice[string]
		want  []string
	}{
		{
			name:  "NewSlice non-empty",
			slice: views.NewSlice([]string{"a", "b", "c"}),
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "NewSlice empty",
			slice: views.NewSlice([]string{}),
			want:  nil,
		},
		{
			name:  "NewSlice nil backing",
			slice: views.NewSlice([]string(nil)),
			want:  nil,
		},
		{
			name:  "NewSlice named slice type",
			slice: views.NewSlice(namedStrings{"a", "b"}),
			want:  []string{"a", "b"},
		},
		{
			name:  "NewSliceVals non-empty",
			slice: views.NewSliceVals([]int{1, 2, 3}, strconv.Itoa),
			want:  []string{"1", "2", "3"},
		},
		{
			name:  "NewSliceVals empty",
			slice: views.NewSliceVals([]int{}, strconv.Itoa),
			want:  nil,
		},
		{
			name:  "NewSliceVals nil backing",
			slice: views.NewSliceVals([]int(nil), strconv.Itoa),
			want:  nil,
		},
		{
			name:  "NewSliceVals named slice type",
			slice: views.NewSliceVals(namedInts{1, 2}, strconv.Itoa),
			want:  []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.slice.Len(); got != len(tt.want) {
				t.Errorf("Len() = %d, want %d", got, len(tt.want))
			}

			if got := slices.Collect(tt.slice.Values()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}

			for i, want := range tt.want {
				if got := tt.slice.At(i); got != want {
					t.Errorf("At(%d) = %q, want %q", i, got, want)
				}
			}

			var gotIndexes []int
			var gotValues []string
			for i, v := range tt.slice.All() {
				gotIndexes = append(gotIndexes, i)
				gotValues = append(gotValues, v)
			}
			var wantIndexes []int
			for i := range tt.want {
				wantIndexes = append(wantIndexes, i)
			}
			if !reflect.DeepEqual(gotIndexes, wantIndexes) {
				t.Errorf("All() yielded indexes %v, want %v", gotIndexes, wantIndexes)
			}
			if !reflect.DeepEqual(gotValues, tt.want) {
				t.Errorf("All() yielded values %v, want %v", gotValues, tt.want)
			}
		})
	}
}

// TestSliceAtOutOfRange pins At to the panic semantics of the builtin index
// operation, which is why At returns a bare T rather than a value and an ok.
func TestSliceAtOutOfRange(t *testing.T) {
	tests := []struct {
		name  string
		slice views.Slice[string]
		index int
	}{
		{"NewSlice negative", views.NewSlice([]string{"a"}), -1},
		{"NewSlice past end", views.NewSlice([]string{"a"}), 1},
		{"NewSlice nil backing", views.NewSlice([]string(nil)), 0},
		{"NewSliceVals negative", views.NewSliceVals([]int{1}, strconv.Itoa), -1},
		{"NewSliceVals past end", views.NewSliceVals([]int{1}, strconv.Itoa), 1},
		{"NewSliceVals nil backing", views.NewSliceVals([]int(nil), strconv.Itoa), 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("At(%d) did not panic", tt.index)
				}
			}()
			tt.slice.At(tt.index)
		})
	}
}

// TestNewSliceValsConvertsLazily pins the documented cost model: nothing is
// converted up front, each read converts exactly once, and abandoning an
// iteration converts nothing further.
func TestNewSliceValsConvertsLazily(t *testing.T) {
	var calls int
	count := func(i int) string {
		calls++
		return strconv.Itoa(i)
	}
	slice := views.NewSliceVals([]int{1, 2, 3}, count)

	if calls != 0 {
		t.Errorf("construction made %d conversions, want 0", calls)
	}

	if slice.Len(); calls != 0 {
		t.Errorf("Len() made %d conversions, want 0", calls)
	}

	calls = 0
	slice.At(0)
	slice.At(0)
	if calls != 2 {
		t.Errorf("two At calls made %d conversions, want 2", calls)
	}

	calls = 0
	for range slice.Values() {
	}
	if calls != 3 {
		t.Errorf("full Values() range made %d conversions, want 3", calls)
	}

	tests := []struct {
		name string
		stop func(views.Slice[string])
	}{
		{
			name: "Values",
			stop: func(s views.Slice[string]) {
				for range s.Values() {
					break
				}
			},
		},
		{
			name: "All",
			stop: func(s views.Slice[string]) {
				for range s.All() {
					break
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run("break during "+tt.name, func(t *testing.T) {
			calls = 0
			tt.stop(slice)
			if calls != 1 {
				t.Errorf("breaking after one element made %d conversions, want 1", calls)
			}
		})
	}
}

// TestSliceAliasesBacking pins the documented aliasing contract: a view
// prevents modification through itself, but does not freeze the backing slice.
func TestSliceAliasesBacking(t *testing.T) {
	backing := []string{"a", "b"}
	slice := views.NewSlice(backing)

	backing[0] = "changed"

	if got := slice.At(0); got != "changed" {
		t.Errorf("At(0) = %q after backing write, want %q", got, "changed")
	}
	if got := slices.Collect(slice.Values()); !reflect.DeepEqual(got, []string{"changed", "b"}) {
		t.Errorf("Values() = %v after backing write", got)
	}
}
