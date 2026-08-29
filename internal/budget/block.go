package budget

import "errors"

// ErrBudgetBlocked indicates execution was halted because a daily/monthly
// budget cap was exceeded (Master Architecture §13, Phase 7).
var ErrBudgetBlocked = errors.New("budget exceeded: execution blocked")

// BlockReporter summarises the blocking state for the caller (e.g. REPL).
type BlockReporter struct {
	DailyUsed   float64
	DailyCap    float64
	MonthlyUsed float64
	MonthlyCap  float64
	Blocked     bool
}
