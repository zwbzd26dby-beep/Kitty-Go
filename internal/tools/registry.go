// Package tools implements the tool registry, allowlist and execution
// framework used by the orchestrator (Master Architecture §17, §22).
package tools

import "context"

// Param describes a single tool argument.
type Param struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

// Tool is a callable capability the agent can invoke.
type Tool interface {
	// Name is the unique identifier used by the model to call the tool.
	Name() string
	// Description explains when to use the tool.
	Description() string
	// Params declares the expected arguments.
	Params() []Param
	// Execute runs the tool with validated string args.
	Execute(ctx context.Context, args map[string]string) (string, error)
}

// Result wraps the outcome of a tool call.
type Result struct {
	Tool    string
	Args    map[string]string
	Output  string
	Errored bool
}

// Registry tracks available tools plus allow/deny state.
type Registry struct {
	tools   map[string]Tool
	allowed map[string]bool
}

// NewRegistry creates an empty tool Registry.
func NewRegistry() *Registry {
	return &Registry{
		tools:   make(map[string]Tool),
		allowed: make(map[string]bool),
	}
}

// Register adds a tool and allows it by default.
func (r *Registry) Register(t Tool) {
	r.tools[t.Name()] = t
	r.allowed[t.Name()] = true
}

// Get returns a tool by name.
func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// List returns all registered tools.
func (r *Registry) List() []Tool {
	out := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, t)
	}
	return out
}

// Allow grants access to a tool.
func (r *Registry) Allow(name string) {
	r.allowed[name] = true
}

// Deny revokes access to a tool.
func (r *Registry) Deny(name string) {
	r.allowed[name] = false
}

// IsAllowed reports whether a tool may be executed.
func (r *Registry) IsAllowed(name string) bool {
	if _, ok := r.tools[name]; !ok {
		return false
	}
	return r.allowed[name]
}
