package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func run(rate uint64, window time.Duration, bucketNum uint32,
	routineNum, loopCount int, timeSpan time.Duration) (uint64, uint64, uint64) {

	r, _ := NewRateLimiter(rate, window, bucketNum)
	var success, fail uint64

	wg := sync.WaitGroup{}
	wg.Add(routineNum)
	start := time.Now()
	for i := 0; i < routineNum; i++ {
		go func() {
			for j := 0; j < loopCount; j++ {
				if r.GetOne() {
					atomic.AddUint64(&success, 1)
				} else {
					atomic.AddUint64(&fail, 1)
				}
				time.Sleep(timeSpan)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	elapse := time.Now().Sub(start).Milliseconds()

	expectedMin := rate / uint64(window/time.Second) * uint64(elapse/1000)
	expectedMax := expectedMin + rate
	return success, expectedMin, expectedMax
}
func Test_1(t *testing.T) {

	actual, expectedMin, expectedMax := run(100000, time.Second, 10, 100, 10000, 100*time.Microsecond)
	if actual < expectedMin || actual > expectedMax {
		t.Errorf("expected: (%d,%d), actual: %d\n", expectedMin, expectedMax, actual)
	}
}

func Test_2(t *testing.T) {

	actual, expectedMin, expectedMax := run(10000, time.Second, 100, 100, 10000, 100*time.Microsecond)
	if actual < expectedMin || actual > expectedMax {
		t.Errorf("expected: (%d,%d), actual: %d\n", expectedMin, expectedMax, actual)
	}
}

func Test_3(t *testing.T) {

	actual, expectedMin, expectedMax := run(1000, time.Second, 100, 100, 5000, time.Millisecond)
	if actual < expectedMin || actual > expectedMax {
		t.Errorf("expected: (%d,%d), actual: %d\n", expectedMin, expectedMax, actual)
	}
}

func Test_4(t *testing.T) {

	actual, expectedMin, expectedMax := run(1000, 10*time.Second, 100, 100, 5000, time.Millisecond)
	if actual < expectedMin || actual > expectedMax {
		t.Errorf("expected: (%d,%d), actual: %d\n", expectedMin, expectedMax, actual)
	}
}

func Test_5(t *testing.T) {

	actual, expectedMin, expectedMax := run(1000, 100*time.Second, 100, 100, 5000, time.Millisecond)
	if actual < expectedMin || actual > expectedMax {
		t.Errorf("expected: (%d,%d), actual: %d\n", expectedMin, expectedMax, actual)
	}
}

func Test_6(t *testing.T) {

	actual, expectedMin, expectedMax := run(100, 100*time.Second, 100, 100, 5000, time.Millisecond)
	if actual < expectedMin || actual > expectedMax {
		t.Errorf("expected: (%d,%d), actual: %d\n", expectedMin, expectedMax, actual)
	}
}

func Test_7(t *testing.T) {

	actual, expectedMin, expectedMax := run(10, 100*time.Second, 100, 100, 5000, time.Millisecond)
	if actual < expectedMin || actual > expectedMax {
		t.Errorf("expected: (%d,%d), actual: %d\n", expectedMin, expectedMax, actual)
	}
}
