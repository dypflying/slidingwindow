package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	ratelimit "github.com/dypflying/slidingwindow"
)

func testAcquire() {
	rate := uint64(100000)
	routineNum := 100
	loopCount := 10000
	timeSpan := time.Millisecond

	r, _ := ratelimit.NewRateLimiter(rate, time.Second, 10)
	var success, fail uint64

	wg := sync.WaitGroup{}
	wg.Add(routineNum)
	start := time.Now()
	for i := 0; i < routineNum; i++ {
		go func() {
			for j := 0; j < loopCount; j++ {
				n := r.Acquire(15)
				atomic.AddUint64(&success, n)
				atomic.AddUint64(&fail, 15-n)
				time.Sleep(timeSpan)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	elapse := time.Now().Sub(start).Milliseconds()
	fmt.Printf("success = %d, fail = %d in %d milliseconds\n", success, fail, elapse)
}

func testGetOne() {
	rate := uint64(100000)
	routineNum := 100
	loopCount := 10000
	timeSpan := 100 * time.Microsecond

	r, _ := ratelimit.NewRateLimiter(rate, time.Second, 10)
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
	fmt.Printf("success = %d, fail = %d in %d milliseconds\n", success, fail, elapse)
}

func main() {

	fmt.Println("The number of CPU Cores:", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	//testGetOne()
	testAcquire()
}
