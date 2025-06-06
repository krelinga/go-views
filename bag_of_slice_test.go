package views_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-views"
)

func TestBagOfSlice(t *testing.T) {
	tests := []struct {
		name          string
		bag           views.Bag[int]
		wantLen       int
		wantValues    []int
		missingValues []int
	}{
		{
			name:          "non-empty slice",
			bag:           views.BagOfSlice[int]{S: []int{1, 2, 3}},
			wantLen:       3,
			wantValues:    []int{1, 2, 3},
			missingValues: []int{4, 5},
		},
		{
			name:          "empty slice",
			bag:           views.BagOfSlice[int]{S: []int{}},
			wantLen:       0,
			wantValues:    nil,
			missingValues: []int{1, 2, 3},
		},
		{
			name:          "nil slice",
			bag:           views.BagOfSlice[int]{},
			wantLen:       0,
			wantValues:    nil,
			missingValues: []int{1, 2, 3},
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