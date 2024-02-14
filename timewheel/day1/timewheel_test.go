package timewheel_test

import (
	timewheel "d_timewheel/timewheel/day1"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
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

func benchmarkCron(t *testing.B, n int) {
	e := cron.New()
	e.Start()
	defer e.Stop()
	for i := 0; i < n; i++ {
		e.AddFunc("@every 1s", func() {
			i = i
		})
	}

}

func BenchmarkCron1000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkCron(t, 1000)
	}
}

func BenchmarkCron10000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkCron(t, 10000)
	}
}
func BenchmarkCron20000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkCron(t, 20000)
	}
}

func BenchmarkCron50000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkCron(t, 50000)
	}
}

func BenchmarkCron100000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkCron(t, 100000)
	}
}

