// Dummy implementation of a thread safe rate limiting.
// It is used to limit the number of requests to showcase the rate limiting.
package dummy_rate_limiter

import (
	"sync"
	"time"
)

// LimitCounter holds the required variables to compose an in dummy rate limiter.
type LimitCounter struct {
	// Mutex is used for handling the concurrent
	// read/write requests for counter
	mu sync.RWMutex

	// count holds the number of requests
	count int
}

// NewLimitCounter returns a new LimitCounter instance with count 0.
func NewLimitCounter() *LimitCounter {
	return &LimitCounter{
		count: 0,
	}
}

// Limit returns true if the limit is reached.
// If the limit is reached, the caller should return error with code codes.ResourceExhausted.
// Limit is used by UnaryServerInterceptor and StreamServerInterceptor.
func (lc *LimitCounter) Limit() bool {
	lc.increase()
	return lc.isLimitReached()
}

// StartLimiter starts the limiter with a ticker.
// The limiter will be decreased every 1 second.
func (lc *LimitCounter) StartLimiter() {
	ticker := time.NewTicker(time.Second * 1)
	quit := make(chan struct{})
	go tickerToDecrease(ticker, quit, lc)
}

// get returns the count.
func (lc *LimitCounter) get() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.count
}

// increase increases the count.
// If limit is reached, count is not increased.
func (lc *LimitCounter) increase() {
	if lc.isLimitReached() {
		return
	}
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.count++
}

// decrease decreases the count.
// If count is 0, it does nothing.
func (lc *LimitCounter) decrease() {
	if lc.get() == 0 {
		return
	}
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.count--
}

// isLimitReached returns true if the limit is reached.
// Limit is 10 for showcase.
func (lc *LimitCounter) isLimitReached() bool {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.count >= 10
}

// tickerToDecrease decreases the count every given ticker.
func tickerToDecrease(ticker *time.Ticker, quit chan struct{}, lc *LimitCounter) {
	for {
		select {
		case <-ticker.C:
			lc.decrease()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
