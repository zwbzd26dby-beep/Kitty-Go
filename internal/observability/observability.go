package observability

import (
	"io"
	"time"
)

// Observability bundles all cross-cutting concerns in one root type
// (Master Architecture §23).
type Observability struct {
	Log     *Logger
	Metrics *Metrics
	Trace   *Trace
	Usage   *Usage
	Cost    *Cost
	Audit   *Audit
	Started time.Time
}

// New creates an Observability root with sane defaults. A nil out writes
// to io.Discard.
func New(out io.Writer) *Observability {
	if out == nil {
		out = io.Discard
	}
	return &Observability{
		Log:     NewLogger(out, LevelInfo),
		Metrics: NewMetrics(),
		Trace:   StartTrace(),
		Usage:   NewUsage(),
		Cost:    NewCost(nil),
		Audit:   NewAudit(),
		Started: time.Now(),
	}
}

// Start marks the current operation in the trace.
func (o *Observability) Start(name string) {
	if o.Trace != nil {
		o.Trace.Start(name)
	}
}

// End finishes the current trace span.
func (o *Observability) End() {
	if o.Trace != nil {
		o.Trace.End()
	}
}

// Uptime returns seconds since procurement.
func (o *Observability) Uptime() float64 {
	return time.Since(o.Started).Seconds()
}
