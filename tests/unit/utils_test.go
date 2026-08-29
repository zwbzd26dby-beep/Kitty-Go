package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/utils"
)

func TestRetrySucceedsEventually(t *testing.T) {
	attempts := 0
	err := utils.Retry(context.Background(), 3, time.Millisecond, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestRetryExhausts(t *testing.T) {
	attempts := 0
	err := utils.Retry(context.Background(), 3, time.Millisecond, func() error {
		attempts++
		return errors.New("always")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestWrapfPreservesNil(t *testing.T) {
	if got := utils.Wrapf(nil, "ctx"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
	e := errors.New("base")
	wrapped := utils.Wrapf(e, "ctx %d", 1)
	if !errors.Is(wrapped, e) {
		t.Fatal("expected wrapped to unwrap to base")
	}
}
