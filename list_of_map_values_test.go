package views_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/krelinga/go-views"
)

func TestListMapValues(t *testing.T) {
	tests := []struct {
		name        string
		list        views.List[string]
		wantLen     int
		wantValues  []string
	}{
		{
			name:    "non-empty list",
			list:    views.ListOfMapValues[int, string]{M: map[int]string{0: "a", 1: "b", 2: "c"}},
			wantLen: 3,
			wantValues:  []string{"a", "b", "c"},
		},
		{
			name:        "empty list",
			list:        views.ListOfMapValues[int, string]{M: map[int]string{}},
			wantLen:     0,
			wantValues:  nil,
		},
		{
			name:        "nil list",
			list:        views.ListOfMapValues[int, string]{},
			wantLen:     0,
			wantValues:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.Len(); got != tt.wantLen {
				t.Errorf("Len() = %d, want %d", got, tt.wantLen)
			}
			gotValues := slices.Collect(tt.list.Values())
			slices.Sort(gotValues)
			slices.Sort(tt.wantValues)
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("Values() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}