package base

import (
	"reflect"
	"testing"
)

func TestLabels_OneOf(t *testing.T) {
	type args struct {
		name   string
		values []string
	}
	tests := []struct {
		name string
		l    Labels
		args args
		want Labels
	}{{
		name: "dc label",
		l:    Labels{},
		args: args{
			name:   "datacenter",
			values: []string{"GRA1", "BHS3"},
		},
		want: Labels{"datacenter": "~(GRA1|BHS3)"},
	}, {
		name: "dc label",
		l:    nil,
		args: args{
			name:   "datacenter",
			values: []string{"GRA1", "BHS3"},
		},
		want: Labels{"datacenter": "~(GRA1|BHS3)"},
	}, {
		name: "nil values",
		l:    nil,
		args: args{
			name:   "datacenter",
			values: nil,
		},
		want: Labels{},
	}, {
		name: "no name",
		l:    Labels{},
		args: args{
			name:   "",
			values: []string{"A", "B"},
		},
		want: Labels{},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.OneOf(tt.args.name, tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Labels.OneOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
