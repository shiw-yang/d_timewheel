package timewheel_test

import (
	timewheel "d_timewheel/timewheel/day1"
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	type args struct {
		interval time.Duration
		slotNums int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case1",
			args: args{
				interval: time.Minute,
				slotNums: 60,
			},
		},
		{
			name: "case2",
			args: args{
				interval: 0,
				slotNums: 0,
			},
		},
		{
			name: "case3",
			args: args{
				interval: time.Second,
				slotNums: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timewheel.NewTimeWheel(tt.args.interval, tt.args.slotNums)
		})
	}
}
