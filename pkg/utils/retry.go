// Package utils provides small reusable helpers (retry, timeouts).
package utils

import (
	"context"
	"time"
)

// Retry invokes fn up to maxAttempts times, waiting a fixed backoff between
// attempts, until it returns nil or the context is cancelled.
func Retry(ctx context.Context, maxAttempts int, backoff time.Duration, fn func() error) error {
	return RetryWith(ctx, RetryPolicy{
		MaxAttempts: maxAttempts,
		Initial:     backoff,
		MaxBackoff:  backoff,
	}, nil, fn)
}

// RetryPolicy configures exponential-backoff retry behaviour.
type RetryPolicy struct {
	// MaxAttempts is the maximum number of executions of fn.
	MaxAttempts int
	// Initial is the first backoff wait.
	Initial time.Duration
	// MaxBackoff caps the exponential growth of the backoff wait.
	MaxBackoff time.Duration
	// Multiplier grows the backoff between attempts (defaults to 2.0).
	Multiplier float64
}

// RetryWith invokes fn up to MaxAttempts times with exponential backoff,
// until it returns nil, a non-retryable error, or the context is cancelled.
// When isRetryable is nil, every error is treated as retryable.
func RetryWith(ctx context.Context, policy RetryPolicy, isRetryable func(error) bool, fn func() error) error {
	wait := policy.Initial
	if wait <= 0 {
		wait = time.Millisecond
	}
	max := policy.MaxBackoff
	if max <= 0 {
		max = wait
	}
	m := policy.Multiplier
	if m <= 0 {
		m = 2.0
	}

	var err error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if isRetryable != nil && !isRetryable(err) {
			return err
		}
		if attempt == policy.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
		wait = nextBackoff(wait, max, m)
	}
	return err
}

func nextBackoff(wait, max time.Duration, m float64) time.Duration {
	next := time.Duration(float64(wait) * m)
	if next > max || next <= wait {
		return max
	}
	return next
}
