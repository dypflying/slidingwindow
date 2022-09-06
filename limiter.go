package ratelimit

//Limiter defines a default rate limiter
type Limiter interface {
	//try to acquire n
	//return n if absolutely acquired, less than n if part is acquired, 0 if none is acquired
	Acquire(n uint64) uint64

	//try to acquire one
	GetOne() bool
}
