package timewheel_test

import (
	timewheel "d_timewheel/timewheel/day1"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestNewTimeWheel(t *testing.T) {
	tw := timewheel.NewTimeWheel(time.Second, 60)
	defer tw.Stop()
	time.Sleep(time.Second * 3)
	t.Log("sleep 3 seconds")
	tw.AddTask(1, func() {
		t.Log("cron task execute")
	}, time.Now().Add(time.Second*10))
	time.Sleep(15 * time.Second)
}

func benchmarkCron(t *testing.B, n int) {
	e := cron.New()
	e.Start()
	defer e.Stop()
	for i := 0; i < n; i++ {
		e.AddFunc("@every 1s", func() {
			_ = i
		})
	}
}

func benchmarkTimeWheel(t *testing.B, n int) {
	tw := timewheel.NewTimeWheel(time.Second, 60)
	defer tw.Stop()
	for i := range n {
		tw.AddTask(int64(i), func() {}, time.Now().Add(time.Duration(i)*time.Second))
	}
}

func BenchmarkTimeWheel1000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkTimeWheel(t, 1000)
	}
}

func BenchmarkTimeWheel10000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkTimeWheel(t, 10000)
	}
}

func BenchmarkTimeWheel20000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkTimeWheel(t, 20000)
	}
}

func BenchmarkTimeWheel50000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkTimeWheel(t, 50000)
	}
}

func BenchmarkTimeWheel100000(t *testing.B) {
	for i := 0; i < t.N; i++ {
		benchmarkTimeWheel(t, 100000)
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
