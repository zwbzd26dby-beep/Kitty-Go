package cost

import "sync"

// CostManager is the central cost tracker (Master Architecture §11).
type CostManager interface {
	// TrackCost records a cost event and returns the running total.
	TrackCost(rec Record) (float64, error)
	// GetTotalCost returns the cumulative cost.
	GetTotalCost() float64
	// GetModelCost returns the cumulative cost attributed to a model.
	GetModelCost(model string) float64
	// Reset clears all tracked cost.
	Reset()
}

// Manager is a thread-safe in-memory CostManager.
type Manager struct {
	mu      sync.Mutex
	total   float64
	byModel map[string]float64
	calc    Calculator
}

// NewManager creates an empty CostManager.
func NewManager() *Manager {
	return &Manager{byModel: make(map[string]float64)}
}

// TrackCost records a cost event.
func (m *Manager) TrackCost(rec Record) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.total += rec.Cost
	m.byModel[rec.Model] += rec.Cost
	return m.total, nil
}

// GetTotalCost returns cumulative cost.
func (m *Manager) GetTotalCost() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.total
}

// GetModelCost returns cumulative cost for a model.
func (m *Manager) GetModelCost(model string) float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.byModel[model]
}

// Reset clears tracked cost.
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.total = 0
	m.byModel = make(map[string]float64)
}

var _ CostManager = (*Manager)(nil)
