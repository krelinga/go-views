package views_test

import (
	"cmp"
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-views"
)

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

func TestDictOfMap(t *testing.T) {
	tests := []struct {
		name        string
		dict        views.Dict[string, int]
		wantLen     int
		wantEntries []KV[string, int]
		wantKeys    []string
		wantValues  []int
		missingKeys []string
	}{
		{
			name:    "non-empty map",
			dict:    views.DictOfMap[string, int]{M: map[string]int{"a": 1, "b": 2, "c": 3}},
			wantLen: 3,
			wantEntries: []KV[string, int]{
				{Key: "a", Value: 1},
				{Key: "b", Value: 2},
				{Key: "c", Value: 3},
			},
			wantKeys:    []string{"a", "b", "c"},
			wantValues:  []int{1, 2, 3},
			missingKeys: []string{"d", "e"},
		},
		{
			name:        "empty map",
			dict:        views.DictOfMap[string, int]{M: map[string]int{}},
			wantLen:     0,
			wantEntries: nil,
			wantKeys:    nil,
			wantValues:  nil,
			missingKeys: []string{"a", "b", "c"},
		},
		{
			name:        "nil map",
			dict:        views.DictOfMap[string, int]{},
			wantLen:     0,
			wantEntries: nil,
			wantKeys:    nil,
			wantValues:  nil,
			missingKeys: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dict.Len(); got != tt.wantLen {
				t.Errorf("Len() = %d, want %d", got, tt.wantLen)
			}
			gotKeys := slices.Collect(tt.dict.Keys())
			slices.Sort(gotKeys)
			slices.Sort(tt.wantKeys)
			if !reflect.DeepEqual(gotKeys, tt.wantKeys) {
				t.Errorf("Keys() = %v, want %v", gotKeys, tt.wantKeys)
			}
			gotValues := slices.Collect(tt.dict.Values())
			slices.Sort(gotValues)
			slices.Sort(tt.wantValues)
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, tt.wantValues)
			}
			for _, wantEntry := range tt.wantEntries {
				if ok := tt.dict.Has(wantEntry.Key); !ok {
					t.Errorf("Has(%q) = false, want true", wantEntry.Key)
				}
				if _, ok := tt.dict.Get(wantEntry.Key); !ok {
					t.Errorf("Get(%q) = false, want true", wantEntry.Key)
				}
			}
			for _, missingK := range tt.missingKeys {
				if ok := tt.dict.Has(missingK); ok {
					t.Errorf("Has(%q) = true, want false", missingK)
				}
				if _, ok := tt.dict.Get(missingK); ok {
					t.Errorf("Get(%q) = true, want false", missingK)
				}
			}
			var gotEntries []KV[string, int]
			for gotK, gotV := range tt.dict.All() {
				gotEntries = append(gotEntries, KV[string, int]{Key: gotK, Value: gotV})
			}
			slices.SortFunc(gotEntries, func(a, b KV[string, int]) int {
				return cmp.Compare(a.Key, b.Key)
			})
			slices.SortFunc(tt.wantEntries, func(a, b KV[string, int]) int {
				return cmp.Compare(a.Key, b.Key)
			})
			if !reflect.DeepEqual(gotEntries, tt.wantEntries) {
				t.Errorf("Entries() = %v, want %v", gotEntries, tt.wantEntries)
			}
		})
	}
}
