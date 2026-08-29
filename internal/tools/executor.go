package tools

import (
	"context"
	"errors"
	"sort"
	"strings"
)

// ErrToolNotFound is returned when the tool does not exist.
var ErrToolNotFound = errors.New("tool not found")

// ErrToolDenied is returned when the tool is blocked by the allowlist.
var ErrToolDenied = errors.New("tool denied by allowlist")

// Executor runs tools through the registry with allowlist enforcement.
type Executor struct {
	reg *Registry
	// simulate records the call without running Execute.
	simulate bool
}

// NewExecutor creates an Executor over a registry.
func NewExecutor(reg *Registry) *Executor {
	return &Executor{reg: reg}
}

// Simulate returns an executor that records intent without executing.
func (e *Executor) Simulate() *Executor {
	return &Executor{reg: e.reg, simulate: true}
}

// Execute invokes a named tool, validating args against its schema.
func (e *Executor) Execute(ctx context.Context, name string, args map[string]string) (*Result, error) {
	if !e.reg.IsAllowed(name) {
		if _, ok := e.reg.Get(name); !ok {
			return nil, ErrToolNotFound
		}
		return nil, ErrToolDenied
	}
	t, _ := e.reg.Get(name)
	if err := validateArgs(t, args); err != nil {
		return nil, err
	}
	if e.simulate {
		return &Result{Tool: name, Args: args, Output: "(simulated)"}, nil
	}
	out, err := t.Execute(ctx, args)
	if err != nil {
		return &Result{Tool: name, Args: args, Errored: true, Output: err.Error()}, nil
	}
	return &Result{Tool: name, Args: args, Output: out}, nil
}

// Called renders the last tool call in the model-facing function-call format.
func (r *Result) Called() string {
	keys := make([]string, 0, len(r.Args))
	for k := range r.Args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString("-> " + r.Tool + "(")
	for i, k := range keys {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(k + "=" + r.Args[k])
	}
	sb.WriteString(")")
	if r.Errored {
		sb.WriteString(" [error: " + r.Output + "]")
	}
	return sb.String()
}

func validateArgs(t Tool, args map[string]string) error {
	required := map[string]bool{}
	for _, p := range t.Params() {
		if p.Required {
			required[p.Name] = true
		}
	}
	for name := range required {
		if _, ok := args[name]; !ok {
			return errors.New("missing required arg: " + name)
		}
	}
	return nil
}
