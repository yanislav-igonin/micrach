package repositories

import (
	"reflect"
	"testing"
)

func TestToInt32s(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int32
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: []int{}, want: []int32{}},
		{name: "values", in: []int{1, 2, 2147483647}, want: []int32{1, 2, 2147483647}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toInt32s(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("toInt32s() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
