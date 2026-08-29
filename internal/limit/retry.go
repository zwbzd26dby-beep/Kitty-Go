package limit

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/utils"
)

// IsRateLimit reports whether err is a soft rate-limit error that warrants a
// retry after backoff.
func IsRateLimit(err error) bool {
	var e errLimitExceeded
	if errors.As(err, &e) {
		return true
	}
	return strings.Contains(err.Error(), "rate limit")
}

// RetryWithBackoff runs fn, retrying rate-limit errors with exponential
// backoff up to maxAttempts times.
func RetryWithBackoff(ctx context.Context, maxAttempts int, initial time.Duration, fn func() error) error {
	return utils.RetryWith(ctx, utils.RetryPolicy{
		MaxAttempts: maxAttempts,
		Initial:     initial,
		MaxBackoff:  5 * time.Second,
		Multiplier:  2,
	}, IsRateLimit, fn)
}
