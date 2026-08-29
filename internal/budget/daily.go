package budget

// NewDaily returns a BudgetManager enforcing only a daily cap.
func NewDaily(cap float64) BudgetManager {
	return NewManager(cap, 0)
}
