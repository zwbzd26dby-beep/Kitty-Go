package observability

import (
	"sync"
)

// Usage tracks token usage per model.
type Usage struct {
	mu      sync.Mutex
	total   map[string]UsagePoint
	records []UsageRecord
}

// UsagePoint aggregates tokens for a model.
type UsagePoint struct {
	PromptTokens     int64
	CompletionTokens int64
	Calls            int64
}

// UsageRecord is an immutable snapshot of a single usage event.
type UsageRecord struct {
	Model            string
	PromptTokens     int64
	CompletionTokens int64
}

// NewUsage creates a Usage tracker.
func NewUsage() *Usage {
	return &Usage{total: make(map[string]UsagePoint)}
}

// LogP records a usage event.
func (u *Usage) LogP(model string, prompt, completion int64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	p := u.total[model]
	p.PromptTokens += prompt
	p.CompletionTokens += completion
	p.Calls++
	u.total[model] = p
	u.records = append(u.records, UsageRecord{Model: model, PromptTokens: prompt, CompletionTokens: completion})
}

// Total returns aggregated usage per model.
func (u *Usage) Total() map[string]UsagePoint {
	u.mu.Lock()
	defer u.mu.Unlock()
	out := make(map[string]UsagePoint, len(u.total))
	for k, v := range u.total {
		out[k] = v
	}
	return out
}

// Records returns a copy of all usage events.
func (u *Usage) Records() []UsageRecord {
	u.mu.Lock()
	defer u.mu.Unlock()
	out := make([]UsageRecord, len(u.records))
	copy(out, u.records)
	return out
}
