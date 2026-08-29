package observability

import (
	"sync"
)

// Cost tracks estimated spend per model.
type Cost struct {
	mu      sync.Mutex
	prices  map[string]PricePoint
	spend   map[string]float64
	records []CostRecord
}

// PricePoint is per-1M-token pricing for a model.
type PricePoint struct {
	PromptIn     float64
	CompletionIn float64
}

// CostRecord is an immutable snapshot of a single cost event.
type CostRecord struct {
	Model string
	Cost  float64
}

// NewCost creates a Cost tracker, optionally seeded with prices.
func NewCost(prices map[string]PricePoint) *Cost {
	c := &Cost{
		prices: make(map[string]PricePoint),
		spend:  make(map[string]float64),
	}
	for k, v := range prices {
		c.prices[k] = v
	}
	return c
}

// Price registers pricing for a model.
func (c *Cost) Price(model string, promptIn, completionIn float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prices[model] = PricePoint{PromptIn: promptIn, CompletionIn: completionIn}
}

// LogP records estimated cost for a usage event.
func (c *Cost) LogP(model string, prompt, completion int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	pp := c.prices[model]
	cost := float64(prompt)*pp.PromptIn/1e6 + float64(completion)*pp.CompletionIn/1e6
	c.spend[model] += cost
	c.records = append(c.records, CostRecord{Model: model, Cost: cost})
}

// TotalSpend returns accumulated spend per model.
func (c *Cost) TotalSpend() map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]float64, len(c.spend))
	for k, v := range c.spend {
		out[k] = v
	}
	return out
}

// Records returns a copy of all cost events.
func (c *Cost) Records() []CostRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]CostRecord, len(c.records))
	copy(out, c.records)
	return out
}
