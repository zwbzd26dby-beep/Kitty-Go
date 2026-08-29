package costtest

import (
	"errors"
	"testing"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/budget"
)

func TestBudgetCheckUnderCap(t *testing.T) {
	b := budget.NewDaily(10.0)
	if err := b.Check(5.0); err != nil {
		t.Fatalf("expected check to pass, got %v", err)
	}
}

func TestBudgetCheckOverCap(t *testing.T) {
	b := budget.NewDaily(10.0)
	if err := b.Check(15.0); err == nil {
		t.Fatal("expected budget check to fail")
	}
}

func TestBudgetSpendBlocks(t *testing.T) {
	b := budget.NewDaily(10.0)
	if err := b.Spend(6.0); err != nil {
		t.Fatalf("spend 6: %v", err)
	}
	if err := b.Spend(4.0); err != nil {
		t.Fatalf("spend 4 (exactly to cap): %v", err)
	}
	// Pulling over cap blocks.
	if err := b.Spend(1.0); !errors.Is(err, budget.ErrBudgetBlocked) {
		t.Fatalf("expected ErrBudgetBlocked, got %v", err)
	}
}

func TestBudgetMonthlyBlocks(t *testing.T) {
	b := budget.NewMonthly(100.0)
	b.Spend(80.0)
	if err := b.Spend(30.0); !errors.Is(err, budget.ErrBudgetBlocked) {
		t.Fatalf("expected block on monthly overage, got %v", err)
	}
}

func TestBudgetRemaining(t *testing.T) {
	b := budget.NewDaily(10.0)
	b.Spend(4.0)
	d, m := b.GetRemaining()
	if d != 6.0 {
		t.Fatalf("expected daily remaining 6, got %v", d)
	}
	if m != 0 {
		t.Fatalf("no monthly cap => remaining 0, got %v", m)
	}
}

func TestBudgetResetUnblocks(t *testing.T) {
	b := budget.NewDaily(5.0)
	b.Spend(3.0)
	if err := b.Spend(3.0); !errors.Is(err, budget.ErrBudgetBlocked) {
		t.Fatalf("expected block, got %v", err)
	}
	b.Reset()
	if err := b.Spend(2.0); err != nil {
		t.Fatalf("expected unblocked after reset, got %v", err)
	}
}
