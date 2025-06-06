package views_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-views"
)

func TestBagOfMapKeys(t *testing.T) {
	t.Run("Len", func(t *testing.T) {
		tests := []struct {
			name string
			bag  views.Bag[string]
			want int
		}{
			{
				name: "non-empty map",
				bag:  views.BagOfMapKeys[string, int]{M: map[string]int{"a": 1, "b": 2, "c": 3}},
				want: 3,
			},
			{
				name: "empty map",
				bag:  views.BagOfMapKeys[string, int]{M: map[string]int{}},
				want: 0,
			},
			{
				name: "nil map",
				bag:  views.BagOfMapKeys[string, int]{},
				want: 0,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := tt.bag.Len(); got != tt.want {
					t.Errorf("Len() = %d, want %d", got, tt.want)
				}
			})
		}
	})

	t.Run("Values", func(t *testing.T) {
		m := map[int]string{1: "one", 2: "two", 3: "three"}
		bag := views.BagOfMapKeys[int, string]{M: m}
		got := slices.Collect(bag.Values())
		want := []int{1, 2, 3}
		slices.Sort(got)
		slices.Sort(want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Values() = %v, want %v", got, want)
		}

		var emptyBag views.BagOfMapKeys[int, string]
		got = slices.Collect(emptyBag.Values())
		if len(got) != 0 {
			t.Errorf("Values() for empty map = %v, want []", got)
		}
	})

	t.Run("Has", func(t *testing.T) {
		m := map[string]int{"foo": 1, "bar": 2}
		bag := views.BagOfMapKeys[string, int]{M: m}

		if !bag.Has("foo") {
			t.Errorf("Has(\"foo\") = false, want true")
		}
		if !bag.Has("bar") {
			t.Errorf("Has(\"bar\") = false, want true")
		}
		if bag.Has("baz") {
			t.Errorf("Has(\"baz\") = true, want false")
		}

		var nilBag views.BagOfMapKeys[string, int]
		if nilBag.Has("foo") {
			t.Errorf("Has(\"foo\") on nil map = true, want false")
		}
	})
}
