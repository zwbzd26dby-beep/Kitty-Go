package limit

import (
	"time"
)

// LimitManager guards per-provider execution against rate and quota caps
// (Master Architecture §12).
type LimitManager interface {
	// CheckAndTrack verifies the caps allow a request consuming tokens, then
	// if permitted records it. It returns an error when rolled back.
	CheckAndTrack(tokens uint64) error
	// ResetWindows resets per-minute and periodic counters.
	ResetWindows()
}

// Manager is a LimitManager backed by a RateLimiter and a Quota.
type Manager struct {
	rate  *RateLimiter
	quota *Quota
}

// NewManager builds a Manager from the given limits.
func NewManager(l Limits) *Manager {
	return &Manager{
		rate:  NewRateLimiter(l.MaxRequestsPerMinute, time.Minute),
		quota: NewQuota(l.TokensMonthly, l.RequestsDaily),
	}
}

// CheckAndTrack returns an error if the rate or quota avoids the request,
// otherwise records its tokens.
func (m *Manager) CheckAndTrack(tokens uint64) error {
	if !m.rate.Allow() {
		return errLimitExceeded{}
	}
	if !m.quota.CanConsume(tokens) {
		return m.quota.Exceeded()
	}
	m.quota.Consume(tokens)
	return nil
}

// ResetWindows resets both the rate window and the quota period.
func (m *Manager) ResetWindows() {
	m.rate.Reset()
	m.quota.Reset()
}

// Used exposes current quota usage.
func (m *Manager) Used() (tokens uint64, requests uint32) {
	return m.quota.Used()
}

var _ LimitManager = (*Manager)(nil)
