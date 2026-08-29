package limittest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/limit"
)

func TestRateLimiterAllowsUpToLimit(t *testing.T) {
	rl := limit.NewRateLimiter(2, time.Minute)
	if !rl.Allow() {
		t.Fatal("expected first allow")
	}
	if !rl.Allow() {
		t.Fatal("expected second allow")
	}
	if rl.Allow() {
		t.Fatal("expected third to be denied")
	}
}

func TestRateLimiterResets(t *testing.T) {
	rl := limit.NewRateLimiter(1, time.Minute)
	rl.Allow()
	if rl.Allow() {
		t.Fatal("expected denial before reset")
	}
	rl.Reset()
	if !rl.Allow() {
		t.Fatal("expected allow after reset")
	}
}

func TestQuotaBlocksOnTokens(t *testing.T) {
	q := limit.NewQuota(10, 0)
	if !q.CanConsume(6) {
		t.Fatal("expected 6 within 10")
	}
	q.Consume(6)
	if !q.CanConsume(4) {
		t.Fatal("expected 4 exactly at cap")
	}
	q.Consume(4)
	if q.CanConsume(1) {
		t.Fatal("expected token overage to be denied")
	}
	if q.Exceeded() == nil {
		t.Fatal("expected quota exceeded error")
	}
}

func TestQuotaBlocksOnRequests(t *testing.T) {
	q := limit.NewQuota(0, 2)
	for i := 0; i < 2; i++ {
		if !q.CanConsume(1) {
			t.Fatalf("expected allow at %d", i)
		}
		q.Consume(1)
	}
	if q.CanConsume(1) {
		t.Fatal("expected request overage denied")
	}
}

func TestManagerCheckAndTrack(t *testing.T) {
	m := limit.NewManager(limit.Limits{MaxRequestsPerMinute: 2, TokensMonthly: 1000})
	if err := m.CheckAndTrack(100); err != nil {
		t.Fatalf("first: %v", err)
	}
	if err := m.CheckAndTrack(100); err != nil {
		t.Fatalf("second: %v", err)
	}
	if err := m.CheckAndTrack(100); err == nil {
		t.Fatal("expected rate limit after 2 requests")
	}
}

func TestRetryWithBackoffOnRateLimit(t *testing.T) {
	attempts := 0
	err := limit.RetryWithBackoff(context.Background(), 3, time.Millisecond, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("rate limit exceeded")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected retry to recover, got %v", err)
	}
}

func TestIsRateLimitSubstring(t *testing.T) {
	if !limit.IsRateLimit(errors.New("rate limit exceeded")) {
		t.Fatal("expected substring match")
	}
	if limit.IsRateLimit(errors.New("some other error")) {
		t.Fatal("expected no match")
	}
}
