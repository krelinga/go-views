package views_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-views"
)

func TestBagOfMapKeys_Len(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	var bag views.Bag[string] = views.BagOfMapKeys[string, int]{M: m}
	if got := bag.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}

	emptyBag := views.BagOfMapKeys[string, int]{M: map[string]int{}}
	if got := emptyBag.Len(); got != 0 {
		t.Errorf("Len() = %d, want 0", got)
	}

	var nilBag views.BagOfMapKeys[string, int]
	if got := nilBag.Len(); got != 0 {
		t.Errorf("Len() = %d, want 0 for nil map", got)
	}
}

func TestBagOfMapKeys_Values(t *testing.T) {
	m := map[int]string{1: "one", 2: "two", 3: "three"}
	bag := views.BagOfMapKeys[int, string]{M: m}
	got := []int{}
	for v := range bag.Values() {
		got = append(got, v)
	}
	want := []int{1, 2, 3}
	slices.Sort(got)
	slices.Sort(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Values() = %v, want %v", got, want)
	}

	var emptyBag views.BagOfMapKeys[int, string]
	got = []int{}
	for v := range emptyBag.Values() {
		got = append(got, v)
	}
	if len(got) != 0 {
		t.Errorf("Values() for empty map = %v, want []", got)
	}
}

func TestBagOfMapKeys_Has(t *testing.T) {
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
}
