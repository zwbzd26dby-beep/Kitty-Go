package cost

import (
	"fmt"
	"sync"
)

// Alert is invoked whenever a configured cost threshold is crossed.
type Alert func(model string, total float64)

// ThresholdAlert fires an Alert callback the first time cumulative cost for a
// model exceeds a threshold since the last reset.
type ThresholdAlert struct {
	mu        sync.Mutex
	threshold float64
	fired     map[string]bool
	onAlert   Alert
}

// NewThresholdAlert builds an alert that fires when a model's cost exceeds
// threshold, invoking onAlert once per crossing.
func NewThresholdAlert(threshold float64, onAlert Alert) *ThresholdAlert {
	return &ThresholdAlert{threshold: threshold, fired: make(map[string]bool), onAlert: onAlert}
}

// Check inspects a cost event and fires the alert if the threshold is crossed.
func (a *ThresholdAlert) Check(model string, cost float64) {
	if a.threshold <= 0 || a.onAlert == nil {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if cost >= a.threshold && !a.fired[model] {
		a.fired[model] = true
		a.onAlert(model, cost)
	}
}

// Error returns an error explaining that the spend threshold was exceeded.
func (a *ThresholdAlert) Error(model string, cost float64) error {
	return fmt.Errorf("cost threshold exceeded for %q: %.4f >= %.4f", model, cost, a.threshold)
}
