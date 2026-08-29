// Package execution implements the Execution Layer. In Phase 1 it exposes
// the Executor interface and a local execution path that delegates to an
// LLM provider. Remote/distributed execution arrives in Phases 9 and 11.
package execution

import (
	"context"

	"github.com/zwbzd26dby-beep/Kitty-Go/pkg/types"
)

// Executor runs a task in a chosen location (local in Phase 1).
type Executor interface {
	// ExecuteLocal runs the task locally using the provided LLM client.
	ExecuteLocal(ctx context.Context, task types.Task, client LLMClient) (*Result, error)
	// GetStatus returns the current status of a job.
	GetStatus(jobID string) (*JobStatus, error)
	// Cancel cancels a running job.
	Cancel(jobID string) error
}

// LLMClient is the minimal generation contract needed by the local executor.
type LLMClient interface {
	Generate(ctx context.Context, msg types.Message, history []types.Turn) (types.Response, error)
}

// Result is the outcome of executing a task.
type Result struct {
	TaskID  string
	Content string
	JobID   string
	Err     error
}
