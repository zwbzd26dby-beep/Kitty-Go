// Package orchestrator coordinates a task through decision, model routing,
// compute routing and execution. In Phase 1 it provides a skeleton that
// passes tasks straight through the local executor; the full flow arrives
// in Phase 17.
package orchestrator

import (
	"context"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/execution"
	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// Orchestrator coordinates the execution of a task.
type Orchestrator interface {
	// Orchestrate runs a task and returns its result.
	Orchestrate(ctx context.Context, task types.Task) (*execution.Result, error)
}

// simpleOrchestrator is the Phase 1 implementation: a pass-through to the
// local executor using an LLM client.
type simpleOrchestrator struct {
	exec   execution.Executor
	client execution.LLMClient
}

// New creates an Orchestrator with the given executor and LLM client.
func New(exec execution.Executor, client execution.LLMClient) Orchestrator {
	return &simpleOrchestrator{exec: exec, client: client}
}

// Orchestrate executes a task locally.
func (o *simpleOrchestrator) Orchestrate(ctx context.Context, task types.Task) (*execution.Result, error) {
	return o.exec.ExecuteLocal(ctx, task, o.client)
}
