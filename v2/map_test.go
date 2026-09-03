package views_test

import (
	"maps"
	"reflect"
	"slices"
	"strconv"
	"testing"

	views "github.com/krelinga/go-views/v2"
)

// namedLabels and namedAges exercise the ~map[K]V constraint.
type namedLabels map[string]string

type namedAges map[string]int

// kvConv derives the exposed value from both the key and the value, which is
// the capability a pair-wise converter has and two independent converters do
// not.
func kvConv(key int, value string) (string, string) {
	return strconv.Itoa(key), strconv.Itoa(key) + ":" + value
}

// kvUnconv recovers a backing key, reporting false for exposed keys that have
// no backing counterpart.
func kvUnconv(key string) (int, bool) {
	i, err := strconv.Atoi(key)
	return i, err == nil
}

// TestMap covers the behavior every Map implementation shares. Declaring the
// table field as views.Map[string, string] makes each entry a compile-time
// conformance check, and lets all three constructors share one table by having
// them produce the same key and value types.
//
// Map iteration order is unspecified, so keys and values are sorted before
// comparison.
func TestMap(t *testing.T) {
	tests := []struct {
		name    string
		m       views.Map[string, string]
		want    map[string]string
		missing []string
	}{
		{
			name:    "NewMap non-empty",
			m:       views.NewMap(map[string]string{"a": "1", "b": "2"}),
			want:    map[string]string{"a": "1", "b": "2"},
			missing: []string{"z", ""},
		},
		{
			name:    "NewMap empty",
			m:       views.NewMap(map[string]string{}),
			want:    map[string]string{},
			missing: []string{"a"},
		},
		{
			name:    "NewMap nil backing",
			m:       views.NewMap(map[string]string(nil)),
			want:    map[string]string{},
			missing: []string{"a"},
		},
		{
			name:    "NewMap named map type",
			m:       views.NewMap(namedLabels{"a": "1"}),
			want:    map[string]string{"a": "1"},
			missing: []string{"z"},
		},
		{
			name:    "NewMapVals non-empty",
			m:       views.NewMapVals(map[string]int{"a": 1, "b": 2}, strconv.Itoa),
			want:    map[string]string{"a": "1", "b": "2"},
			missing: []string{"z"},
		},
		{
			name:    "NewMapVals empty",
			m:       views.NewMapVals(map[string]int{}, strconv.Itoa),
			want:    map[string]string{},
			missing: []string{"a"},
		},
		{
			name:    "NewMapVals nil backing",
			m:       views.NewMapVals(map[string]int(nil), strconv.Itoa),
			want:    map[string]string{},
			missing: []string{"a"},
		},
		{
			name:    "NewMapVals named map type",
			m:       views.NewMapVals(namedAges{"a": 1}, strconv.Itoa),
			want:    map[string]string{"a": "1"},
			missing: []string{"z"},
		},
		{
			name:    "NewMapKeyValues non-empty",
			m:       views.NewMapKeyValues(map[int]string{1: "one", 2: "two"}, kvConv, kvUnconv),
			want:    map[string]string{"1": "1:one", "2": "2:two"},
			missing: []string{"99", "not-a-number"},
		},
		{
			name:    "NewMapKeyValues empty",
			m:       views.NewMapKeyValues(map[int]string{}, kvConv, kvUnconv),
			want:    map[string]string{},
			missing: []string{"1", "not-a-number"},
		},
		{
			name:    "NewMapKeyValues nil backing",
			m:       views.NewMapKeyValues(map[int]string(nil), kvConv, kvUnconv),
			want:    map[string]string{},
			missing: []string{"1", "not-a-number"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Len(); got != len(tt.want) {
				t.Errorf("Len() = %d, want %d", got, len(tt.want))
			}

			for key, want := range tt.want {
				got, ok := tt.m.Get(key)
				if !ok {
					t.Errorf("Get(%q) reported false, want true", key)
					continue
				}
				if got != want {
					t.Errorf("Get(%q) = %q, want %q", key, got, want)
				}
			}

			for _, key := range tt.missing {
				got, ok := tt.m.Get(key)
				if ok {
					t.Errorf("Get(%q) = %q, true; want false", key, got)
				}
				if got != "" {
					t.Errorf("Get(%q) returned %q on miss, want the zero value", key, got)
				}
			}

			gotKeys := slices.Sorted(tt.m.Keys())
			wantKeys := slices.Sorted(maps.Keys(tt.want))
			if !reflect.DeepEqual(gotKeys, wantKeys) {
				t.Errorf("Keys() = %v, want %v", gotKeys, wantKeys)
			}

			gotValues := slices.Sorted(tt.m.Values())
			wantValues := slices.Sorted(maps.Values(tt.want))
			if !reflect.DeepEqual(gotValues, wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, wantValues)
			}

			if got := maps.Collect(tt.m.All()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewMapKeyValuesGet separates the two ways a lookup can fail, since they
// take different paths: unconv can reject the key outright, or recover a key
// that the backing map does not hold.
//
// The backing map deliberately holds the zero key. That is what makes the
// unrecoverable cases meaningful: an implementation that ignored unconv's bool
// would fall back to the zero key and return this entry, and a backing map
// without one would let that bug pass unnoticed.
func TestNewMapKeyValuesGet(t *testing.T) {
	m := views.NewMapKeyValues(map[int]string{0: "zero", 1: "one"}, kvConv, kvUnconv)

	tests := []struct {
		name   string
		key    string
		want   string
		wantOK bool
	}{
		{name: "recoverable key present in backing map", key: "1", want: "1:one", wantOK: true},
		{name: "recoverable zero key present in backing map", key: "0", want: "0:zero", wantOK: true},
		{name: "recoverable key absent from backing map", key: "99", wantOK: false},
		{name: "unrecoverable key", key: "not-a-number", wantOK: false},
		{name: "unrecoverable empty key", key: "", wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := m.Get(tt.key)
			if ok != tt.wantOK {
				t.Errorf("Get(%q) ok = %v, want %v", tt.key, ok, tt.wantOK)
			}
			if got != tt.want {
				t.Errorf("Get(%q) = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

// TestNewMapKeyValuesNonInjectiveKeys pins the documented behavior when a key
// conversion collapses distinct backing keys: Len keeps counting backing
// entries, and iteration keeps yielding the duplicates.
func TestNewMapKeyValuesNonInjectiveKeys(t *testing.T) {
	collapse := func(key int, value string) (string, string) { return "same", value }
	recover := func(string) (int, bool) { return 1, true }
	m := views.NewMapKeyValues(map[int]string{1: "a", 2: "b"}, collapse, recover)

	if got := m.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2 (backing entries)", got)
	}
	if got := slices.Sorted(m.Keys()); !reflect.DeepEqual(got, []string{"same", "same"}) {
		t.Errorf("Keys() = %v, want the duplicate to be yielded twice", got)
	}
	if got := slices.Sorted(m.Values()); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("Values() = %v, want both backing values", got)
	}
}

// TestMapConvertsLazily pins the documented cost model for both converting
// views: nothing is converted up front, each read converts exactly once, and
// abandoning an iteration converts nothing further.
func TestMapConvertsLazily(t *testing.T) {
	tests := []struct {
		name string
		// newMap builds a three-entry view that calls count once per conversion.
		newMap func(count func()) views.Map[string, string]
		// presentKey and missingKey are exposed keys, which differ by
		// implementation: NewMapVals passes backing keys through, while
		// NewMapKeyValues converts them.
		presentKey string
		missingKey string
		// keysConverts records whether Keys() has to convert, which differs by
		// implementation: NewMapVals passes keys through untouched, while
		// NewMapKeyValues must run the pair-wise converter to produce one.
		keysConverts bool
	}{
		{
			name: "NewMapVals",
			newMap: func(count func()) views.Map[string, string] {
				return views.NewMapVals(map[string]int{"a": 1, "b": 2, "c": 3}, func(i int) string {
					count()
					return strconv.Itoa(i)
				})
			},
			presentKey:   "a",
			missingKey:   "zz",
			keysConverts: false,
		},
		{
			name: "NewMapKeyValues",
			newMap: func(count func()) views.Map[string, string] {
				return views.NewMapKeyValues(map[int]string{1: "a", 2: "b", 3: "c"},
					func(k int, v string) (string, string) {
						count()
						return strconv.Itoa(k), v
					},
					kvUnconv,
				)
			},
			presentKey:   "1",
			missingKey:   "99",
			keysConverts: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls int
			m := tt.newMap(func() { calls++ })

			if calls != 0 {
				t.Errorf("construction made %d conversions, want 0", calls)
			}

			if m.Len(); calls != 0 {
				t.Errorf("Len() made %d conversions, want 0", calls)
			}

			calls = 0
			m.Get(tt.presentKey)
			m.Get(tt.presentKey)
			if calls != 2 {
				t.Errorf("two Get calls made %d conversions, want 2", calls)
			}

			calls = 0
			m.Get(tt.missingKey)
			if calls != 0 {
				t.Errorf("Get on a missing key made %d conversions, want 0", calls)
			}

			calls = 0
			for range m.Values() {
			}
			if calls != 3 {
				t.Errorf("full Values() range made %d conversions, want 3", calls)
			}

			calls = 0
			for range m.Keys() {
			}
			wantKeyCalls := 0
			if tt.keysConverts {
				wantKeyCalls = 3
			}
			if calls != wantKeyCalls {
				t.Errorf("full Keys() range made %d conversions, want %d", calls, wantKeyCalls)
			}

			for _, seq := range []struct {
				name string
				stop func()
			}{
				{"Values", func() {
					for range m.Values() {
						break
					}
				}},
				{"All", func() {
					for range m.All() {
						break
					}
				}},
			} {
				calls = 0
				seq.stop()
				if calls != 1 {
					t.Errorf("breaking after one element of %s() made %d conversions, want 1", seq.name, calls)
				}
			}
		})
	}
}

// TestMapAliasesBacking pins the documented aliasing contract: a view prevents
// modification through itself, but does not freeze the backing map.
func TestMapAliasesBacking(t *testing.T) {
	backing := map[string]string{"a": "1"}
	m := views.NewMap(backing)

	backing["a"] = "changed"
	backing["b"] = "2"

	if got, ok := m.Get("a"); got != "changed" || !ok {
		t.Errorf("Get(a) = %q, %v after backing write, want %q, true", got, ok, "changed")
	}
	if got, ok := m.Get("b"); got != "2" || !ok {
		t.Errorf("Get(b) = %q, %v after backing insert, want %q, true", got, ok, "2")
	}
	if got := m.Len(); got != 2 {
		t.Errorf("Len() = %d after backing insert, want 2", got)
	}
}
