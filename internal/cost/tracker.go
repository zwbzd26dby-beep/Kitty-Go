package cost

// Tracker combines Calculator and CostManager to record a request's usage as
// a cost (the cost-tracking half of Master Architecture §35).
type Tracker struct {
	calc Calculator
	mgr  CostManager
}

// NewTracker builds a Tracker over mgr.
func NewTracker(mgr CostManager) *Tracker {
	return &Tracker{calc: Calculator{}, mgr: mgr}
}

// TrackUsage computes cost from usage+pricing and records it against the
// manager, returning the new running total.
func (t *Tracker) TrackUsage(provider, model string, u Usage, p Pricing) (float64, error) {
	c := t.calc.Calculate(u, p)
	return t.mgr.TrackCost(Record{Provider: provider, Model: model, Usage: u, Cost: c})
}
