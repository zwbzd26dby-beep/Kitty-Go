package observability

import (
	"context"
	"sync"
	"time"
)

type traceKey struct{}

// span is a timed operation within a trace.
type span struct {
	Name     string
	Start    time.Time
	Duration time.Duration
}

// Trace is a collection of spans.
type Trace struct {
	mu     sync.Mutex
	spans  []span
	active []span
}

// StartTrace begins a new trace.
func StartTrace() *Trace {
	return &Trace{}
}

// WithTrace attaches a trace to ctx.
func WithTrace(ctx context.Context, t *Trace) context.Context {
	return context.WithValue(ctx, traceKey{}, t)
}

// TraceFromCtx returns the trace stored in ctx, if any.
func TraceFromCtx(ctx context.Context) (*Trace, bool) {
	t, ok := ctx.Value(traceKey{}).(*Trace)
	return t, ok
}

// Start begins a named span.
func (t *Trace) Start(name string) {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.active = append(t.active, span{Name: name, Start: time.Now()})
}

// End finishes the most recent span.
func (t *Trace) End() {
	if t == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.active) == 0 {
		return
	}
	last := t.active[len(t.active)-1]
	t.active = t.active[:len(t.active)-1]
	last.Duration = time.Since(last.Start)
	t.spans = append(t.spans, last)
}

// Drain returns completed spans and clears them.
func (t *Trace) Drain() []SpanInfo {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]SpanInfo, 0, len(t.spans))
	for _, s := range t.spans {
		out = append(out, SpanInfo{Name: s.Name, Duration: s.Duration})
	}
	t.spans = t.spans[:0]
	return out
}

// SpanInfo is the exported, immutable span shape.
type SpanInfo struct {
	Name     string
	Duration time.Duration
}

var _ = SpanInfo{}
