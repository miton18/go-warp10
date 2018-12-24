package base

import (
	"testing"
	"time"
)

func TestDatapoints_Add(t *testing.T) {
	now := time.Now()

	type args struct {
		ts    time.Time
		value interface{}
	}
	tests := []struct {
		name string
		dps  Datapoints
		args args
	}{{
		name: "Add an int datapoint",
		dps:  Datapoints{},
		args: args{
			ts:    now,
			value: 1,
		},
	}, {
		name: "Add a string datapoint",
		dps:  Datapoints{},
		args: args{
			ts:    now,
			value: "test",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dps.Add(tt.args.ts, tt.args.value)

			if len(tt.dps) != 1 {
				t.Error("NewEmptyGTS() Datapoints must have 1 element")
				return
			}
			if len(tt.dps[0]) != 5 {
				t.Errorf("NewEmptyGTS() Datapoint must have 5 element: %+v", tt.dps[0])
				return
			}
			if tt.dps[0][0] != now.UnixNano()/1000 {
				t.Errorf("NewEmptyGTS() time = %+v, want %+v", tt.dps[0][0], tt.args.ts)
				return
			}
			if tt.dps[0][4] != tt.args.value {
				t.Errorf("NewEmptyGTS() value = %+v, want %+v", tt.dps[0][4], tt.args.value)
				return
			}
		})
	}
}
