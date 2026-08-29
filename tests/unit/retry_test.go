package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/utils"
)

func TestRetryWithNonRetryableStops(t *testing.T) {
	attempts := 0
	boom := errors.New("permanent")
	err := utils.RetryWith(context.Background(), utils.RetryPolicy{
		MaxAttempts: 5,
		Initial:     time.Millisecond,
		MaxBackoff:  4 * time.Millisecond,
	}, func(err error) bool { return false }, func() error {
		attempts++
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("expected boom, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-retryable error, got %d", attempts)
	}
}

func TestRetryWithRetryableSucceeds(t *testing.T) {
	attempts := 0
	err := utils.RetryWith(context.Background(), utils.RetryPolicy{
		MaxAttempts: 4,
		Initial:     time.Millisecond,
		MaxBackoff:  2 * time.Millisecond,
	}, func(err error) bool { return true }, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryContextCancelled(t *testing.T) {
	attempts := 0
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := utils.RetryWith(ctx, utils.RetryPolicy{
		MaxAttempts: 10,
		Initial:     10 * time.Millisecond,
		MaxBackoff:  10 * time.Millisecond,
	}, nil, func() error {
		attempts++
		return errors.New("transient")
	})
	if err == nil {
		t.Fatal("expected context error")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt before cancellation, got %d", attempts)
	}
}
