// Package budget implements daily/monthly spend caps with blocking
// (Master Architecture §13).
package budget

import (
	"fmt"
	"sync"
	"time"
)

// BudgetManager enforces daily and monthly spend limits (Master Arch §13).
type BudgetManager interface {
	// Check verifies that spending cost is allowed this period.
	Check(cost float64) error
	// Spend records a cost against the current period, blocking when over.
	Spend(cost float64) error
	// GetRemaining returns (daily, monthly) remaining budget.
	GetRemaining() (daily, monthly float64)
	// Reset resets period counters (called on period rollover by the owner).
	Reset()
}

// Manager is a thread-safe BudgetManager.
type Manager struct {
	mu          sync.Mutex
	dailyCap    float64
	monthlyCap  float64
	dailyUsed   float64
	monthlyUsed float64
	dayStart    time.Time
	monthStart  time.Time
	now         func() time.Time
	blocked     bool
}

// NewManager builds a BudgetManager with the given daily/monthly caps
// (0 = unlimited).
func NewManager(dailyCap, monthlyCap float64) *Manager {
	return &Manager{
		dailyCap:   dailyCap,
		monthlyCap: monthlyCap,
		dayStart:   time.Now(),
		monthStart: time.Now(),
		now:        time.Now,
	}
}

// Check verifies cost fits within the remaining budget.
func (m *Manager) Check(cost float64) error {
	m.rollover()
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.blocked {
		return ErrBudgetBlocked
	}
	if m.dailyCap > 0 && m.dailyUsed+cost > m.dailyCap {
		return fmt.Errorf("daily budget would be exceeded: %.4f + %.4f > %.4f", m.dailyUsed, cost, m.dailyCap)
	}
	if m.monthlyCap > 0 && m.monthlyUsed+cost > m.monthlyCap {
		return fmt.Errorf("monthly budget would be exceeded: %.4f + %.4f > %.4f", m.monthlyUsed, cost, m.monthlyCap)
	}
	return nil
}

// Spend records cost, blocking the manager if any cap is exceeded.
func (m *Manager) Spend(cost float64) error {
	m.rollover()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dailyUsed += cost
	m.monthlyUsed += cost
	if (m.dailyCap > 0 && m.dailyUsed > m.dailyCap) || (m.monthlyCap > 0 && m.monthlyUsed > m.monthlyCap) {
		m.blocked = true
		return ErrBudgetBlocked
	}
	return nil
}

// GetRemaining returns (daily, monthly) remaining budget.
func (m *Manager) GetRemaining() (daily, monthly float64) {
	m.rollover()
	m.mu.Lock()
	defer m.mu.Unlock()
	d := m.dailyCap
	if d == 0 {
		d = 0
	}
	return sub(m.dailyCap, m.dailyUsed), sub(m.monthlyCap, m.monthlyUsed)
}

// Reset clears usage counters and unblocks.
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dailyUsed = 0
	m.monthlyUsed = 0
	m.blocked = false
	m.dayStart = m.now()
	m.monthStart = m.now()
}

func sub(cap, used float64) float64 {
	if cap <= 0 {
		return 0
	}
	if used >= cap {
		return 0
	}
	return cap - used
}

func (m *Manager) rollover() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	if now.Day() != m.dayStart.Day() {
		m.dailyUsed = 0
		m.dayStart = now
	}
	if now.Month() != m.monthStart.Month() {
		m.monthlyUsed = 0
		m.monthStart = now
		// Budgets roll over; blocked state resets on new period.
		m.blocked = false
	}
}

var _ BudgetManager = (*Manager)(nil)
