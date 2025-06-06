package views_test

import (
	"cmp"
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-views"
)

func TestDictOfSlice(t *testing.T) {
	tests := []struct {
		name        string
		dict        views.Dict[int, string]
		wantLen     int
		wantEntries []KV[int, string]
		wantKeys    []int
		wantValues  []string
		missingKeys []int
	}{
		{
			name:    "non-empty slice",
			dict:    views.DictOfSlice[string]{S: []string{"a", "b", "c"}},
			wantLen: 3,
			wantEntries: []KV[int, string]{
				{Key: 0, Value: "a"},
				{Key: 1, Value: "b"},
				{Key: 2, Value: "c"},
			},
			wantKeys:    []int{0, 1, 2},
			wantValues:  []string{"a", "b", "c"},
			missingKeys: []int{-1, 5},
		},
		{
			name:        "empty slice",
			dict:        views.DictOfSlice[string]{S: []string{}},
			wantLen:     0,
			wantEntries: nil,
			wantKeys:    nil,
			wantValues:  nil,
			missingKeys: []int{0, 1, 2},
		},
		{
			name:        "nil slice",
			dict:        views.DictOfSlice[string]{},
			wantLen:     0,
			wantEntries: nil,
			wantKeys:    nil,
			wantValues:  nil,
			missingKeys: []int{0, 1, 2},
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
			for _, wantK := range tt.wantKeys {
				if !tt.dict.Has(wantK) {
					t.Errorf("HasKey(%q) = false, want true", wantK)
				}
			}
			for _, missingK := range tt.missingKeys {
				if tt.dict.Has(missingK) {
					t.Errorf("HasKey(%q) = true, want false", missingK)
				}
			}
			for _, wantEntry := range tt.wantEntries {
				if gotV, ok := tt.dict.Get(wantEntry.Key); !ok || gotV != wantEntry.Value {
					t.Errorf("Get(%q) = %v, want %v", wantEntry.Key, gotV, wantEntry.Value)
				}
			}
			for _, missingK := range tt.missingKeys {
				if _, ok := tt.dict.Get(missingK); ok {
					t.Errorf("Get(%q) = true, want false", missingK)
				}
			}
			var gotEntries []KV[int, string]
			for gotK, gotV := range tt.dict.All() {
				gotEntries = append(gotEntries, KV[int, string]{Key: gotK, Value: gotV})
			}
			slices.SortFunc(gotEntries, func(a, b KV[int, string]) int {
				return cmp.Compare(a.Key, b.Key)
			})
			slices.SortFunc(tt.wantEntries, func(a, b KV[int, string]) int {
				return cmp.Compare(a.Key, b.Key)
			})
			if !reflect.DeepEqual(gotEntries, tt.wantEntries) {
				t.Errorf("Entries() = %v, want %v", gotEntries, tt.wantEntries)
			}
		})
	}
}
