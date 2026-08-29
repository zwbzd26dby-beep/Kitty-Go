package observability

import (
	"sync"
	"time"
)

// Metrics is a thread-safe counter/gauges store.
type Metrics struct {
	mu       sync.Mutex
	counters map[string]int64
	gauges   map[string]float64
	start    time.Time
}

// NewMetrics creates a Metrics store.
func NewMetrics() *Metrics {
	return &Metrics{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
		start:    time.Now(),
	}
}

// Inc increments a counter by 1.
func (m *Metrics) Inc(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name]++
}

// Add increments a counter by delta.
func (m *Metrics) Add(name string, delta int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += delta
}

// Set sets a gauge value.
func (m *Metrics) Set(name string, v float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = v
}

// Get returns a counter value (0 if missing).
func (m *Metrics) Get(name string) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.counters[name]
}

// Snapshot returns copies of counters and gauges.
func (m *Metrics) Snapshot() (map[string]int64, map[string]float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cs := make(map[string]int64, len(m.counters))
	for k, v := range m.counters {
		cs[k] = v
	}
	gs := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		gs[k] = v
	}
	return cs, gs
}

// Reset clears counters and gauges (uptime preserved).
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters = make(map[string]int64)
	m.gauges = make(map[string]float64)
}
