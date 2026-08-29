// Package limit implements rate limiting and quotas (Master Architecture §12).
package limit

import (
	"sync"
	"time"
)

// Limits configures rate and quota constraints.
type Limits struct {
	MaxRequestsPerMinute uint32
	MaxTokensPerMinute   uint32
	TokensMonthly        uint64
	RequestsDaily        uint32
	// TokensUsedMonthly is the running count of tokens consumed this period.
	TokensUsedMonthly uint64
	// RequestsUsedDaily is the running count of requests this period.
	RequestsUsedDaily uint32
}

// RateLimiter enforces max events per time window using a sliding-window
// counter. It is a single-window limiter (per-minute) suitable for per-provider
// throttling.
type RateLimiter struct {
	mu     sync.Mutex
	window time.Duration
	limit  uint32
	start  time.Time
	count  uint32
	now    func() time.Time
}

// NewRateLimiter creates a limiter allowing limit events per window.
func NewRateLimiter(limit uint32, window time.Duration) *RateLimiter {
	return &RateLimiter{window: window, limit: limit, start: time.Now(), now: time.Now}
}

// Allow reports whether one more event fits in the current window.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	if now.Sub(r.start) >= r.window {
		r.start = now
		r.count = 0
	}
	if r.limit == 0 || r.count < r.limit {
		r.count++
		return true
	}
	return false
}

// Reset resets the window.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.start = r.now()
	r.count = 0
}

// errLimitExceeded is returned when the per-minute rate is exhausted.
type errLimitExceeded struct{}

func (errLimitExceeded) Error() string { return "rate limit exceeded" }
