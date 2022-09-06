/*
Copyright (c) Yunpeng Deng(dypflying)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package ratelimit

import (
	"errors"
	"sync/atomic"
	"time"
)

var zeroTime time.Time

type window struct {
	//the total arrrive in current bucket window
	arrives uint64
	//_rearPadding [7]uint64 //just for padding the 64-bytes based cache line
	//the create time of the bucket window
	createTime time.Time
}

type rateLimiter struct {
	//current total arrives in the ring buffer
	_forePadding [8]uint64 //just for padding the 64-bytes based cache line
	totalArrives uint64    //atomic
	_rearPadding [7]uint64 //just for padding the 64-bytes based cache line

	rate uint64
	//the ring buffer struct
	ringBuffer []window
	//the size of the ring
	ringSize uint32
	//the current pointer of the ring
	ringHead uint32
	//the time span of a bucket window
	bucketWindowDuration time.Duration
}

//NewRateLimiter is the contructor for a rate limiter
func NewRateLimiter(rate uint64, timeWindow time.Duration, bucketNum uint32) (Limiter, error) {

	//zero bucket is not allowed
	if bucketNum == 0 {
		return nil, errors.New("bucketNum can not be zero")
	}

	return &rateLimiter{
		rate:                 rate,
		ringSize:             bucketNum,
		ringBuffer:           make([]window, bucketNum, bucketNum),
		ringHead:             0,
		bucketWindowDuration: timeWindow / time.Duration(bucketNum),
	}, nil
}

//try to acquire n
//return n if absolutely acquired, less than n if part is acquired, 0 if none is acquired
func (r *rateLimiter) Acquire(n uint64) uint64 {

	now := time.Now()
	if r.ringBuffer[r.ringHead].createTime == zeroTime {
		r.ringBuffer[r.ringHead].createTime = now
	}

	for {

		oldTotal := atomic.LoadUint64(&r.totalArrives)
		currentTotal := oldTotal

		//it falls in one of the previous buckets
		if now.Sub(r.ringBuffer[r.ringHead].createTime) > r.bucketWindowDuration {

			//the pointer to the tail of the ring
			lastBucket := (r.ringHead + 1) % r.ringSize

			//another solution is just to deduce the currentArrives from the head bucket
			//however the bucket.arrives may not be consistent with the total
			//currentTotal = currentTotal - r.ringBuffer[last].arrives

			//re-calculate the ringTotal
			currentTotal := uint64(0)
			for i := 0; i < int(r.ringSize); i++ {
				if i != int(lastBucket) {
					currentTotal += r.ringBuffer[i].arrives
				}
			}
			if atomic.CompareAndSwapUint64(&r.totalArrives, oldTotal, currentTotal) {
				atomic.StoreUint64(&r.ringBuffer[lastBucket].arrives, 0) //empty the last bucket
				atomic.StoreUint32(&r.ringHead, lastBucket)              //move the pointer to the next bucket
				r.ringBuffer[r.ringHead].createTime = now
				continue
			}

		} else { //here, it falls in the current bucket

			acquirableCount := int64(r.rate) - int64(currentTotal)
			if acquirableCount < 0 {
				return 0
			} else if acquirableCount > int64(n) {
				acquirableCount = int64(n)
			}
			currentTotal += uint64(acquirableCount)
			if atomic.CompareAndSwapUint64(&r.totalArrives, oldTotal, currentTotal) {
				//there could introduce slight data inconsistency bettern totalArrives and bucket's arrives
				atomic.AddUint64(&r.ringBuffer[r.ringHead].arrives, uint64(acquirableCount))
				return uint64(acquirableCount)
			}
		}

	}
}

//try to acquire one
func (r *rateLimiter) GetOne() bool {
	return (r.Acquire(1) == 1)
}
