package observability

import (
	"sync"
	"time"
)

// AuditEvent records an auditable action (who, what, when).
type AuditEvent struct {
	Time   time.Time
	Actor  string
	Action string
	Detail string
}

// Audit is an append-only audit log.
type Audit struct {
	mu      sync.Mutex
	events  []AuditEvent
	onEvent func(AuditEvent)
}

// NewAudit creates an empty audit log.
func NewAudit() *Audit {
	return &Audit{}
}

// OnEvent sets a sink invoked for every audit event. Sinks run synchronously.
func (a *Audit) OnEvent(fn func(AuditEvent)) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.onEvent = fn
}

// Record appends an audit event.
func (a *Audit) Record(actor, action, detail string) {
	ev := AuditEvent{Time: time.Now(), Actor: actor, Action: action, Detail: detail}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, ev)
	if a.onEvent != nil {
		a.onEvent(ev)
	}
}

// Events returns a copy of all audit events.
func (a *Audit) Events() []AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]AuditEvent, len(a.events))
	copy(out, a.events)
	return out
}

// Filter returns events matching the action.
func (a *Audit) Filter(action string) []AuditEvent {
	all := a.Events()
	var out []AuditEvent
	for _, e := range all {
		if e.Action == action {
			out = append(out, e)
		}
	}
	return out
}
