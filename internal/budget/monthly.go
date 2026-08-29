package budget

// NewMonthly returns a BudgetManager enforcing only a monthly cap.
func NewMonthly(cap float64) BudgetManager {
	return NewManager(0, cap)
}
