package limit

import (
	"fmt"
	"sync"
)

// Quota enforces a cap on cumulative tokens/requests per period.
type Quota struct {
	mu            sync.Mutex
	tokensCap     uint64
	requestsCap   uint32
	tokensUsed    uint64
	requestsUsed  uint32
	overTokened   bool
	overRequested bool
}

// NewQuota creates a quota with the given caps (0 = unlimited).
func NewQuota(tokensCap uint64, requestsCap uint32) *Quota {
	return &Quota{tokensCap: tokensCap, requestsCap: requestsCap}
}

// CanConsume reports whether consuming tokens+1 request is permitted.
func (q *Quota) CanConsume(tokens uint64) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.tokensCap > 0 && q.tokensUsed+tokens > q.tokensCap {
		q.overTokened = true
		return false
	}
	if q.requestsCap > 0 && q.requestsUsed+1 > q.requestsCap {
		q.overRequested = true
		return false
	}
	return true
}

// Consume records token+request usage (call only after CanConsume is true).
func (q *Quota) Consume(tokens uint64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tokensUsed += tokens
	q.requestsUsed++
}

// Used reports current usage.
func (q *Quota) Used() (tokens uint64, requests uint32) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.tokensUsed, q.requestsUsed
}

// Exceeded returns an error describing which cap was hit first, if any.
func (q *Quota) Exceeded() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	switch {
	case q.overTokened:
		return fmt.Errorf("token quota exceeded (%d/%d)", q.tokensUsed, q.tokensCap)
	case q.overRequested:
		return fmt.Errorf("request quota exceeded (%d/%d)", q.requestsUsed, q.requestsCap)
	}
	return nil
}

// Reset clears usage within the period.
func (q *Quota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tokensUsed = 0
	q.requestsUsed = 0
	q.overTokened = false
	q.overRequested = false
}
