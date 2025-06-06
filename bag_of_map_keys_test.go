package views_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-views"
)

func TestBagOfMapKeys(t *testing.T) {
	tests := []struct {
		name          string
		bag           views.Bag[string]
		wantLen       int
		wantValues    []string
		missingValues []string
	}{
		{
			name:          "non-empty map",
			bag:           views.BagOfMapKeys[string, int]{M: map[string]int{"a": 1, "b": 2, "c": 3}},
			wantLen:       3,
			wantValues:    []string{"a", "b", "c"},
			missingValues: []string{"d", "e"},
		},
		{
			name:          "empty map",
			bag:           views.BagOfMapKeys[string, int]{M: map[string]int{}},
			wantLen:       0,
			wantValues:    nil,
			missingValues: []string{"a", "b", "c"},
		},
		{
			name:          "nil map",
			bag:           views.BagOfMapKeys[string, int]{},
			wantLen:       0,
			wantValues:    nil,
			missingValues: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bag.Len(); got != tt.wantLen {
				t.Errorf("Len() = %d, want %d", got, tt.wantLen)
			}
			gotValues := slices.Collect(tt.bag.Values())
			slices.Sort(gotValues)
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, tt.wantValues)
			}
			for _, wantV := range tt.wantValues {
				if !tt.bag.Has(wantV) {
					t.Errorf("Has(%q) = false, want true", wantV)
				}
			}
			for _, missingV := range tt.missingValues {
				if tt.bag.Has(missingV) {
					t.Errorf("Has(%q) = true, want false", missingV)
				}
			}
		})
	}
}
