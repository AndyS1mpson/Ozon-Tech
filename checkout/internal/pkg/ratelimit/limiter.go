// Rate limiter
package ratelimit

import (
	"context"
	"time"

	"golang.org/x/sync/semaphore"
)

// Define custom rate limiter
type RateLimiter struct {
	queue          chan struct{}
	maxConcurrency uint16
	rps uint16
	sem            *semaphore.Weighted
}

// Create a new rate limiter instance
// Takes two parameters as input:
// 		rps - limit on the number of requests per second
// 		maxConcurrency - limit on the number of simultaneous queries
func New(rps uint16, maxConcurrency uint16) *RateLimiter {
	sem := semaphore.NewWeighted(int64(maxConcurrency))

	r := &RateLimiter{
		queue: make(chan struct{}, rps),
		sem: sem,
		rps: rps,
		maxConcurrency: maxConcurrency,
	}

	go r.emitter()

	return r
}

// Keeps the number of requests per second and frees up resources
func (r *RateLimiter) emitter() {
	for range r.queue {
		time.Sleep(time.Second / time.Duration(r.rps))
		r.sem.Release(1)
	}
}

// Allows the request or waits for permission
func (r *RateLimiter) Acquire(ctx context.Context) error {
	return r.sem.Acquire(ctx, 1)
}

// Notifies you that the request has ended
func (r *RateLimiter) Release() {
	r.queue <- struct{}{}
}

// Frees up resources
func (r *RateLimiter) Close() {
	for i := 0; i < int(r.maxConcurrency); i++ {
		r.Acquire(context.Background())
	}
	close(r.queue)
}
