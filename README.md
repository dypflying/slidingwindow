Go Sliding Window Rate Limiter
==============================
This package provides a Golang lock-free implementation of the sliding-window algorithm, which uses the compare-and-swap (CAS) mechanism to for data consistency, and the implementation uses a configurable bucket number for splicing the time windows. 

Note: This implementation is still experimental. 

Quick Start
=====
### simple rate limiter with the GetOne() method.
```go
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
```


### Rate limiter example with the Acquire() method
```go
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
```

License 
=======
MIT License

Copyright (c) Yunpeng Deng(dypflying)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


Report Issues 
=============
Please report bugs via the [GitHub Issue Tracker](https://github.com/dypflying/slidingwindow/issues) or [Contact Author](#author-and-contact) 

Contact Author
==============
- Author: Yunpeng Deng (dypflying)
- Mailto: dypflying@sina.com

Either English or Chinese is welcome.