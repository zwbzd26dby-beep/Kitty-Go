// Package utils provides small reusable helpers (retry, timeouts).
package utils

import (
	"context"
	"time"
)

// Retry invokes fn up to maxAttempts times, waiting backoff between attempts,
// until it returns nil or the context is cancelled.
func Retry(ctx context.Context, maxAttempts int, backoff time.Duration, fn func() error) error {
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if attempt == maxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}
	return err
}
